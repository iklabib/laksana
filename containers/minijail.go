package containers

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"codeberg.org/iklabib/laksana/model"
)

type Minijail struct {
	Ctx        context.Context
	Path       string
	ConfigPath string
}

func NewMinijail(ctx context.Context, config string) Minijail {
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
		Ctx:        ctx,
		Path:       minijail,
		ConfigPath: file.Name(),
	}
}

func (mn Minijail) Clean() {
	os.Remove(filepath.Dir(mn.ConfigPath))
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

	cmd := exec.CommandContext(mn.Ctx, mn.Path, args...)
	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff
	err := cmd.Run()

	return model.SandboxExecResult{
		Error:  err,
		Stdout: stdoutBuff,
		Stderr: stderrBuff,
	}
}
