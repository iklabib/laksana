package container

import (
	"bytes"
	"os/exec"
	"strings"
)

func RunContainer(src []byte, image string) (string, string, error) {
	cmd := exec.Command(
		"podman",
		"run",
		"-i",
		"--rm",
		image,
	)

	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdin = bytes.NewReader(src)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}
