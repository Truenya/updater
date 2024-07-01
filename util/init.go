package util

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func ProcessWithWarnLog(fn func() error) {
	if err := fn(); err != nil {
		logrus.Warnln(err)
	}
}

func LogNotNilErr(err error) {
	if err != nil {
		logrus.Errorln(err)
	}
}

func InitCustomJSON(dirPath, filePath string) (*os.File, error) {
	// ~/.config/updater/config.json
	_, err := os.Stat(filePath)
	if err != nil {
		os.MkdirAll(dirPath, 0755) //nolint:errcheck

		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			logrus.Panicf("Failed to create default config: %s", err)
		}

		ProcessWithWarnLog(func() error {
			_, err := f.WriteString("{}")

			return fmt.Errorf("failed to write default config: %w", err)
		})
		f.Close()
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}

	return f, err
}

type DefaultFilePathError struct {
	osErr error
}

func (dfe DefaultFilePathError) Error() string {
	return fmt.Sprintf("Failed to get config directory: %s", dfe.osErr)
}

func DefaultFilePath(fileName string) (string, string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		// ./.config.json
		return "./", "./." + fileName, DefaultFilePathError{err}
	}

	dirPath := config + "/updater/"
	// ~/.config/updater/
	// ./updater $PWD./updater/config.json DefaultFilePathErr
	// ~.config/updater ~/.config/updater/config.json nil
	return dirPath, dirPath + fileName, err
}

func WarnNotNilErr(err error) {
	if err != nil {
		logrus.Warnln(err)
	}
}
