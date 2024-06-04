package containers

import (
	"bytes"
	"os"
	"os/exec"
	"path"

	"codeberg.org/iklabib/markisa/model"
)

type Minijail struct {
	Path       string
	ConfigPath string
}

func NewMinijail() Minijail {
	minijail, err := exec.LookPath("minijail0")
	if err != nil {
		panic(err)
	}

	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return Minijail{
		Path:       minijail,
		ConfigPath: path.Join(path.Dir(executable), "configs/minijail.cfg"),
	}
}

func (mn Minijail) argsBuilder(dir string, commands []string) []string {
	// keep in mind that minijail need absolute path
	// is there a way for it to look in path without bash invocation?
	args := []string{"--config", mn.ConfigPath, "-C", dir, "--"}
	return append(args, commands...)
}

func (mn Minijail) PreExec(dir string) {
	// minijail0 -t does not automatically create /tmp
	if err := os.Mkdir(path.Join(dir, "tmp"), 0765); err != nil {
		panic(err)
	}
}

func (mn Minijail) ExecConfined(dir string, commands []string) model.SandboxExecResult {
	mn.PreExec(dir)

	args := mn.argsBuilder(dir, commands)

	var stdoutBuff bytes.Buffer
	var stderrBuff bytes.Buffer

	cmd := exec.Command(mn.Path, args...)
	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff
	err := cmd.Run()

	return model.SandboxExecResult{
		Error:  err,
		Stdout: stdoutBuff.String(),
		Stderr: stderrBuff.String(),
	}
}
