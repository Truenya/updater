package statistic

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

func setup() {
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func TestStatistic(t *testing.T) {
	t.Parallel()

	Set("test", "1_", time.Second)

	res, ok := Get("test", "1_")
	require.True(t, ok)
	require.Equal(t, time.Second, res)
	Set("test", "2_", 2*time.Second)

	res, ok = Get("test", "2_")
	require.True(t, ok)
	require.Equal(t, 2*time.Second, res)
	require.Equal(t, 2*time.Second, GetEstDuration("test", 2, util.UpdateLocal))
	require.Equal(t, 3*time.Second, GetEstDuration("test", 1, util.UpdateLocal))
	require.Equal(t, 3*time.Second, GetEstDuration("test", 0, util.UpdateLocal))
}
