package statistic

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

var m sync.Mutex                              //nolint:gochecknoglobals
var _data map[string]map[string]time.Duration //nolint:gochecknoglobals
const fileName = "statistic.json"             //nolint:gochecknoglobals

func ReadFromDefaultFile() {
	dirPath, filePath, err := util.DefaultFilePath(fileName)
	if err != nil {
		logrus.Errorln("Failed to open user specific config dir:", err)
	}

	f, err := util.InitCustomJSON(dirPath, filePath)
	if err != nil {
		logrus.Panicf("Failed to open statistic for init %s", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&_data)
	if err != nil {
		logrus.Panicf("Failed to read statistic for init %s", err)
	}
}

func Get(s, k string) (time.Duration, bool) {
	m.Lock()
	defer m.Unlock()

	v, ok := _data[s][k]

	return v, ok
}

func Set(s, k string, v time.Duration) {
	m.Lock()
	defer m.Unlock()

	if _data == nil {
		_data = map[string]map[string]time.Duration{}
	}

	if _, ok := _data[s]; !ok {
		_data[s] = make(map[string]time.Duration)
	}

	_data[s][k] = v
}

func GetEstDuration(service string, stage int, local util.UpdaterType) time.Duration {
	var sum time.Duration

	for k, v := range _data[service] {
		splitted := strings.Split(k, "_")
		if local == util.UpdateLocal && len(splitted) > 3 {
			continue
		}

		if local == util.UpdateRemote && len(splitted) < 4 {
			continue
		}

		if IsNextStage(splitted[0], stage) {
			sum += v
		}
	}

	return sum
}

func GetSumDuration(service string, uType util.UpdaterType) time.Duration {
	var sum time.Duration

	for k, v := range _data[service] {
		splitted := strings.Split(k, "_")
		if uType == util.UpdateLocal && len(splitted) > 3 {
			continue
		}

		if uType == util.UpdateRemote && len(splitted) < 4 {
			continue
		}

		sum += v
	}

	return sum
}

func IsNextStage(k string, curStage int) bool {
	if !util.ContainNumber(k) {
		return false
	}

	stage, err := strconv.Atoi(k)
	if err != nil {
		logrus.Error(err)

		return false
	}

	return curStage <= stage
}

func Write() error {
	logrus.Info("Write statistic")

	_, filePath, err := util.DefaultFilePath(fileName)
	if err != nil {
		logrus.Errorln("Failed to open user specific statistic dir:", err)

		return err
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		logrus.Error(err)

		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")

	if err := e.Encode(_data); err != nil {
		logrus.Warnf("[statistic] Failed to encode to file, %s", err)

		return err
	}

	return nil
}
