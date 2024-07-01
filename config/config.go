package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

// service_name to service settings.
var data map[string]map[string]string //nolint:gochecknoglobals
const fileName = "config.json"        //nolint:gochecknoglobals
var m sync.Mutex                      //nolint:gochecknoglobals

func ReadDefaultFile() {
	dirPath, filePath, err := util.DefaultFilePath(fileName)
	if err != nil {
		log.Errorln("Failed to open user specific config dir:", err)
	}

	f, err := util.InitCustomJSON(dirPath, filePath)
	if err != nil {
		log.Panicf("Failed to open config for init %s", err)
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&data); err != nil {
		log.Panicf("Failed to read config for init %s", err)
	}
}

func Write() error {
	_, filePath, err := util.DefaultFilePath(fileName)
	if err != nil {
		return fmt.Errorf("failed to open user specific config dir: %w", err)
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to open config: %w", err)
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")

	if err = e.Encode(data); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func Empty() bool {
	return len(GetSupported()) == 0
}

func GetSelects(k string) ([]string, bool) {
	s, ok := Get("select", k)
	if !ok {
		return []string{}, false
	}

	splitted := strings.Split(s, ";")

	return splitted, ok
}

func GetSelectsFiltered(k, v string) ([]string, bool) {
	s, ok := GetSelects(k)
	if !ok {
		return []string{}, false
	}

	result := make([]string, 0)

	for _, str := range s {
		if strings.Contains(str, v) {
			result = append(result, str)
		}
	}

	return result, true
}
func AddSelect(k, v string) {
	s, ok := Get("select", k)
	if !ok || len(s) == 0 {
		Set("select", k, v)

		return
	}

	splitted := strings.Split(s, ";")
	result := splitted[0]

	for i, str := range splitted {
		if str == v {
			// Dont add same.
			return
		}

		if i == 0 {
			// First already added
			continue
		}

		result = fmt.Sprintf("%s;%s", result, str)
	}

	result = fmt.Sprintf("%s;%s", result, v)
	Set("select", k, result)
}

func Set(k, k2, v string) {
	m.Lock()
	defer m.Unlock()

	if _, ok := data[k]; !ok {
		data[k] = make(map[string]string)
	}

	data[k][k2] = v
}

func Unset(k, k2 string) {
	m.Lock()
	defer m.Unlock()
	delete(data[k], k2)
}

func Get(k, k2 string) (string, bool) {
	m.Lock()
	defer m.Unlock()

	v, ok := data[k][k2]

	return v, ok
}

func Args(service string) map[string]string {
	return data[service]
}

func ClearDirs() {
	for k := range data {
		Unset(k, "dir")
	}
}

func GetSupported() map[string]string {
	output := make(map[string]string)

	for k, v := range data {
		if supported, ok := v["supported"]; ok {
			output[k] = supported
		}
	}

	return output
}

func GetServicesFiltered(v string) []string {
	services := make([]string, 0)

	for k := range data {
		if strings.Contains(k, v) && k != v {
			services = append(services, k)
		}
	}

	return services
}
