package execute

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func Local(cmd string) ([]byte, error) {
	logrus.Debugf("[local] %s ", cmd)
	p := exec.Command("bash", "-c", cmd)

	return p.Output()
}
