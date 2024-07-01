package execute

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func Rsync(r Remoter, localDir string, file string, serviceName string) error {
	uploadBinCmd, err := r.RsyncCommand(localDir, file, serviceName)
	if err != nil {
		return err
	}

	logrus.Debugln("[rsync]: ", uploadBinCmd)
	p := exec.Command("bash", "-c", uploadBinCmd)
	out, err := p.Output()

	if err != nil {
		return fmt.Errorf("%s, %w", string(out), err)
	}

	return err
}
