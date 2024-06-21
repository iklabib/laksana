package containers

import "codeberg.org/iklabib/laksana/model"

type Sandbox interface {
	ExecConfined(dir string, commands []string) model.SandboxExecResult
}
