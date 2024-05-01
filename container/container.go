package container

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/iklabib/markisa/model"
)

func RunContainer(src []byte, image string) model.RunResult {
	cmd := exec.Command(
		"podman",
		"run",
		"-e",
		"MARKISA_RUN_TIME=3000",
		"-i",
		"--rm",
		image,
	)

	var stdout bytes.Buffer
	cmd.Stdin = bytes.NewReader(src)
	cmd.Stdout = &stdout

	cmd.Run()
	runResult := model.RunResult{}
	str := stdout.String()
	err := json.Unmarshal([]byte(str), &runResult)
	if err != nil {
		return model.RunResult{
			ExitCode: -1,
			Status:   "INTERNAL_ERROR",
		}
	}

	return runResult
}
