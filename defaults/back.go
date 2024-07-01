package defaults

import (
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
)

func BackRemote() script.Script {
	SetDefaultIfNotExists("supervisorctl_stuff", SupervisorStuff)

	return IncludesScript(DefaultInclude{"back_prepare_container", BackImageScript}, DefaultInclude{"back_main", BackMain})
}

func BackLocal() script.Script {
	SetDefaultIfNotExists("supervisorctl_stuff", SupervisorStuff)

	return IncludesScript(
		DefaultInclude{"back_prepare_local", BackPrepareLocal},
		DefaultInclude{"back_main", BackMain},
	)
}

const PythonPackages = "Copying python packages"

func BackMain() script.Script {
	s := CmdCreator{"back"}
	cpPythonPackagesCmd := `bash -c 'for PACK in $(ls /tmp/update_dir/back/python_packages);` +
		` do _PACK=/tmp/update_dir/back/python_packages/${PACK}; if [ -d ${_PACK} ];` +
		` then docker cp /tmp/update_dir/back/python_packages/${PACK}/${PACK} ` +
		`vm_box:/usr/local/lib/python3.9/site-packages/; fi; done'`
	chmodScriptsCmd := "docker exec vm_box bash -c 'chmod +x /opt/ispsystem/vm/scripts/*.py " +
		"&& chmod +x /opt/ispsystem/vm/scripts/*/*.py'" // TODO тут вроде попроще через какую-то bash конструкцию можно
		// TODO ./...
	restartCmd := "docker exec vm_box supervisorctl restart vm_reader vm_writer"

	commands := []command.Command{command.MakeSSH("Copying binaries", s.DCpToDefaultBoxDest("bin/", "vm"))}
	commands = append(commands, command.MakeSSH("Copying scripts", s.DCpToDefaultBoxDest("scripts/", "vm")))
	commands = append(commands, command.MakeSSH("Chmoding scripts", chmodScriptsCmd))
	commands = append(commands, command.MakeSSH("Copying api, etc", s.DCpToDefaultBoxDest("etc/", "vm")))
	commands = append(commands, command.MakeSSH(PythonPackages, cpPythonPackagesCmd))
	commands = append(commands, command.MakeInclude("supervisorctl_stuff"))
	commands = append(commands, command.MakeSSH("Restarting", restartCmd))

	return script.Script{
		Ver:      1,
		Commands: commands,
	}
}
func BackPrepareLocal() script.Script {
	s := CmdCreator{"back"}
	commands := make([]command.Command, 0)
	commands = append(commands, command.MakeSSH("Creating dir", s.MkBinDir()))
	commands = append(commands, command.MakeRsync("Uploading vm by rsync", []string{"build", "bin/vm"}))
	commands = append(commands, command.MakeRsync("Uploading vmctl by rsync", []string{"build", "bin/vmctl"}))
	commands = append(commands, command.MakeRsync("Uploading vmdbfixer by rsync", []string{"build", "bin/vmdbfixer"}))
	commands = append(commands, command.MakeRsync("Uploading etc by rsync", []string{"build", "etc/"}))
	commands = append(commands, command.MakeRsync("Uploading scripts by rsync", []string{"", "scripts/"}))
	commands = append(commands, command.MakeRsync("Uploading python_packages by rsync", []string{"", "python_packages/"}))

	return script.Script{
		Ver:      1,
		Args:     map[string]string{"build_dir": "build"},
		Commands: commands,
	}
}

func BackImageScript() script.Script {
	s := CmdCreator{"back"}
	cmds := []command.Command{command.MakeSSH("Removing dir", s.RmDir())}
	cmds = append(cmds, command.MakeSSHArgs("Deleting container", ContainerDeleteCmd, ContainerDeferArgs()))
	cmds = append(cmds, command.MakeSSH("Creating dir", s.MkBinDir()))
	cmds = append(cmds, command.MakeSSHArgs("Pulling image", s.PullImg(), BranchArgs()))
	cmds = append(cmds, command.MakeSSHArgs("Creating container", s.CreateContainer(), ContainerAndBranchArgs()))
	cmds = append(cmds, command.MakeSSHArgs(PythonPackages, s.DCpFrom("/python_packages"), ContainerArgs()))
	cmds = append(cmds, command.MakeSSHArgs("Copying other", s.DCpAllFrom("vm/"), ContainerArgs()))
	cmds = append(cmds, command.MakeSSH("Moving etc", s.MvInService("vm/etc/", ".")))
	cmds = append(cmds, command.MakeSSH("Moving bins", s.MvInService("vm/bin/", ".")))
	cmds = append(cmds, command.MakeSSH("Moving scripts", s.MvInService("vm/scripts/", ".")))

	// TODO Добавить в основной путь
	// Для этого нужно вытаскивать либы из конана.
	// cpLibsCmd := "docker cp ~/update_dir/back/vm/lib vm_box:/opt/ispsystem/vm/"
	// cpLibsCmd := s.DCpToDefaultBoxDest("vm/lib/", "vm")
	cmds = append(cmds, command.MakeSSH("Copying libs", s.DCpToDefaultBoxDest("vm/lib/", "vm")))

	return script.Script{
		Ver:      1,
		Args:     map[string]string{"branch": "master", "container_name": "back_update"},
		Commands: cmds,
	}
}

func BackInPlace() script.Script {
	cpPythonPackagesCmd := `bash -c 'for PACK in $(ls %s/python_packages);` +
		` do _PACK=%s/python_packages/${PACK}; if [ -d ${_PACK} ]; ` +
		`then docker cp %s/python_packages/${PACK}/${PACK} ` +
		`vm_box:/usr/local/lib/python3.9/site-packages/; fi; done'`
	chmodScriptsCmd := "docker exec vm_box bash -c 'chmod +x /opt/ispsystem/vm/scripts/*.py " +
		"&& chmod +x /opt/ispsystem/vm/scripts/*/*.py'" // TODO тут вроде попроще через какую-то bash конструкцию можно
	restartCmd := "docker exec vm_box supervisorctl restart vm_reader vm_writer"

	commands := []command.Command{command.MakeLocalArgs("Copying binaries",
		"docker cp %s/%s/bin/ vm_box:/opt/ispsystem/vm/", []string{"dir", "build"})}
	commands = append(commands, command.MakeLocalArgs("Copying scripts",
		"docker cp %s/%s/ vm_box:/opt/ispsystem/vm/", []string{"dir", "scripts"}))
	commands = append(commands, command.MakeLocal("Chmoding scripts", chmodScriptsCmd))
	commands = append(commands, command.MakeLocalArgs(PythonPackages, cpPythonPackagesCmd, []string{"dir", "dir", "dir"}))
	commands = append(commands, command.MakeLocal("Restarting", restartCmd))

	return script.Script{
		Ver:      1,
		Commands: commands,
	}
}
