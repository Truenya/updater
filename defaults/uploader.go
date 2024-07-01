package defaults

import (
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
)

func UploaderLocal() script.Script {
	// s := CommandGenerator{CmdCreator{"uploader"}}
	// s := CmdCreator{"uploader"}
	//    mkdir := func (){
	//        s.cc.MkDir()
	//    }
	// return GenerateScript(s.MkSupDir,s.MkBinDir,s.RsyncBin,   )
	s := CmdCreator{"uploader"}
	// dirCreateCmd := "mkdir  -p  ~/update_dir/uploader/supervisord.conf.d"
	// mkdirUploaderCmd := "docker exec vm_box mkdir -p /opt/ispsystem/uploader/bin/"
	// cpUploaderCmd := "docker cp /tmp/update_dir/uploader/uploader vm_box:/opt/ispsystem/uploader/bin/uploader"
	// cpEtcCmd := "docker cp /tmp/update_dir/uploader/etc vm_box:/opt/ispsystem/uploader/"
	// cpSupervisorCmd := "docker cp /tmp/update_dir/uploader/supervisord.conf.d vm_box:/etc"
	// stopCmd := "docker exec vm_box supervisorctl stop uploader"
	// startCmd := "docker exec vm_box supervisorctl start uploader"

	commands := []command.Command{command.MakeSSH("Creating dir", s.MkDir("supervisord.conf.d"))}
	commands = append(commands, command.MakeRsync("Uploading uploader by rsync", []string{"", "uploader"}))
	commands = append(commands, command.MakeRsync("Uploading supervisor config by rsync",
		[]string{"", "supervisord.conf.d/uploader.conf"}))
	commands = append(commands, command.MakeRsync("Uploading etc by rsync", []string{"", "etc"}))
	// commands = append(commands, command.MakeSSH("Creating uploader dir in container", mkdirUploaderCmd))
	// commands = append(commands, command.MakeSSH("Copying uploader", cpUploaderCmd))
	// commands = append(commands, command.MakeSSH("Copying etc", cpEtcCmd))
	// commands = append(commands, command.MakeSSH("Copying supervisor config", cpSupervisorCmd))
	commands = append(commands, command.MakeInclude("supervisorctl_stuff"))
	// commands = append(commands, command.MakeSSH("Stopping uploader", stopCmd))
	// commands = append(commands, command.MakeSSH("Starting uploader", startCmd))
	return script.Script{
		Ver:      1,
		Commands: commands,
	}
}
