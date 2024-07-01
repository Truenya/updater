package defaults

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
)

func ScriptForTest(file string) script.Script {
	var script script.Script

	fullPath := "testdata/" + file

	f, err := os.OpenFile(fullPath, os.O_RDONLY, 0755)
	if err != nil {
		logrus.Panicf("Failed to open config %s for init %s", fullPath, err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&script)
	if err != nil {
		logrus.Panicf("Failed to read config %s for init %s", fullPath, err)
	}

	return script
}

func TestMonitor(t *testing.T) {
	t.Parallel()

	require.Equal(t, MonitorMain(), ScriptForTest("monitor/main.json"))
	require.Equal(t, MonitorPrepareLocal(), ScriptForTest("monitor/prepare_local.json"))
	require.Equal(t, MonitorPrepareContainer(), ScriptForTest("monitor/prepare_container.json"))
	require.Equal(t, MonitorLocal(), ScriptForTest("monitor/local.json"))
	require.Equal(t, MonitorRemote(), ScriptForTest("monitor/remote.json"))
}

func TestChecker(t *testing.T) {
	t.Parallel()

	require.Equal(t, CheckerMain(), ScriptForTest("checker/main.json"))
	require.Equal(t, CheckerPrepareLocal(), ScriptForTest("checker/prepare_local.json"))
	require.Equal(t, CheckerPrepareContainer(), ScriptForTest("checker/prepare_container.json"))
	require.Equal(t, CheckerLocal(), ScriptForTest("checker/local.json"))
	require.Equal(t, CheckerRemote(), ScriptForTest("checker/remote.json"))
}

func TestBack(t *testing.T) {
	t.Parallel()

	require.Equal(t, BackMain(), ScriptForTest("back/main.json"))
	require.Equal(t, BackPrepareLocal(), ScriptForTest("back/prepare_local.json"))
	require.Equal(t, BackImageScript(), ScriptForTest("back/prepare_container.json"))
	require.Equal(t, BackLocal(), ScriptForTest("back/local.json"))
	require.Equal(t, BackRemote(), ScriptForTest("back/remote.json"))
}

func TestResowatch(t *testing.T) {
	t.Parallel()

	require.Equal(t, ResowatchMain(), ScriptForTest("resowatch/main.json"))
	require.Equal(t, ResowatchPrepareLocal(), ScriptForTest("resowatch/prepare_local.json"))
	require.Equal(t, ResowatchPrepareContainer(), ScriptForTest("resowatch/prepare_container.json"))
	require.Equal(t, ResowatchLocal(), ScriptForTest("resowatch/local.json"))
	require.Equal(t, ResowatchRemote(), ScriptForTest("resowatch/remote.json"))
}

func TestUploader(t *testing.T) {
	t.Parallel()

	// require.Equal(t, UploaderPrepareLocal(), ScriptForTest("uploader/prepare_local.json"))
	// require.Equal(t, UploaderPrepareContainer(), ScriptForTest("uploader/prepare_container.json"))
	require.Equal(t, UploaderLocal(), ScriptForTest("uploader/all.json"))
}

func TestIfacewatch(t *testing.T) {
	t.Parallel()

	require.Equal(t, IfacewatchMain(), ScriptForTest("ifacewatch/main.json"))
	require.Equal(t, IfacewatchPrepareLocal(), ScriptForTest("ifacewatch/prepare_local.json"))
	require.Equal(t, IfacewatchPrepareRemote(), ScriptForTest("ifacewatch/prepare_container.json"))
	require.Equal(t, IfacewatchLocal(), ScriptForTest("ifacewatch/local.json"))
	require.Equal(t, IfacewatchRemote(), ScriptForTest("ifacewatch/remote.json"))
}

func TestSupervisorStuff(t *testing.T) {
	t.Parallel()

	require.Equal(t, SupervisorStuff(), ScriptForTest("supervisorctl/stuff.json"))
}

func TestUtil(t *testing.T) {
	t.Parallel()

	require.Equal(t, CheckRsyncExists(), ScriptForTest("util/check_rsync_exists.json"))
	require.Equal(t, CheckSSHPassExists(), ScriptForTest("util/check_sshpass_exists.json"))
}
