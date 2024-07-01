package execute

import (
	"errors"
	"fmt"

	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/defaults"
)

type Remoter struct {
	Key  string
	Pass string
	Addr string
	User string
}

func RemoterFromConfig() Remoter {
	key, _ := config.Get("ssh", "key")
	pass, _ := config.Get("ssh", "pass")
	addr, _ := config.Get("ssh", "addr")
	user, _ := config.Get("ssh", "user")

	return Remoter{
		Key:  key,
		Pass: pass,
		Addr: addr,
		User: user,
	}
}

func RemoterFromData(sshData map[string]string) Remoter {
	return Remoter{
		Key:  sshData["key"],
		Pass: sshData["pass"],
		Addr: sshData["addr"],
		User: sshData["user"],
	}
}

func (r *Remoter) AddKey(l []string) []string {
	if r.Key != "" {
		l = append(l, "-i", r.Key)
	}

	return l
}

func (r *Remoter) GetUserAndAddr() (string, error) {
	if r.Addr == "" {
		return "", errors.New("address is not specified")
	}

	if r.User == "" {
		return r.Addr, nil
	}

	return r.User + "@" + r.Addr, nil
}

func (r *Remoter) AddPassword(args []string) ([]string, bool) {
	if r.Pass != "" {
		return append(args, "-p", r.Pass, "ssh"), true
	}

	return args, false
}

func (r *Remoter) GetRsyncSSH() string {
	if r.Key != "" {
		return `-e "ssh -i ` + r.Key + `"`
	}

	if r.Pass != "" {
		return `-e "sshpass -p ` + r.Pass + ` ssh"`
	}

	return ""
}

func (r *Remoter) RsyncCommand(localDir string, file string, serviceName string) (string, error) {
	// TODO move most part command creating to defaults
	rsyncPart := fmt.Sprintf("rsync %s -r ", r.GetRsyncSSH())
	localFilePath := fmt.Sprintf(" %s/%s ", localDir, file)
	addr, err := r.GetUserAndAddr()
	remoteFilePath := fmt.Sprintf(" %s:%s%s/%s", addr, defaults.DefaultRemoteDir, serviceName, file)

	return rsyncPart + localFilePath + remoteFilePath, err
}

func (r *Remoter) SSHCommand() ([]string, bool, error) {
	l, needSshpass := r.AddPassword([]string{})
	l = r.AddKey(l)
	l = AddStrictHostOpts(l)

	ua, err := r.GetUserAndAddr()

	return append(l, ua), needSshpass, err
}
