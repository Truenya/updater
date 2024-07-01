package updater

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/statistic"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

func setup() {
	config.ReadDefaultFile()
	script.Init()
	statistic.ReadFromDefaultFile()

	if _, ok := config.Get("ssh", "addr"); !ok {
		config.Set("ssh", "addr", "127.0.0.1")
	}
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func TestSetUp(t *testing.T) {
	t.Parallel()

	r := Reporter{}
	require.Error(t, r.setUp(nil))
	err := r.setUp(make(chan util.Progress, 1))
	require.NoError(t, err)
	require.Equal(t, 1, r.curStage)
	require.Equal(t, time.Duration(0), r.lastStageDur)
	require.Equal(t, 0, len(r.progress))
	require.Equal(t, 1, cap(r.progress))
}

func TestReport(t *testing.T) {
	t.Parallel()

	t.Run("Report", func(t *testing.T) {
		t.Parallel()

		r := Reporter{}
		err := r.setUp(make(chan util.Progress, 1))
		require.NoError(t, err)

		r.report("testMessage", "testNoInfoReport", util.UpdateRemote)
		require.Equal(t, 1, r.curStage) // not changed
		c := <-r.progress
		require.Equal(t, "testMessage", c.Message)
	})
}

func TestUpdateMeanStage(t *testing.T) {
	t.Parallel()

	statistic.ReadFromDefaultFile()

	t.Run("Duration is last stage is updating correctly", func(t *testing.T) {
		t.Parallel()

		r := Reporter{}
		err := r.setUp(make(chan util.Progress, 1))
		require.NoError(t, err)

		statistic.Set("testDuration", "stage_1_ms", time.Second)
		r.updateMeanStage(time.Second, "testDuration", util.UpdateLocal, "test message")
		require.Equal(t, time.Second, r.lastStageDur)
		require.Equal(t, 2, r.curStage)
	})

	t.Run("Statistic is updated correctly when last stage differs from stored", func(t *testing.T) {
		t.Parallel()

		r := Reporter{}
		err := r.setUp(make(chan util.Progress, 1))
		require.NoError(t, err)

		require.Equal(t, 1, r.curStage)
		statistic.Set("testDuration2", "1_inplace_test", 4*time.Second)
		r.updateMeanStage(2*time.Second, "testDuration2", util.UpdateInPlace, "test")
		require.Equal(t, 2, r.curStage)
		require.Equal(t, 2*time.Second, r.lastStageDur)

		v, ok := statistic.Get("testDuration2", "1_inplace_test")
		require.True(t, ok)
		require.Equal(t, 3*time.Second, v)
	})

	t.Run("Statistic is updated correctly when last stage is same as stored", func(t *testing.T) {
		t.Parallel()

		r := Reporter{}
		err := r.setUp(make(chan util.Progress, 1))
		require.NoError(t, err)

		statistic.Set("test", "container_stage_1_ms", 2*time.Second)
		r.updateMeanStage(2*time.Second, "test", util.UpdateRemote, "test message")

		v, ok := statistic.Get("test", "container_stage_1_ms")
		require.True(t, ok)
		require.Equal(t, 2*time.Second, v)
	})
}
