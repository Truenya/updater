package defaults

import (
	"fmt"

	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
)

type CommandGenerator struct {
	cc CmdCreator
}

func (cg CommandGenerator) MkBinDir() command.Command {
	return command.MakeSSH("Creating dir", cg.cc.MkBinDir())
}

func (cg CommandGenerator) RsyncBin() command.Command {
	return command.MakeRsync(fmt.Sprintf("Uploading %s by rsync", cg.cc.s), cg.cc.RsyncBin())
}

func (cg CommandGenerator) MvBin() command.Command {
	return command.MakeSSH("Moving "+cg.cc.s+" to bin dir", cg.cc.MvInService(cg.cc.s, "bin"))
}

func (cg CommandGenerator) RmDir() command.Command {
	return command.MakeSSH("Removing dir", cg.cc.RmDir())
}

func (cg CommandGenerator) PullImg() command.Command {
	return command.MakeSSHArgs("Pulling containter", cg.cc.PullImg(), BranchArgs())
}

func (cg CommandGenerator) CreateCont() command.Command {
	return command.MakeSSHArgs("Creating container", cg.cc.CreateContainer(), ContainerAndBranchArgs())
}

func (cg CommandGenerator) CpServiceFrom() command.Command {
	return command.MakeSSHArgs("Copying data", DockerCpFromCmd("/"+cg.cc.s+"/"), ContainerArgs())
}

func (cg CommandGenerator) CpBinFromRoot() command.Command {
	return command.MakeSSHArgs("Copying bin", cg.cc.DCpBinFrom(cg.cc.s), ContainerArgs())
}

func (cg CommandGenerator) CpSupervisordFrom() command.Command {
	return command.MakeSSHArgs("Copying supervisord", cg.cc.DCpFrom("/etc/supervisord.conf.d/"), ContainerArgs())
}

func (cg CommandGenerator) CpBinToBox() command.Command {
	return command.MakeSSH("Copying bin", cg.cc.DCpBinToBox())
}

func (cg CommandGenerator) CpEtcToBox() command.Command {
	return command.MakeSSH("Copying etc", cg.cc.DCpToDefaultBox("etc/"))
}

func (cg CommandGenerator) CpSupervisordToBox() command.Command {
	return command.MakeSSH("Copying supervisorctl", cg.cc.DCpToBoxDest("supervisord.conf.d/", "/etc/"))
}

func (cg CommandGenerator) Restart() command.Command {
	return command.MakeSSH("Restarting "+cg.cc.s, cg.cc.DRestartBoxService())
}

func RsyncEtc() command.Command {
	return command.MakeRsync("Uploading etc by rsync", RsyncArg("etc/"))
}

func RsyncSupervisord() command.Command {
	return command.MakeRsync("Uploading supervisoctl conf by rsync", RsyncArg("supervisord.conf.d/"))
}

func RmCont() command.Command {
	return command.MakeSSHArgs("Deleting container", ContainerDeleteCmd, ContainerDeferArgs())
}

func IncludeSuper() command.Command {
	return command.MakeInclude("supervisorctl_stuff")
}

func RereadSuper() command.Command {
	return command.MakeSSHArgs("Rereading supervisor configs",
		"docker exec %s supervisorctl reread", []string{"container_name"})
}

func UpdateSuper() command.Command {
	return command.MakeSSHArgs("Update supervisor configs",
		"docker exec %s supervisorctl update", []string{"container_name"})
}

func CheckRsyncLocal() command.Command {
	return command.MakeLocal("Checking rsync exists on local", "which rsync")
}

func CheckRsyncRemote() command.Command {
	return command.MakeSSH("Checking rsync exists on remote", "which rsync")
}

func CheckSSHPass() command.Command {
	return command.MakeLocal("Checking sshpass exists on local", "which sshpass")
}

func CheckLS() command.Command {
	return command.MakeLocal("Checking ls exists on local", "which ls")
}
