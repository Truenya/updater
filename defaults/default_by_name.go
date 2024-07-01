package defaults

import (
	"fmt"

	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

type SupportedScripts byte

const (
	NONE SupportedScripts = iota
	LOCAL
	REMOTE
	LOCALREMOTE
	INPLACE
	LOCALINPLACE
	REMOTEINPLACE
	ALL
)

var String2SS = map[string]SupportedScripts{ //nolint:gochecknoglobals
	"0": NONE,
	"1": LOCAL,
	"2": REMOTE,
	"3": LOCALREMOTE,
	"4": INPLACE,
	"5": LOCALINPLACE,
	"6": REMOTEINPLACE,
	"7": ALL,
}

type UpdaterConfigs map[string]SupportedScripts

var Supported = UpdaterConfigs{ //nolint:gochecknoglobals
	"back":       ALL,
	"monitor":    ALL,
	"resowatch":  LOCALREMOTE,
	"ifacewatch": LOCALREMOTE,
	"uploader":   LOCAL,
	"checker":    LOCALREMOTE,
	"ha-agent":   INPLACE,
	"node-alert": INPLACE,
}

var defaultFuncsByTypeAndService = map[util.UpdaterType]map[string]func() script.Script{ //nolint:gochecknoglobals
	util.UpdateLocal: {
		"back":       BackLocal,
		"monitor":    MonitorLocal,
		"resowatch":  ResowatchLocal,
		"ifacewatch": IfacewatchLocal,
		"uploader":   UploaderLocal,
		"checker":    CheckerLocal,
	},
	util.UpdateRemote: {
		"back":       BackRemote,
		"monitor":    MonitorRemote,
		"resowatch":  ResowatchRemote,
		"ifacewatch": IfacewatchRemote,
		"checker":    CheckerRemote,
	},
	util.UpdateInPlace: {
		"back":     BackInPlace,
		"ha-agent": HaAgentInPlace,
		// "node-alert":  NodeAlertInPlace,
		"monitor": MonitorInplace,
	},
}

func ScriptForService(service string, uType util.UpdaterType) script.Script {
	scriptName := GetScriptName(service, uType)

	return SetDefaultIfNotExistsAndReturn(scriptName, defaultFuncsByTypeAndService[uType][service])
}

type NotSupportedError struct{}

const NotSupportedMsg = "Not supported"

func (NotSupportedError) Error() string {
	return NotSupportedMsg
}

var uType2String = map[util.UpdaterType]string{ //nolint:gochecknoglobals
	util.UpdateLocal:   "_local",
	util.UpdateRemote:  "_remote",
	util.UpdateInPlace: "_inplace",
}

func GetScriptName(service string, uType util.UpdaterType) string {
	return service + uType2String[uType]
}

func IsServiceSupported(service string, uType util.UpdaterType) bool {
	switch Supported[service] {
	case NONE:
		return false
	case LOCAL:
		return uType == util.UpdateLocal
	case REMOTE:
		return uType == util.UpdateRemote
	case INPLACE:
		return uType == util.UpdateInPlace
	case LOCALREMOTE:
		return uType != util.UpdateInPlace
	case LOCALINPLACE:
		return uType != util.UpdateRemote
	case REMOTEINPLACE:
		return uType != util.UpdateLocal
	case ALL:
		return true
	default:
		return false
	}
}

func FillSupported() {
	for k, v := range Supported {
		config.Set(k, "supported", fmt.Sprint(v))
	}
}
