package container

import (
	"bytes"
	"encoding/json"
	"gitlab.com/iklabib/markisa/model"
	"os/exec"
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

	cmd.Run();
  runResult := model.RunResult{}
  json.Unmarshal(stdout.Bytes(), &runResult)

  return runResult
}
