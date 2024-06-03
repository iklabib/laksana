package containers

import "codeberg.org/iklabib/markisa/model"

type Sandbox interface {
	ExecConfined(dir string, commands []string) model.SandboxExecResult
}
