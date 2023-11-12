package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("", func(t *testing.T) {
		expected := make(Environment, 5)
		expected["BAR"] = EnvValue{
			Value:      "bar",
			NeedRemove: false,
		}
		expected["EMPTY"] = EnvValue{
			Value:      "",
			NeedRemove: false,
		}
		expected["FOO"] = EnvValue{
			Value:      "   foo\nwith new line",
			NeedRemove: false,
		}
		expected["HELLO"] = EnvValue{
			Value:      `"hello"`,
			NeedRemove: false,
		}
		expected["UNSET"] = EnvValue{
			Value:      "",
			NeedRemove: true,
		}

		env, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t, expected, env)
	})
}
