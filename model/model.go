package model

import (
	"bytes"
)

// run stage model
type RunResult struct {
	Message  string       `json:"message"`
	Builds   []BuildError `json:"builds"`
	Tests    []TestResult `json:"tests"`
	ExitCode int          `json:"exit_code"`
}

type Submission struct {
	Type           string       `json:"type"`
	SourceCode     []SourceFile `json:"src"`
	SourceCodeTest string       `json:"src_test"`
}

type SourceFile struct {
	Name       string `json:"name"`
	Path       string `json:"path,omitempty"`
	SourceCode string `json:"src"`
}

type SandboxExecResult struct {
	Error  error
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

type BuildError struct {
	Filename  string `json:"name"`
	Message   string `json:"message"`
	Line      int    `json:"line"`
	Character int    `json:"character"`
}

type TestResult struct {
	Status string // PASS or FAILED
	Name   string
	Output string
	Order  int
}
