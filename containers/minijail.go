package containers

import (
	"bytes"
	"os"
	"os/exec"

	"codeberg.org/iklabib/markisa/model"
)

type Minijail struct {
	Path       string
	ConfigPath string
}

func NewMinijail(config string) Minijail {
	minijail, err := exec.LookPath("minijail0")
	if err != nil {
		panic(err)
	}

	file, err := os.CreateTemp("", "minijail*.cfg")
	if err != nil {
		panic(err)
	}

	file.WriteString(config)
	file.Close()

	return Minijail{
		Path:       minijail,
		ConfigPath: file.Name(),
	}
}

func (mn Minijail) argsBuilder(dir string, commands []string) []string {
	// keep in mind that minijail need absolute path
	// is there a way for it to look in path without bash invocation?
	args := []string{"--config", mn.ConfigPath, "-P", dir, "--"}
	return append(args, commands...)
}

func (mn Minijail) ExecConfined(dir string, commands []string) model.SandboxExecResult {
	args := mn.argsBuilder(dir, commands)

	var stdoutBuff bytes.Buffer
	var stderrBuff bytes.Buffer

	cmd := exec.Command(mn.Path, args...)
	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff
	err := cmd.Run()

	return model.SandboxExecResult{
		Error:  err,
		Stdout: stdoutBuff,
		Stderr: stderrBuff,
	}
}
