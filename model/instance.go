package model

const (
	BUILD_ERROR   = 2
	RUNTIME_ERROR = 3
)

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
