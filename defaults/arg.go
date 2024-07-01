package defaults

import "gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"

func BranchArgs() []string             { return []string{util.Branch} }
func ContainerArgs() []string          { return []string{util.Cont} }
func ContainerAndBranchArgs() []string { return []string{util.Cont, util.Branch} }
func ContainerDeferArgs() []string     { return []string{util.Cont, util.Defer} }
func DeferArg() []string               { return []string{util.Defer} }
