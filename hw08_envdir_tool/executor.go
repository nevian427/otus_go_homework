package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		if err := os.Setenv(k, v.Value); err != nil {
			return 128
		}
	}

	proc := exec.Command(cmd[0], cmd[1:]...)
	proc.Stderr = os.Stderr
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	if err := proc.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
	}
	return
}
