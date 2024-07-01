package execute

import (
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func BySSH(r Remoter, args ...string) (string, error) {
	l, needSshpass, err := r.SSHCommand()
	if err != nil {
		return "", err
	}

	p := GetExecCmd(l, needSshpass, args...)
	out, err := p.Output()

	return string(out), err
}

func AddStrictHostOpts(args []string) []string {
	return append(args, "-o", "StrictHostKeyChecking=no")
}

func GetExecCmd(l []string, needSSHPass bool, args ...string) *exec.Cmd {
	if needSSHPass {
		return WithSSHPass(l, args...)
	}

	return JustSSH(l, args...)
}

func AddBatchModeOpts(args []string) []string {
	return append(args, "-o", "BatchMode=yes")
}

func WithSSHPass(l []string, args ...string) *exec.Cmd {
	l = append(l, args...)
	logrus.Debugln("sshpass", strings.Join(l, " "))

	return exec.Command("sshpass", l...)
}

func JustSSH(l []string, args ...string) *exec.Cmd {
	l = AddBatchModeOpts(l)
	l = append(l, args...)
	logrus.Debugln("ssh", strings.Join(l, " "))

	return exec.Command("ssh", l...)
}

func CheckConnection(_ string) bool {
	return true
}
