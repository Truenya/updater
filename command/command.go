package command

import "fmt"

type ctype int8

const (
	SSH ctype = iota
	RSYNC
	LOCAL
	INCLUDE
)

type Command struct {
	Msg  string
	Type ctype
	Args Argst
	Data string
}
type Argst []string

func MakeSSHArgs(msg, c string, args []string) Command {
	return Command{Msg: msg, Type: SSH, Data: c, Args: args}
}

func MakeSSH(msg, c string) Command {
	return Command{Msg: msg, Type: SSH, Data: c}
}

func MakeRsync(msg string, args []string) Command {
	return Command{Msg: msg, Type: RSYNC, Args: args}
}

func MakeLocal(msg, c string) Command {
	return Command{Msg: msg, Type: LOCAL, Data: c}
}

func MakeLocalArgs(msg, c string, args []string) Command {
	return Command{Msg: msg, Type: LOCAL, Data: c, Args: args}
}

func MakeInclude(scriptKey string) Command {
	return Command{Type: INCLUDE, Data: scriptKey}
}

func MakeIncludes(includes ...string) []Command {
	commands := make([]Command, 0)
	for _, i := range includes {
		commands = append(commands, MakeInclude(i))
	}

	return commands
}

func (c Command) GetResultingCmdWithArgs() string {
	if c.Args == nil {
		return c.Data
	}

	argv := []interface{}{}

	for i := range c.Args {
		if c.Args[i] != "defer" {
			argv = append(argv, c.Args[i])
		}
	}

	// FIX cannot use c.Args (variable of type []string) as []any value in argument to fmt.Sprintf
	// cmd = fmt.Sprintf(cmd, c.Args...)
	return fmt.Sprintf(c.Data, argv...)
}

func (c Command) IsDeferred() bool {
	for _, arg := range c.Args {
		if arg == "defer" {
			return true
		}
	}

	return false
}
