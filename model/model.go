package model

import (
	"bytes"
)

// run stage model
type RunResult struct {
	Message  string         `json:"message"`
	Builds   []CompileError `json:"builds"`
	Tests    []TestResult   `json:"tests"`
	ExitCode int            `json:"exit_code"`
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

type CompileError struct {
	Filename string
	Message  string
	Line     int
	Column   int
}

type TestResult struct {
	Status string // PASS or FAILED
	Name   string
	Output string
	Order  int
}
