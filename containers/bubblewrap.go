package containers

import (
	"bytes"
	"os/exec"

	"codeberg.org/iklabib/markisa/model"
)

// TODO: resource limit

type Bubblewrap struct {
	Path string
}

func NewBubblewrap() Bubblewrap {
	var bwrap = Bubblewrap{}
	path, err := exec.LookPath("bwrap")
	if err != nil {
		panic(err)
	}
	bwrap.Path = path
	return bwrap
}

// TODO: use --args FD instead of this mess
func (bw Bubblewrap) argsBuilder(dir string, commands []string) []string {
	homePath := "/home/user"
	args := []string{
		"--unshare-user",
		"--unshare-ipc",
		"--unshare-pid",
		"--unshare-cgroup",
		"--clearenv",
		"--ro-bind", "/usr/bin", "/usr/bin",
		"--ro-bind", "/usr/lib", "/usr/lib",
		// "--ro-bind", "/usr/lib64", "/usr/lib64",
		"--ro-bind", "/usr/include", "/usr/include",
		"--ro-bind", "/bin", "/bin",
		"--ro-bind", "/lib", "/lib",
		// "--ro-bind", "/lib64", "/lib64",
		"--bind", dir, homePath,
		"--proc", "/proc",
		"--tmpfs", "/tmp",
		"--chdir", homePath,
		"--setenv", "PWD", homePath,
		"--setenv", "HOME", homePath,
		"--setenv", "USER", "user",
		"--setenv", "USERNAME", "user",
		"--setenv", "PATH", "/bin:/usr/bin",
		"--new-session",
		"--",
	}

	return append(args, commands...)
}

func (bw Bubblewrap) ExecConfined(dir string, commands []string) model.SandboxExecResult {
	args := bw.argsBuilder(dir, commands)

	var stdoutBuff bytes.Buffer
	var stderrBuff bytes.Buffer

	cmd := exec.Command(bw.Path, args...)
	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff
	cmd.Dir = dir
	err := cmd.Run()

	return model.SandboxExecResult{
		Error:  err,
		Stdout: stdoutBuff.String(),
		Stderr: stderrBuff.String(),
	}
}
