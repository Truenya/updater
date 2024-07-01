package main

import (
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/gui/gui"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/startup"
)

func main() {
	defer startup.RecoverMain()
	startup.Prepare()
	gui.Application()
}
