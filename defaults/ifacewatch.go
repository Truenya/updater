package defaults

import (
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
)

//
// func IfacewatchAll() script.Script {
// 	s := CommandGenerator{CmdCreator{"ifacewatch"}}
// 	return GenerateScript(s.MkBinDir, s.RsyncBin, s.MvBin, s.CpBinToBox, s.Restart)
// }

func IfacewatchLocal() script.Script {
	return IncludesScript(
		DefaultInclude{"ifacewatch_prepare_local", IfacewatchPrepareLocal},
		DefaultInclude{"ifacewatch_main", IfacewatchMain},
	)
}

func IfacewatchRemote() script.Script {
	return IncludesScript(
		DefaultInclude{"ifacewatch_prepare_container", IfacewatchPrepareRemote},
		DefaultInclude{"ifacewatch_main", IfacewatchMain},
	)
}

func IfacewatchPrepareLocal() script.Script {
	s := CommandGenerator{CmdCreator{"ifacewatch"}}

	return GenerateScript(s.MkBinDir, s.RsyncBin, s.MvBin)
}

func IfacewatchPrepareRemote() script.Script {
	cg := CommandGenerator{CmdCreator{"ifacewatch"}}
	f := func() command.Command {
		return command.MakeSSHArgs("Copying bin", cg.cc.DCpFrom("/bin"), []string{"container_name"})
	}

	return GenerateScript(RmCont, cg.RmDir, cg.PullImg, cg.CreateCont, f).AddDefaultRemoteArgs(cg.cc.s)
}

func IfacewatchMain() script.Script {
	s := CommandGenerator{CmdCreator{"ifacewatch"}}

	return GenerateScript(s.CpBinToBox, s.Restart)
}
