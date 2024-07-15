package toolchains

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"codeberg.org/iklabib/laksana/containers"
	"codeberg.org/iklabib/laksana/model"
	"codeberg.org/iklabib/laksana/util"
)

var executableName = "LittleRosie"

type CSharp struct {
	Ctx        context.Context
	Workdir    string
	Submission model.Submission
}

func NewCSharp(workdir string, Submission model.Submission) *CSharp {
	return &CSharp{
		Workdir:    workdir,
		Submission: Submission,
	}
}

func (cs CSharp) Run(sandbox containers.Sandbox) model.RunResult {
	dir, err := cs.Prep()
	if err != nil {
		return model.RunResult{Message: err.Error()}
	}

	defer func() {
		os.RemoveAll(dir)
	}()

	if buildErrors, err := cs.Build(dir, sandbox); err != nil {
		return model.RunResult{
			Message: err.Error(),
			Builds:  buildErrors,
		}
	}

	testResult, err := cs.Eval(dir, sandbox)
	if err != nil {
		return model.RunResult{Message: err.Error()}
	}

	return model.RunResult{
		Success: true,
		Tests:   testResult,
	}
}

func (cs CSharp) Prep() (string, error) {
	tempDir, err := CreateBox(cs.Workdir)
	if err != nil {
		return tempDir, err
	}

	os.Mkdir(filepath.Join(tempDir, "/dev"), 0o751)
	os.Mkdir(filepath.Join(tempDir, "/etc"), 0o751)

	// NUnit looking for /etc/passwd for whatever reason
	passwd := "ubuntu:x:1000:1000:Ubuntu:/home/ubuntu:/bin/bash"
	util.CreateROFile(filepath.Join(tempDir, "etc", "passwd"), passwd)

	runnerDest := filepath.Join(tempDir, "CSharp")
	runnerSource := filepath.Join("runner", "CSharp")

	// FIXME: it would better not to copy runner for every box
	// linking and binding doesn't seem to work
	err = exec.Command("cp", "-r", runnerSource, runnerDest).Run()
	if err != nil {
		return tempDir, err
	}

	return tempDir, nil
}

func (cs CSharp) Build(dir string, sandbox containers.Sandbox) ([]model.BuildError, error) {
	submissions := model.SourceCode{
		SourceCodeTest: cs.Submission.SourceCodeTest,
		SourceCodes:    cs.Submission.SourceFiles,
	}

	marshaled, err := json.Marshal(submissions)
	if err != nil {
		return nil, err
	}

	commands := []string{filepath.Join("CSharp", executableName), "build"}
	result := sandbox.ExecConfinedWithStdin(dir, commands, bytes.NewReader(marshaled))

	if result.Error != nil {
		return nil, fmt.Errorf("internal error: %s", &result.Stderr)
	}

	buildResult := struct {
		Status  int                `json:"status"`
		Message string             `json:"message"`
		Builds  []model.BuildError `json:"compilation_errors"`
	}{}

	if err := json.Unmarshal(result.Stdout.Bytes(), &buildResult); err != nil {
		return nil, err
	}

	switch buildResult.Status {
	case 0: // success
		return nil, nil
	case 1: // compilation error
		return buildResult.Builds, nil
	case 2: // internal error
		return nil, fmt.Errorf("internal error: %s", buildResult.Message)
	default:
		return nil, fmt.Errorf("internal error: unknown status '%d'", buildResult.Status)
	}
}

func (cs CSharp) Eval(dir string, sandbox containers.Sandbox) ([]model.TestResult, error) {
	executable := filepath.Join("CSharp", executableName)
	result := sandbox.ExecConfined(dir, []string{executable, "execute", "main.dll"})

	if result.Error != nil {
		if exitCode := util.GetExitCode(result.Error); exitCode > 1 {
			return nil, errors.New(result.Stderr.String())
		}
	}

	var testResult []model.TestResult
	if err := json.Unmarshal(result.Stdout.Bytes(), &testResult); err != nil {
		return nil, err
	}

	return testResult, nil
}
