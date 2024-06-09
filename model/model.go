package model

type Box struct {
	Id      string
	Type    string
	Version string
}

// buil stage model
type BuildResult struct {
	ExitCode   int
	Status     string
	Stdout     string
	Stderr     string
	Executable []byte
}

type BuilderResponse struct {
	Err string // internal error reason
	BuildResult
}

// run stage model
type RunResult struct {
	ExitCode int
	Status   string
	Stdout   string
	Stderr   string
}

// Run endpoint model
type RunResponse struct {
	Build  BuildResult
	Run    RunResult
	Status string
}

type Submission struct {
	Type          string `json:"type"`
	Src           string `json:"src"`
	SrcTest       string `json:"src_test"`
	SandboxConfig string `json:"sandbox_config" validate:"required"`
}

type SandboxExecResult struct {
	Error  error
	Stdout string
	Stderr string
}
