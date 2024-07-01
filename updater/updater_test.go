package updater

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

func TestGet(t *testing.T) {
	t.Parallel()

	u := Get("testl", util.UpdateLocal, map[string]string{"test": "test"}, map[string]string{"addr": "localhost"})
	require.NotNil(t, u)
	require.Equal(t, util.UpdateLocal, u.uType)
	require.Equal(t, "testl", u.service)
}

func TestGetKeys(t *testing.T) {
	t.Parallel()

	keys := getKeys(map[string]struct{}{"test": {}})
	require.Equal(t, 1, len(keys))
}

func TestBusy(t *testing.T) {
	t.Parallel()

	u := Get("test", util.UpdateLocal, map[string]string{"test": "test"}, map[string]string{"addr": "localhost"})
	require.False(t, u.IsBusy())
	u.SetBusy(true)
	require.True(t, u.IsBusy())
}

func TestUploadByRsync(t *testing.T) {
	t.Parallel()

	t.Run("When dir is not set it is error", func(t *testing.T) {
		t.Parallel()

		config.Unset("test", "dir")

		u := Get("test", util.UpdateLocal, map[string]string{"test": "test"}, map[string]string{"addr": "localhost"})
		err := u.uploadByRsync("test", "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get dir for service: test")
	})
}

func TestExecuteLocally(t *testing.T) {
	t.Parallel()

	u := Get("test", util.UpdateInPlace, map[string]string{"test": "test"}, map[string]string{})
	err := u.executeLocally("ls")
	require.NoError(t, err)
}

//
// func TestUpdate(t *testing.T) {
// 	t.Parallel()
//
// 	script.Init()
// 	config.Set("test", "supported", "3")
// 	defaults.SetDefaultIfNotExistsAndReturn("test_inplace", defaults.CheckLSExists)
//
// 	u := Get("test", util.UpdateLocal, map[string]string{"test": "test"}, map[string]string{})
// 	err := u.Update(make(chan string, 1))
// 	require.Error(t, err)
// 	require.Contains(t, err.Error(), "failed to get dir for service: test")
//
// 	u2 := Get("test", util.UpdateInPlace, map[string]string{"test": "test"}, map[string]string{})
// 	err = u2.Update(make(chan string, 1))
// 	require.NoError(t, err)
// }

func TestSetArgs(t *testing.T) {
	t.Parallel()

	u := Get("test", util.UpdateLocal, map[string]string{}, map[string]string{})
	u.SetArgs(map[string]string{"test": "test"})
	require.Equal(t, "test", u.args["test"])
}

// func TestSsh(t *testing.T) {
// 	t.Parallel()
//
// 	oldAddr, _ := config.Get("ssh", "addr")
// 	defer func(oldAddr string) {
// 		config.Set("ssh", "addr", oldAddr)
// 	}(oldAddr)
//
// 	config.Set("ssh", "addr", "localhost")
// 	u := Get("test", util.UpdateLocal)
// 	err := u.r.setUp(make(chan string, 1))
// 	require.NoError(t, err)
// 	err = u.processWithMeanSSH(command.Command{})
// 	require.NoError(t, err)
// }
//
// func TestMean(t *testing.T) {
// 	t.Parallel()
//
// 	u := Get("test", util.UpdateInPlace)
// 	err := u.r.setUp(make(chan string, 1))
// 	require.NoError(t, err)
// 	err = u.processWithMeanLocal(command.Command{Data: "ls"})
// 	require.NoError(t, err)
// }
//
// func TestRsync(t *testing.T) {
// 	t.Parallel()
//
// 	u := Get("test", util.UpdateRemote)
// 	err := u.r.setUp(make(chan string, 1))
// 	require.NoError(t, err)
// 	err = u.processWithMeanRsync(command.Command{})
// 	require.Error(t, err)
// }
