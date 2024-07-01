package main

import (
	"os"

	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/cli/cli"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/startup"
)

func main() {
	defer startup.RecoverMain()
	startup.Prepare()
	cli.Run(os.Args)
}
