package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("without env", func(t *testing.T) {
		cmd := append([]string{}, "echo", "HELLO", "BEAUTIFUL", "WORLD")

		resCode := -1
		resOut := captureStdout(func() {
			resCode = RunCmd(cmd, nil)
		})

		require.Equal(t, 0, resCode)
		require.Equal(t, "HELLO BEAUTIFUL WORLD\n", resOut)
	})

	t.Run("with env to set", func(t *testing.T) {
		cmd := append([]string{}, "printenv", "FOO", "BAR")
		env := make(Environment, 2)
		env["FOO"] = EnvValue{
			Value:      "foo",
			NeedRemove: false,
		}
		env["BAR"] = EnvValue{
			Value:      "bar",
			NeedRemove: false,
		}

		resCode := -1
		resOut := captureStdout(func() {
			resCode = RunCmd(cmd, env)
		})

		require.Equal(t, 0, resCode)
		require.Equal(t, "foo\nbar\n", resOut)
	})

	t.Run("with env to unset", func(t *testing.T) {
		if err := os.Setenv("TEST", "test"); err != nil {
			return
		}

		cmd := append([]string{}, "printenv", "TEST")
		resCode := -1
		resOut := captureStdout(func() {
			resCode = RunCmd(cmd, nil)
		})
		require.Equal(t, 0, resCode)
		require.Equal(t, "test\n", resOut)

		env := make(Environment, 2)
		env["TEST"] = EnvValue{
			Value:      "",
			NeedRemove: true,
		}

		cmd = append([]string{}, "printenv", "TEST")
		resOut = captureStdout(func() {
			resCode = RunCmd(cmd, env)
		})
		require.Equal(t, 1, resCode)
		require.Equal(t, "", resOut)
	})
}

func captureStdout(f func()) string {
	// save curr stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// call func
	f()

	// restore stdout
	w.Close()
	os.Stdout = old

	// read from pipe
	out := make(chan string)
	go func() {
		var buf [4024]byte
		n, _ := r.Read(buf[:])
		out <- string(buf[:n])
	}()

	actual := <-out

	return actual
}
