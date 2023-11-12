package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envName, envVal := range env {
		if envVal.NeedRemove {
			err := os.Unsetenv(envName)
			if err != nil {
				fmt.Println("failed to unset env:", err)
			}
			continue
		}

		err := os.Setenv(envName, envVal.Value)
		if err != nil {
			fmt.Println("failed to set env:", err)
		}
	}

	cmdName := cmd[0]
	args := cmd[1:]
	command := createCommandWithStd(cmdName, args)

	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}

	return 0
}

func createCommandWithStd(name string, args []string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd
}
