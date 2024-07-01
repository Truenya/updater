package script

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

type Script struct {
	Ver      int               `json:"ver"`
	Name     string            `json:"-"`
	Args     map[string]string `json:"args"`
	Commands []command.Command `json:"commands"`
}

var _scripts map[string]Script //nolint:gochecknoglobals
var dirPath = "./scripts/"     //nolint:gochecknoglobals
var m sync.Mutex               //nolint:gochecknoglobals

func GetUserPath() (string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		logrus.Errorln(err)

		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}

	return config, nil
}

func Init() {
	path, err := GetUserPath()
	if err == nil {
		dirPath = path + "/updater/scripts/"
	} else {
		logrus.Errorln("using", "./scripts/", err)
	}

	if _, err = os.Stat(dirPath); err != nil {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			panic(err)
		}
	}

	services, err := os.ReadDir(dirPath)
	util.LogNotNilErr(err)

	_scripts = make(map[string]Script)

	for _, e := range services {
		if !e.IsDir() {
			panic("\n\nBreaking change in update.\nPlease delete all scripts in scripts dir: " + dirPath)
		}

		serviceDir := e.Name()

		scripts, err := os.ReadDir(dirPath + serviceDir)
		if err != nil {
			logrus.Errorln(err)

			continue
		}

		for _, s := range scripts {
			fileName := s.Name()
			script := Script{}
			fullPath := fmt.Sprintf("%s%s/%s", dirPath, serviceDir, fileName)

			f, err := os.OpenFile(fullPath, os.O_RDONLY, 0755)
			if err != nil {
				logrus.Panicf("Failed to open script %s for init %s", fullPath, err)
			}
			defer f.Close()

			err = json.NewDecoder(f).Decode(&script)
			if err != nil {
				logrus.Panicf("Failed to read script %s for init %s", fullPath, err)
			}

			scriptName := serviceDir + "_" + strings.Split(fileName, ".")[0]
			logrus.Debugf("Script %s read successfully", scriptName)
			script.Name = scriptName // Для проверок в include команде

			_scripts[scriptName] = script
		}
	}
}

func Get(name string) (Script, bool) {
	m.Lock()
	defer m.Unlock()

	v, ok := _scripts[name]
	v.Name = name

	return v, ok
}

func CreateFileIfNotExists(file string) {
	if _, err := os.Stat(file); err != nil {
		m.Lock()
		defer m.Unlock()

		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		util.LogNotNilErr(err)
		_, err = f.WriteString("{}")
		util.WarnNotNilErr(err)
		f.Close()
	}
}

func CheckService(scriptName string) string {
	splitted := strings.Split(scriptName, "_")
	service := splitted[0]

	servicePath := fmt.Sprintf("%s/%s", dirPath, service)
	if _, err := os.Stat(servicePath); err != nil {
		m.Lock()
		defer m.Unlock()
		os.MkdirAll(servicePath, 0755) //nolint:errcheck
	}

	return fmt.Sprintf("%s/%s.json", service, strings.Join(splitted[1:], "_"))
}

func Set(name string, c Script) {
	scriptRelPath := CheckService(name)
	scriptFullPath := fmt.Sprintf("%s/%s", dirPath, scriptRelPath)
	CreateFileIfNotExists(scriptFullPath)

	c.Name = name

	m.Lock()
	defer m.Unlock()

	_scripts[name] = c
	f, err := os.OpenFile(scriptFullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	util.LogNotNilErr(err)

	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", " ")

	if err = e.Encode(c); err != nil {
		logrus.Warn(err)
	}
}

func (s Script) ExtendArgsByDefaultAndGiven(args []string, service string, uArgs map[string]string) []string {
	for i, arg := range args {
		// Not exactly arg, but instruction
		if arg == util.Defer {
			continue
		}

		logrus.Debugf("Looking for arg %s in service %s", arg, service)

		if TrySubstituteFromArgs(args, i, uArgs) {
			continue
		}

		if s.TrySubstituteFromScript(args, i) {
			continue
		}

		logrus.Debugf("Not found in config and script, using exact arg, %s", arg)
	}

	return args
}

func TrySubstituteFromArgs(args []string, i int, uArgs map[string]string) bool {
	val, ok := uArgs[args[i]]
	if ok {
		args[i] = val
		logrus.Debugf("found %s", val)
	}

	return ok
}

func (s Script) TrySubstituteFromScript(args []string, i int) bool {
	val, ok := s.Args[args[i]]
	if ok {
		args[i] = val
		logrus.Debugf("found %s", val)
	}

	return ok
}

func (s Script) AddArgs(k, v []string) Script {
	if s.Args == nil {
		s.Args = make(map[string]string)
	}

	for i, key := range k {
		s.Args[key] = v[i]
	}

	return s
}

func (s Script) AddDefaultRemoteArgs(service string) Script {
	k := []string{"branch", "container_name"}
	v := []string{"master", service + "_update"}

	if s.Args == nil {
		s.Args = make(map[string]string)
	}

	for i, key := range k {
		s.Args[key] = v[i]
	}

	return s
}
