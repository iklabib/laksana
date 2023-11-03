package container

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"gitlab.com/iklabib/markisa/model"
)

func RunContainer(src []byte, image string) model.RunResult {
	cmd := exec.Command(
		"podman",
		"run",
		"-i",
		"--rm",
		image,
	)

	var stdout bytes.Buffer
	cmd.Stdin = bytes.NewReader(src)
	cmd.Stdout = &stdout

	cmd.Run()
	runResult := model.RunResult{}
	err := json.Unmarshal(stdout.Bytes(), &runResult)
	if err != nil {
		return model.RunResult{
			ExitCode: -1,
			Status:   "INTERNAL_ERROR",
		}
	}

	return runResult
}
