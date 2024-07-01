package defaults

import (
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

func HaAgentInPlace() script.Script {
	return script.Script{
		Commands: []command.Command{
			command.MakeLocalArgs("Pulling containter",
				"docker pull registry-dev.ispsystem.net/team/vm/ha-agent:%s", []string{"branch"}),
			command.MakeLocalArgs("Creating container",
				"docker create --name %s registry-dev.ispsystem.net/team/vm/ha-agent:%s sh",
				[]string{"container_name", "branch"}),
			command.MakeLocalArgs("Removing container", "docker rm %s", []string{"container_name", "defer"}),
			command.MakeLocalArgs("Copying data", "docker cp %s:/bin/ha-agent .", []string{"container_name"}),
			command.MakeLocal("Md5local creating", "md5sum ./ha-agent > md5local.txt"),
			command.MakeLocalArgs("Removing local", "rm -rf ha-agent", []string{util.Defer}),
			command.MakeLocal("Copying to container", "docker cp ./ha-agent vm_box:/root/"),
			command.MakeLocal("Md5container creating", "docker exec vm_box md5sum /root/ha-agent > md5cont.txt"),
			command.MakeLocalArgs("Removing in container", "docker exec vm_box rm -rf /root/ha-agent",
				[]string{util.Defer}),
			command.MakeLocal("Copying to node7",
				"docker exec vm_box vssh 7 -U --src /root/ha-agent --dst /opt/ispsystem/vm/bin/"),
			command.MakeLocal("Copying to node8",
				"docker exec vm_box vssh 8 -U --src /root/ha-agent --dst /opt/ispsystem/vm/bin/"),
			command.MakeLocal("Copying to node9",
				"docker exec vm_box vssh 9 -U --src /root/ha-agent --dst /opt/ispsystem/vm/bin/"),
			command.MakeLocal("Restarting on node7",
				"docker exec vm_box vssh 7 'systemctl restart ha-agent'"),
			command.MakeLocal("Restarting on node8",
				"docker exec vm_box vssh 8 'systemctl restart ha-agent'"),
			command.MakeLocal("Restarting on node9",
				"docker exec vm_box vssh 9 'systemctl restart ha-agent'"),
		},
		Args: map[string]string{
			util.Cont: "ha_agent_upd",
		},
	}
}

// import (
// 	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
// 	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
// )
//
// func HaAgentRemote() script.Script {
// 	createDirCmd := "mkdir -p ~/update_dir/ha-agent/"
// 	commands := []command.Command{command.MakeSSH("Creating directory", createDirCmd)}
//
// 	rmExistingCmd := "rm -rf ~/update_dir/ha-agent/*"
// 	commands = append(commands, command.MakeSSH("Removing dir", rmExistingCmd))
//
// 	args := []string{"branch"}
// 	pullImageCmd := "docker pull registry-dev.ispsystem.net/team/vm/ha-agent:%s"
// 	commands = append(commands, command.MakeSSHArgs("Pulling containter", pullImageCmd, args))
//
// 	args = []string{"container_name", "branch"}
// 	createContainerCmd := "docker create --name %s registry-dev.ispsystem.net/team/vm/ha-agent:%s sh"
// 	commands = append(commands, command.MakeSSHArgs("Creating container", createContainerCmd, args))
//
// 	args = []string{"container_name"}
// 	copyFromContainerCmd := "docker cp %s:/bin/ha-agent /opt/ispsystem/vm/bin/"
// 	commands = append(commands, command.MakeSSHArgs("Copying data", copyFromContainerCmd, args))
//
// 	containerDeleteCmd := "docker rm %s" // container_name
// 	args = []string{"container_name"}
// 	commands = append(commands, command.MakeSSHArgs("Deleting container", containerDeleteCmd, args))
// 	return script.Script{
// 		Ver:      1,
// 		Args:     map[string]string{"branch": "master", "container_name": "ha-agent_update"},
// 		Commands: commands,
// 	}
// }
