package updater

import (
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/defaults"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/execute"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/statistic"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

type Updater struct {
	r       Reporter
	service string
	busy    bool
	uType   util.UpdaterType
	rem     execute.Remoter
	args    map[string]string
}

var once sync.Once //nolint:gochecknoglobals
var m sync.Mutex   //nolint:gochecknoglobals

func updater(k string) (*Updater, bool) {
	m.Lock()
	defer m.Unlock()
	once.Do(func() {
		updatersM = make(map[string]*Updater)
	})

	v, ok := updatersM[k]

	return v, ok
}

func setUpdater(k string, u *Updater) {
	m.Lock()
	defer m.Unlock()

	updatersM[k] = u
}

func Get(service string, uType util.UpdaterType, args, sshData map[string]string) *Updater {
	u, ok := updater(service)
	if !ok {
		u = newUpdater(service)
		setUpdater(service, u)
	}

	u.SetArgs(args)
	u.setRemoter(sshData)

	u.uType = uType
	updatersM[service].uType = uType

	return u
}

func newUpdater(service string) *Updater {
	return &Updater{
		service: service,
		args:    map[string]string{},
	}
}

func (u *Updater) SetArgs(m map[string]string) {
	for k, v := range m {
		u.args[k] = v
	}
}

func (u *Updater) setRemoter(sshData map[string]string) {
	u.rem = execute.RemoterFromData(sshData)
}

func (u *Updater) Update(progress chan util.Progress) error {
	defer func() {
		err := statistic.Write()
		if err != nil {
			log.Errorln(err)
		}
	}()

	if _, ok := u.args["dir"]; !ok && u.uType == util.UpdateLocal {
		return fmt.Errorf("failed to get dir for service: %s", u.service)
	}

	if err := u.setUp(progress); err != nil {
		return err
	}

	curScript, err := u.GetScript()
	if err != nil {
		return err
	}

	// Выполняем полученный скрипт.
	// если хотя-бы одна команда возвращает ошибку, сразу вернем её для отображения.
	return u.execute(curScript)
}

func (u *Updater) GetScript() (script.Script, error) {
	// Сначала ищем в памяти скрипт
	scriptName := defaults.GetScriptName(u.service, u.uType)

	curScript, found := script.Get(scriptName)
	if found {
		return curScript, nil
	}

	log.Warnf("Script for %s not found", scriptName)

	// Не нашли скрипта. Посмотрим в скриптах по умолчанию
	if defaults.IsServiceSupported(u.service, u.uType) {
		// Нашли скрипт по умолчанию
		log.Infoln("Using default")

		return defaults.ScriptForService(u.service, u.uType), nil
	}

	return script.Script{}, defaults.NotSupportedError{}
}

func (u *Updater) setUp(progress chan util.Progress) error {
	// Первичная настройка статуса обновления (сбрасываем с предыдущего раза)
	// Заодно устанавливаем канал для передачи информации назад
	err := u.r.setUp(progress)
	if err != nil {
		return err
	}

	if u.rem.Addr == "" && u.uType != util.UpdateInPlace {
		return fmt.Errorf("please set up ssh address in settings window")
	}

	// Проверим наличие rsync и sshpass при необходимости
	needSSHPass := u.rem.Pass != ""
	if u.uType != util.UpdateInPlace && needSSHPass {
		s := defaults.SetDefaultIfNotExistsAndReturn("util_check_sshpass_exists", defaults.CheckSSHPassExists)
		if err := u.execute(s); err != nil {
			return fmt.Errorf("please install sshpass to connect with password")
		}
	}

	if u.uType == util.UpdateLocal {
		s := defaults.SetDefaultIfNotExistsAndReturn("util_check_rsync_exists", defaults.CheckRsyncExists)
		if err := u.execute(s); err != nil {
			return fmt.Errorf("please install rsync locally and on remote")
		}
	}

	return nil
}

// Общий случай выполнения, когда инклюды еще не известны.
func (u *Updater) execute(s script.Script) error {
	return u.executeWithIncludes(s, map[string]struct{}{})
}

// Частный случай выполнения скрипта, подключенного из другого скрипта.
func (u *Updater) executeWithIncludes(s script.Script, includes map[string]struct{}) error {
	log.Debugf("Script %v ", s)

	var err error

	for _, c := range s.Commands {
		// Аргументы у команды могут как быть так и отсутствовать
		if c.Args != nil {
			// Если они есть, мы их попробуем заполнить в следующем приоритете:
			// 1. Ключ в config.json для этого сервиса
			// 2. Аргумент по умолчанию в скрипте
			// 3. Сам текст ключа аргумента (вместо "branch" подставить "branch")
			// 0. "defer" это не аргумент, а инструкция
			c.Args = s.ExtendArgsByDefaultAndGiven(c.Args, u.service, u.args)
		}

		if c.IsDeferred() {
			log.Debugf("Deferring command: %v ", c)
			defer func(c command.Command) {
				err = u.executeByType(c, includes)
			}(c)

			continue
		}

		if err = u.executeByType(c, includes); err != nil {
			return err
		}
	}

	return err
}

// Тут бы лучше использовать интерфейс с методом process.
func (u *Updater) executeByType(c command.Command, includes map[string]struct{}) error {
	var err error
	// Выполняем команду в зависимости от типа.
	// Заполняем процесс обновления
	switch c.Type {
	case command.SSH:
		err = u.processWithMeanSSH(c)
	case command.RSYNC:
		err = u.processWithMeanRsync(c)
	case command.LOCAL:
		err = u.processWithMeanLocal(c)
	// Команда может быть на подключение другого скрипта
	case command.INCLUDE:
		// Проверим что скрипт еще не подключался
		if _, ok := includes[c.Data]; ok {
			return fmt.Errorf("include circle is forbidden, current includes: %v, including: %s", getKeys(includes), c.Data)
		}

		includes[c.Data] = struct{}{}

		// Берем скрипт
		curScript, ok := script.Get(c.Data)
		if !ok {
			return fmt.Errorf("included script %s not found", c.Data)
		}

		log.Infof("Including script: %s", c.Data)
		// И просто выполняем, передавая в него текущий набор подключенных скриптов
		// Он виртуально вставляется в этом месте
		err = u.executeWithIncludes(curScript, includes)
	}

	return err
}

func (u *Updater) processWithMeanSSH(c command.Command) error {
	return u.processWithMeanFn(
		func(cmd string) error {
			out, err := execute.BySSH(u.rem, cmd)
			if err != nil {
				log.Errorf("cmd: %s, out: %s, err: %s", cmd, out, err)

				return fmt.Errorf("cmd: %v, err: %w", cmd, err)
			}

			return nil
		}, c)
}

var updatersM map[string]*Updater //nolint:gochecknoglobals
const rsyncArgsCount = 2

func (u *Updater) processWithMeanRsync(c command.Command) error {
	return u.processWithMeanFn(
		func(string) error {
			if len(c.Args) != rsyncArgsCount {
				return fmt.Errorf("expected %d arguments, got %d", rsyncArgsCount, len(c.Args))
			}
			err := u.uploadByRsync(c.Args[1], c.Args[0])
			if err != nil {
				return err
			}

			return nil
		}, c)
}

func (u *Updater) processWithMeanLocal(c command.Command) error {
	return u.processWithMeanFn(
		func(cmd string) error {
			return u.executeLocally(cmd)
		}, c)
}

func (u *Updater) processWithMeanFn(fn func(cmd string) error, c command.Command) error {
	// Передаем текущий прогресс
	u.r.report(c.Msg, u.service, u.uType)
	// Заполняем аргументы если они там есть
	cmd := c.GetResultingCmdWithArgs()
	// Выполняем команду
	start := time.Now()

	if err := fn(cmd); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			log.Errorf("cmd: %s, out: %s, err: %s", cmd, ee.Stderr, ee)
			err = fmt.Errorf("err: %s", ee.Stderr)
		}

		return err
	}
	// Замеряя время выполнения
	// u.r.updateMeanStage(time.Since(start), u.service, u.uType)
	u.r.updateMeanStage(time.Since(start), u.service, u.uType, c.Msg)

	return nil
}

func (u *Updater) executeLocally(cmd string) error {
	_, err := execute.Local(cmd)

	return err
}

func (u *Updater) uploadByRsync(file string, extPath string) error {
	dir, ok := u.args["dir"]
	// dir, ok := config.Get(u.service, "dir")
	if !ok {
		log.Errorf("failed to get dir for service: %s", u.service)

		return fmt.Errorf("failed to get dir for service: %s", u.service)
	}

	r := execute.RemoterFromConfig()
	if extPath == "" {
		return execute.Rsync(r, dir, file, u.service)
	}

	return execute.Rsync(r, dir+"/"+extPath, file, u.service)
}

func (u Updater) IsBusy() bool {
	return u.busy
}

func (u *Updater) SetBusy(b bool) {
	u.busy = b
}

func getKeys(includes map[string]struct{}) []string {
	i := 0
	keys := make([]string, len(includes))

	for k := range includes {
		keys[i] = k
		i++
	}

	return keys
}
