package util

import (
	"os"
	"os/exec"
)

func GetExitCode(err *error) int {
	if (*err) == nil {
		return 0
	}

	if exitError, ok := (*err).(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return 1
}

func CreateTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "box_")
	if err != nil {
		return "", err
	}
	return tempDir, nil
}
