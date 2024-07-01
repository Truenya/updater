package defaults

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

func TestIsServiceSupported(t *testing.T) {
	t.Parallel()
	t.Run("Test is service supported", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, true, IsServiceSupported("back", util.UpdateLocal))
		require.Equal(t, true, IsServiceSupported("monitor", util.UpdateLocal))
		require.Equal(t, true, IsServiceSupported("resowatch", util.UpdateLocal))
		require.Equal(t, true, IsServiceSupported("ifacewatch", util.UpdateLocal))
		require.Equal(t, true, IsServiceSupported("uploader", util.UpdateLocal))
		require.Equal(t, true, IsServiceSupported("checker", util.UpdateLocal))
		require.Equal(t, true, IsServiceSupported("back", util.UpdateRemote))
		require.Equal(t, true, IsServiceSupported("monitor", util.UpdateRemote))
		require.Equal(t, true, IsServiceSupported("resowatch", util.UpdateRemote))
		require.Equal(t, true, IsServiceSupported("ifacewatch", util.UpdateRemote))
		require.Equal(t, true, IsServiceSupported("checker", util.UpdateRemote))
		require.Equal(t, false, IsServiceSupported("unknown", util.UpdateRemote))
		require.Equal(t, false, IsServiceSupported("unknown", util.UpdateInPlace))
		require.Equal(t, true, IsServiceSupported("back", util.UpdateInPlace))
		Supported["test"] = ALL
		require.Equal(t, true, IsServiceSupported("test", util.UpdateLocal))
		require.Equal(t, true, IsServiceSupported("test", util.UpdateRemote))
		require.Equal(t, true, IsServiceSupported("test", util.UpdateInPlace))
		Supported["test"] = LOCALREMOTE
		require.Equal(t, true, IsServiceSupported("test", util.UpdateLocal))
		require.Equal(t, true, IsServiceSupported("test", util.UpdateRemote))
		require.Equal(t, false, IsServiceSupported("test", util.UpdateInPlace))
		Supported["test"] = LOCALINPLACE
		require.Equal(t, true, IsServiceSupported("test", util.UpdateLocal))
		require.Equal(t, false, IsServiceSupported("test", util.UpdateRemote))
		require.Equal(t, true, IsServiceSupported("test", util.UpdateInPlace))
		Supported["test"] = REMOTEINPLACE
		require.Equal(t, false, IsServiceSupported("test", util.UpdateLocal))
		require.Equal(t, true, IsServiceSupported("test", util.UpdateRemote))
		require.Equal(t, true, IsServiceSupported("test", util.UpdateInPlace))
		// unknown test
		require.Equal(t, false, IsServiceSupported("UNKNOWN", util.UpdateLocal))
		require.Equal(t, false, IsServiceSupported("UNKNOWN", util.UpdateRemote))
		require.Equal(t, false, IsServiceSupported("UNKNOWN", util.UpdateInPlace))
		Supported["test"] = NONE
		require.Equal(t, false, IsServiceSupported("test", util.UpdateLocal))
		require.Equal(t, false, IsServiceSupported("test", util.UpdateRemote))
		require.Equal(t, false, IsServiceSupported("test", util.UpdateInPlace))
		Supported["test"] = REMOTE
		require.Equal(t, false, IsServiceSupported("test", util.UpdateLocal))
		require.Equal(t, true, IsServiceSupported("test", util.UpdateRemote))
		require.Equal(t, false, IsServiceSupported("test", util.UpdateInPlace))
		Supported["test"] = LOCAL
		require.Equal(t, true, IsServiceSupported("test", util.UpdateLocal))
		require.Equal(t, false, IsServiceSupported("test", util.UpdateRemote))
		require.Equal(t, false, IsServiceSupported("test", util.UpdateInPlace))
		Supported["test"] = INPLACE
		require.Equal(t, false, IsServiceSupported("test", util.UpdateLocal))
		require.Equal(t, false, IsServiceSupported("test", util.UpdateRemote))
		require.Equal(t, true, IsServiceSupported("test", util.UpdateInPlace))
	})
}

func TestGetScriptName(t *testing.T) {
	t.Parallel()
	t.Run("Test get script name", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "back_local", GetScriptName("back", util.UpdateLocal))
		require.Equal(t, "monitor_local", GetScriptName("monitor", util.UpdateLocal))
		require.Equal(t, "resowatch_local", GetScriptName("resowatch", util.UpdateLocal))
		require.Equal(t, "ifacewatch_local", GetScriptName("ifacewatch", util.UpdateLocal))
		require.Equal(t, "uploader_local", GetScriptName("uploader", util.UpdateLocal))
		require.Equal(t, "checker_local", GetScriptName("checker", util.UpdateLocal))
		require.Equal(t, "back_remote", GetScriptName("back", util.UpdateRemote))
		require.Equal(t, "monitor_remote", GetScriptName("monitor", util.UpdateRemote))
		require.Equal(t, "resowatch_remote", GetScriptName("resowatch", util.UpdateRemote))
		require.Equal(t, "ifacewatch_remote", GetScriptName("ifacewatch", util.UpdateRemote))
		require.Equal(t, "checker_remote", GetScriptName("checker", util.UpdateRemote))
		require.Equal(t, "back_inplace", GetScriptName("back", util.UpdateInPlace))
	})
}

func TestError(t *testing.T) {
	t.Parallel()
	t.Run("Test error", func(t *testing.T) {
		t.Parallel()
		n := NotSupportedError{}
		require.Equal(t, n.Error(), NotSupportedMsg)
	})
}

func TestScriptForService(t *testing.T) {
	t.Parallel()
	SetDefaultIfNotExists("back_local", defaultFuncsByTypeAndService[util.UpdateLocal]["back"])
	t.Run("Test script for service", func(t *testing.T) {
		t.Parallel()
		back := BackLocal()
		back.Name = "back_local"
		require.Equal(t, back, ScriptForService("back", util.UpdateLocal))
	})
}

func TestFillSupported(t *testing.T) {
	config.ReadDefaultFile()
	t.Parallel()
	t.Run("Test fill supported", func(t *testing.T) {
		t.Parallel()
		FillSupported()
		for k, v := range Supported {
			res, ok := config.Get(k, "supported")
			require.Equal(t, ok, true)
			require.Equal(t, fmt.Sprint(v), res)
		}
	})
}
