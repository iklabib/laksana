package containers

import (
	"io"

	"codeberg.org/iklabib/laksana/model"
)

type Sandbox interface {
	ExecConfined(dir string, commands []string) model.SandboxExecResult
	ExecConfinedWithStdin(dir string, commands []string, stdin io.Reader) model.SandboxExecResult
}
