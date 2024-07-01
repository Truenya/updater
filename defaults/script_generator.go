package defaults

import (
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
)

func GenerateScript(fns ...func() command.Command) script.Script {
	commands := []command.Command{}
	for _, fn := range fns {
		commands = append(commands, fn())
	}

	return script.Script{
		Ver:      1,
		Commands: commands,
	}
}
