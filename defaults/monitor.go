package defaults

import (
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
)

func MonitorLocal() script.Script {
	SetDefaultIfNotExists("supervisorctl_stuff", SupervisorStuff)

	return IncludesScript(
		DefaultInclude{"monitor_prepare_local", MonitorPrepareLocal},
		DefaultInclude{"monitor_main", MonitorMain},
	)
}

func MonitorRemote() script.Script {
	SetDefaultIfNotExists("supervisorctl_stuff", SupervisorStuff)

	return IncludesScript(
		DefaultInclude{"monitor_prepare_container", MonitorPrepareContainer},
		DefaultInclude{"monitor_main", MonitorMain},
	)
}

func MonitorMain() script.Script {
	m := CommandGenerator{CmdCreator{"monitor"}}

	return GenerateScript(m.CpBinToBox, m.CpEtcToBox, m.CpSupervisordToBox, IncludeSuper, m.Restart)
}

func MonitorPrepareLocal() script.Script {
	m := CommandGenerator{CmdCreator{"monitor"}}

	return GenerateScript(m.MkBinDir, m.RsyncBin, m.MvBin, RsyncEtc, RsyncSupervisord)
}

func MonitorPrepareContainer() script.Script {
	m := CommandGenerator{CmdCreator{"monitor"}}

	return GenerateScript(m.RmDir, RmCont, m.MkBinDir, m.PullImg, m.CreateCont, m.CpServiceFrom,
		m.CpSupervisordFrom).AddDefaultRemoteArgs(m.cc.s)
}

func MonitorInplace() script.Script {
	return script.Script{
		Commands: []command.Command{
			command.MakeLocalArgs("Pulling containter",
				"docker pull registry-dev.ispsystem.net/team/vm/monitor:%s", BranchArgs()),
			command.MakeLocalArgs("Creating container",
				"docker create --name %s registry-dev.ispsystem.net/team/vm/monitor:%s sh",
				ContainerAndBranchArgs()),
			command.MakeLocalArgs("Removing container", "docker rm %s", ContainerDeferArgs()),
			command.MakeLocalArgs("Copying data", "docker cp %s:/bin/monitor .", ContainerArgs()),
			command.MakeLocal("Md5local creating", "md5sum ./monitor > md5local.txt"),
			command.MakeLocalArgs("Removing local", "rm -rf monitor", DeferArg()),
			command.MakeLocal("Copying to container", "docker cp ./monitor vm_box:/root/"),
			command.MakeLocal("Md5container creating", "docker exec vm_box md5sum /root/monitor > md5cont.txt"),
			command.MakeLocalArgs("Removing in container", "docker exec vm_box rm -rf /root/monitor",
				DeferArg()),
		},
	}
}
