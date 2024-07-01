package defaults

import (
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
)

func ResowatchLocal() script.Script {
	return IncludesScript(
		DefaultInclude{"resowatch_prepare_local", ResowatchPrepareLocal},
		DefaultInclude{"resowatch_main", ResowatchMain},
	)
}

func ResowatchRemote() script.Script {
	return IncludesScript(
		DefaultInclude{"resowatch_prepare_container", ResowatchPrepareContainer},
		DefaultInclude{"resowatch_main", ResowatchMain},
	)
}

func ResowatchPrepareContainer() script.Script {
	cg := CommandGenerator{CmdCreator{"resowatch"}}

	return GenerateScript(RmCont, cg.RmDir, cg.MkBinDir, cg.PullImg,
		cg.CreateCont, cg.CpServiceFrom).AddDefaultRemoteArgs(cg.cc.s)
}

func ResowatchPrepareLocal() script.Script {
	cg := CommandGenerator{CmdCreator{"resowatch"}}

	return GenerateScript(cg.MkBinDir, cg.RsyncBin, cg.MvBin, RsyncSupervisord, RsyncEtc)
}

func ResowatchMain() script.Script {
	cg := CommandGenerator{CmdCreator{"resowatch"}}

	return GenerateScript(cg.CpBinToBox, cg.Restart)
}
