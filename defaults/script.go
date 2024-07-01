package defaults

import (
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
)

func SetDefaultIfNotExists(scriptName string, fn func() script.Script) {
	if _, ok := script.Get(scriptName); !ok {
		script.Set(scriptName, fn())
	}
}

var defaultMux sync.Mutex //nolint:gochecknoglobals

func SetDefaultIfNotExistsAndReturn(scriptName string, fn func() script.Script) script.Script {
	curScript, ok := script.Get(scriptName)
	if ok {
		return curScript
	}

	defaultMux.Lock()
	defer defaultMux.Unlock()

	script.Set(scriptName, fn())

	return fn()
}

func CheckRsyncExists() script.Script {
	return GenerateScript(CheckRsyncLocal, CheckRsyncRemote)
}

func CheckSSHPassExists() script.Script {
	return GenerateScript(CheckSSHPass)
}

func CheckLSExists() script.Script {
	return GenerateScript(CheckLS)
}

func SupervisorStuff() script.Script {
	return GenerateScript(RereadSuper, UpdateSuper).AddArgs([]string{"container_name"}, []string{"vm_box"})
}

type DefaultInclude struct {
	Include string
	Fn      func() script.Script
}

func IncludesScript(includes ...DefaultInclude) script.Script {
	commands := []command.Command{}

	for _, i := range includes {
		logrus.Debugln("Creating default include script", i.Include)
		SetDefaultIfNotExists(i.Include, i.Fn)
		commands = append(commands, command.MakeInclude(i.Include))
	}

	return script.Script{
		Ver:      1,
		Commands: commands,
	}
}
