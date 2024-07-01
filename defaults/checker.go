package defaults

import (
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
)

func CheckerRemote() script.Script {
	SetDefaultIfNotExists("supervisorctl_stuff", SupervisorStuff)

	return IncludesScript(
		DefaultInclude{"checker_prepare_container", CheckerPrepareContainer},
		DefaultInclude{"checker_main", CheckerMain},
	)
}

func CheckerLocal() script.Script {
	return IncludesScript(
		DefaultInclude{"checker_prepare_local", CheckerPrepareLocal},
		DefaultInclude{"checker_main", CheckerMain},
	)
}

func CheckerPrepareLocal() script.Script {
	c := CommandGenerator{CmdCreator{"checker"}}

	return GenerateScript(c.MkBinDir, c.RsyncBin, c.MvBin)
}

func CheckerMain() script.Script {
	c := CommandGenerator{CmdCreator{"checker"}}

	return GenerateScript(c.CpBinToBox, c.Restart)
}

func CheckerPrepareContainer() script.Script {
	c := CommandGenerator{CmdCreator{"checker"}}

	return GenerateScript(RmCont, c.RmDir, c.MkBinDir, c.PullImg, c.CreateCont, c.CpBinFromRoot).AddArgs(
		[]string{"branch", "container_name"}, []string{"master", "checker_update"})
}
