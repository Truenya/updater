package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultFilePath(t *testing.T) {
	t.Parallel()

	_, _, err := DefaultFilePath("config.json")
	if err != nil {
		t.Error(err)
	}
}

func TestLogNotNilErr(t *testing.T) {
	t.Parallel()

	LogNotNilErr(nil)
}

func TestProcessWithWarnLog(t *testing.T) {
	t.Parallel()

	ProcessWithWarnLog(func() error {
		return nil
	})
}

func TestInitCustomJson(t *testing.T) {
	t.Parallel()

	_, err := InitCustomJSON("./", "test.json")
	if err != nil {
		t.Error(err)
	}
}

func TestContainNumber(t *testing.T) {
	t.Parallel()

	require.True(t, ContainNumber("123"))
	require.False(t, ContainNumber("absdsdfsdfa"))
	require.True(t, ContainNumber("123.456"))
	require.True(t, ContainNumber("123.456.789"))
	require.True(t, ContainNumber("sadfkljsdhgkljdhfgsdklfjhg7"))
}

func TestError(t *testing.T) {
	t.Parallel()

	s := DefaultFilePathError{fmt.Errorf("test")}.Error()
	require.Equal(t, "Failed to get config directory: test", s)
}
