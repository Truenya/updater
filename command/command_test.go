package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeSSHArgs(t *testing.T) {
	t.Parallel()

	t.Run("Test make ssh args", func(t *testing.T) {
		t.Parallel()

		testArgs := []string{"arg1", "arg2"}
		result := MakeSSHArgs("msg", "c", testArgs)
		require.Equal(t, Command{Msg: "msg", Type: SSH, Data: "c", Args: testArgs}, result)
	})
}

func TestMakeSSH(t *testing.T) {
	t.Parallel()

	t.Run("Test make ssh", func(t *testing.T) {
		t.Parallel()

		result := MakeSSH("msg", "c")
		require.Equal(t, Command{Msg: "msg", Type: SSH, Data: "c"}, result)
	})
}

func TestMakeRsync(t *testing.T) {
	t.Parallel()

	t.Run("Test make rsync", func(t *testing.T) {
		t.Parallel()

		result := MakeRsync("msg", []string{"arg1", "arg2"})
		require.Equal(t, Command{Msg: "msg", Type: RSYNC, Args: []string{"arg1", "arg2"}}, result)
	})
}

func TestMakeLocal(t *testing.T) {
	t.Parallel()

	t.Run("Test make local", func(t *testing.T) {
		t.Parallel()

		result := MakeLocal("msg", "c")
		require.Equal(t, Command{Msg: "msg", Type: LOCAL, Data: "c"}, result)
	})
}

func TestMakeLocalArgs(t *testing.T) {
	t.Parallel()

	t.Run("Test make local args", func(t *testing.T) {
		t.Parallel()

		testArgs := []string{"arg1", "arg2"}
		result := MakeLocalArgs("msg", "c", testArgs)
		require.Equal(t, Command{Msg: "msg", Type: LOCAL, Data: "c", Args: testArgs}, result)
	})
}

func TestMakeInclude(t *testing.T) {
	t.Parallel()

	t.Run("Test make include", func(t *testing.T) {
		t.Parallel()

		result := MakeInclude("c")
		require.Equal(t, Command{Type: INCLUDE, Data: "c"}, result)
	})
}

func TestMakeIncludes(t *testing.T) {
	t.Parallel()

	t.Run("Test make includes", func(t *testing.T) {
		t.Parallel()

		result := MakeIncludes("c", "d")
		require.Equal(t, []Command{{Type: INCLUDE, Data: "c"}, {Type: INCLUDE, Data: "d"}}, result)
	})
}

func TestCommand_GetResultingCmdWithArgs(t *testing.T) {
	t.Parallel()

	t.Run("Test get resulting cmd with args", func(t *testing.T) {
		t.Parallel()

		result := Command{Msg: "msg", Type: SSH, Data: "c %s %s", Args: []string{"arg1", "arg2", "defer"}}
		require.Equal(t, "c arg1 arg2", result.GetResultingCmdWithArgs())
		result = Command{Msg: "msg", Type: SSH, Data: "c %s %s", Args: []string{"arg1", "arg2"}}
		require.Equal(t, "c arg1 arg2", result.GetResultingCmdWithArgs())
		result = Command{Msg: "msg", Type: SSH, Data: "c", Args: []string{}}
		require.Equal(t, "c", result.GetResultingCmdWithArgs())
		result = Command{Msg: "msg", Type: SSH, Data: "c", Args: nil}
		require.Equal(t, "c", result.GetResultingCmdWithArgs())
	})
}

func TestIsDeferred(t *testing.T) {
	t.Parallel()

	t.Run("Test is deferred", func(t *testing.T) {
		t.Parallel()

		result := Command{Msg: "msg", Type: SSH, Data: "c %s %s", Args: []string{"arg1", "arg2", "defer"}}
		require.True(t, result.IsDeferred())
		result = Command{Msg: "msg", Type: SSH, Data: "c %s %s", Args: []string{"arg1", "arg2"}}
		require.False(t, result.IsDeferred())
		result = Command{Msg: "msg", Type: SSH, Data: "c", Args: []string{}}
		require.False(t, result.IsDeferred())
		result = Command{Msg: "msg", Type: SSH, Data: "c", Args: nil}
		require.False(t, result.IsDeferred())
	})
}
