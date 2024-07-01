package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setup() {
	ReadDefaultFile()
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()
	os.Exit(code)
}

// func TestNotEmpty(t *testing.T) {
// 	t.Run("Empty", func(t *testing.T) {
// 		require.False(t, Empty())
// 	})
// }

func TestSet(t *testing.T) {
	t.Parallel()
	t.Run("Set", func(t *testing.T) {
		t.Parallel()
		Set("test", "foo", "bar")
		result, ok := Get("test", "foo")
		require.True(t, ok)
		require.Equal(t, "bar", result)
	})
}

func TestUnset(t *testing.T) {
	t.Parallel()
	t.Run("Unset", func(t *testing.T) {
		t.Parallel()
		Unset("test", "foo")
		result, ok := Get("test", "foo")
		require.False(t, ok)
		require.Equal(t, "", result)
	})
}

func TestSelect(t *testing.T) {
	t.Parallel()
	t.Run("AddSelect", func(t *testing.T) {
		t.Parallel()
		AddSelect("test", "foo")
		result, ok := GetSelects("test")
		require.True(t, ok)
		require.Equal(t, []string{"foo"}, result)
		// reset
		AddSelect("test", "foo")
		result, ok = GetSelects("test")
		require.True(t, ok)
		require.Equal(t, []string{"foo"}, result)
	})
}
