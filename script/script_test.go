package script

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/command"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
)

const fileName string = "test_create_file"

func setup() {
	Init()
	config.ReadDefaultFile()
}

func teardown() {
	os.Remove(fileName)
	config.Unset("test", "branch")
	os.RemoveAll(dirPath + "test/")
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func TestSetGet(t *testing.T) {
	t.Parallel()

	gold := Script{Ver: 1, Commands: []command.Command{command.MakeInclude("test_include")}}
	goldName := "test_set"
	Set(goldName, gold)
	s, ok := Get(goldName)
	gold.Name = goldName

	require.True(t, ok)
	require.Equal(t, gold, s)
}

func TestAddArgs(t *testing.T) {
	t.Parallel()

	gold := Script{Ver: 1, Commands: []command.Command{command.MakeInclude("test_include")}}
	// gold_name := "test_args"
	result := gold.AddArgs([]string{"test"}, []string{"args"})
	require.NotEqual(t, gold, result)
	require.Nil(t, gold.Args)
	gold.Args = map[string]string{}
	gold.Args["test"] = "args"
	require.Equal(t, gold, result)
	require.NotEqual(t, gold, result.AddDefaultRemoteArgs("test"))
	gold.Args["branch"] = "master"
	gold.Args["container_name"] = "test_update"
	require.Equal(t, gold, result.AddDefaultRemoteArgs("test"))
}

func TestAddDefaultRemoteArgsToNil(t *testing.T) {
	t.Parallel()

	gold := Script{Ver: 1, Commands: []command.Command{command.MakeInclude("test_include")}}
	// gold_name := "test_args"
	require.NotEqual(t, gold, gold.AddDefaultRemoteArgs("test"))
	result := gold
	gold.Args = map[string]string{}
	gold.Args["branch"] = "master"
	gold.Args["container_name"] = "test_update"
	require.Equal(t, gold, result.AddDefaultRemoteArgs("test"))
}

func TestCreateFile(t *testing.T) {
	t.Parallel()

	_, err := os.Stat(fileName)
	require.NotNil(t, err)
	CreateFileIfNotExists(fileName)
	_, err = os.Stat(fileName)
	require.Nil(t, err)
	CreateFileIfNotExists(fileName)
	_, err = os.Stat(fileName)
	require.Nil(t, err)
}

func TestExtendArgsByConfig(t *testing.T) {
	t.Parallel()

	service := "test"
	arg := "branch"
	goldVal := "in_test"
	config.Set(service, arg, goldVal)
	args := []string{arg}
	tmp := Script{}
	result := tmp.ExtendArgsByDefaultAndGiven(args, service, config.Args(service))
	require.Equal(t, []string{goldVal}, result)
}

func TestExtendArgsByScript(t *testing.T) {
	t.Parallel()

	service := "test"
	arg := "container_name"
	config.Unset(service, arg)

	goldVal := "in_test_script"

	args := []string{arg}
	tmp := Script{Args: map[string]string{
		arg: goldVal,
	}}
	result := tmp.ExtendArgsByDefaultAndGiven(args, service, nil)
	require.Equal(t, []string{goldVal}, result)
}

func TestExtendArgsByScriptAndConfig(t *testing.T) {
	t.Parallel()

	service := "test2"
	arg := "container_name"
	goldVal := "in_test_script"
	goldVal2 := "in_test"
	config.Set(service, "branch", goldVal2)

	args := []string{arg, "branch"}
	tmp := Script{Args: map[string]string{
		arg: goldVal,
	}}
	result := tmp.ExtendArgsByDefaultAndGiven(args, service, config.Args(service))
	require.Equal(t, []string{goldVal, goldVal2}, result)
	config.Unset(service, "branch")
	nilArgs := tmp.ExtendArgsByDefaultAndGiven(nil, service, config.Args(service))
	require.Nil(t, nilArgs)
}

func TestExtendArgsNoVals(t *testing.T) {
	t.Parallel()

	service := "test3"
	arg := "not_existing_arg"
	args := []string{arg}
	tmp := Script{Args: map[string]string{}}
	result := tmp.ExtendArgsByDefaultAndGiven(args, service, nil)
	require.Equal(t, args, result)
}

func TestScipDefer(t *testing.T) {
	t.Parallel()

	service := "test4"
	arg := "defer"
	goldVal := "defer"
	config.Set(service, arg, goldVal)
	args := []string{arg}
	tmp := Script{Args: map[string]string{
		arg: goldVal,
	}}
	result := tmp.ExtendArgsByDefaultAndGiven(args, service, nil)
	require.Equal(t, []string{goldVal}, result)
}
