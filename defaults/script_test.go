package defaults

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
)

func setup() {
	script.Init()
}

func teardown() {
	path, err := script.GetUserPath()
	if err == nil {
		logrus.Traceln("removing", path, "/updater/scripts/test")
		os.RemoveAll(path + "/updater/scripts/test")
	} else {
		logrus.Errorln("removing", "./scripts/", err)
		os.RemoveAll("./scripts/")
	}
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func TestIncludes(t *testing.T) {
	t.Parallel()

	require.Equal(t,
		IncludesScript(DefaultInclude{"back_prepare_container", BackImageScript}, DefaultInclude{"back_main", BackMain}),
		ScriptForTest("back/remote.json"))
}

func TestSetDefault(t *testing.T) {
	t.Parallel()

	goldName := "test_back_prepare_local2"
	script := SetDefaultIfNotExistsAndReturn(goldName, BackPrepareLocal)
	gold := BackPrepareLocal()
	require.Equal(t, gold, script)
}

func TestSetNotExistingDefault(t *testing.T) {
	t.Parallel()

	goldName := "test_not_exist_supervisorctl_stuff"
	_, ok := script.Get(goldName)
	require.False(t, ok)

	goldf := SupervisorStuff
	SetDefaultIfNotExists(goldName, goldf)
	result, ok := script.Get(goldName)
	require.True(t, ok)

	gold := goldf()
	require.NotEqual(t, gold, result)
	gold.Name = goldName
	require.Equal(t, gold, result)
	result = SetDefaultIfNotExistsAndReturn(goldName+"2", SupervisorStuff)
	gold.Name = ""
	require.Equal(t, gold, result)
}

func TestSetNotExistingDefaultAndReturn(t *testing.T) {
	t.Parallel()

	goldName := "test_not_exist_ss_stuff"
	_, ok := script.Get(goldName)
	require.False(t, ok)

	goldf := SupervisorStuff
	result := SetDefaultIfNotExistsAndReturn(goldName, goldf)
	gold := goldf()
	require.Equal(t, gold, result)
}
