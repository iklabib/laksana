package toolchains

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"path/filepath"

	"codeberg.org/iklabib/laksana/containers"
	"codeberg.org/iklabib/laksana/model"
	"codeberg.org/iklabib/laksana/util"
)

type CSharp struct {
	Ctx        context.Context
	Workdir    string
	Submission model.Submission
}

func NewCSharp(workdir string) *CSharp {
	return &CSharp{
		Workdir: workdir,
	}
}

func (cs CSharp) Run(sandbox containers.Sandbox) model.RunResult {
	dir, err := cs.Prep()
	if err != nil {
		err = errors.New("preparation failure")
		return model.RunResult{
			ExitCode: model.INTERNAL_ERROR,
			Message:  err.Error(),
		}
	}

	if buildErrors, err := cs.Build(dir, sandbox); err != nil {
		return model.RunResult{
			ExitCode: util.GetExitCode(&err),
			Message:  err.Error(),
			Builds:   buildErrors,
		}
	}

	runResult := cs.Eval(dir, sandbox)

	// os.RemoveAll(dir)

	return runResult
}

func (cs CSharp) Prep() (string, error) {
	tempDir, err := CreateBox(cs.Workdir)
	if err != nil {
		return tempDir, err
	}

	scriptDest := filepath.Join(tempDir, "run.bash")
	scriptSource := filepath.Join("runner", "CSharp", "run.bash")
	if err := util.Copy(scriptSource, scriptDest); err != nil {
		return tempDir, err
	}

	runnerDest := filepath.Join(tempDir, "output")
	runnerSource := filepath.Join("runner", "CSharp", "output")
	if err := exec.Command("cp", "-r", runnerSource, runnerDest).Run(); err != nil {
		return tempDir, err
	}

	return tempDir, nil
}

func (cs CSharp) Build(dir string, sandbox containers.Sandbox) ([]model.BuildError, error) {
	var compileErrors []model.BuildError

	testSubmission := model.SourceFile{
		Name:       "Test.cs",
		Path:       "",
		SourceCode: cs.Submission.SourceCodeTest,
	}
	submissions := append(cs.Submission.SourceCode, testSubmission)
	marshaled, err := json.Marshal(submissions)
	if err != nil {
		log.Println("marshal failure")
		return compileErrors, err
	}

	commands := []string{"/bin/bash", "run.bash", "build"}
	result := sandbox.ExecConfinedWithStdin(dir, commands, bytes.NewReader(marshaled))

	exitCode := util.GetExitCode(&result.Error)
	if exitCode > 1 {
		return compileErrors, result.Error
	}

	err = json.Unmarshal(result.Stdout.Bytes(), &compileErrors)
	return compileErrors, err
}

func (cs CSharp) Eval(dir string, sandbox containers.Sandbox) model.RunResult {
	commands := []string{"/bin/bash", "run.bash", "execute"}
	result := sandbox.ExecConfined(dir, commands)
	exitCode := util.GetExitCode(&result.Error)
	if exitCode > 1 {
		return model.RunResult{
			ExitCode: exitCode,
		}
	}

	var testResult []model.TestResult
	if err := json.Unmarshal(result.Stdout.Bytes(), &testResult); err != nil {
		return model.RunResult{
			ExitCode: util.GetExitCode(&err),
			Message:  "unmarshal failure",
		}
	}

	return model.RunResult{
		ExitCode: 0,
		Tests:    testResult,
	}
}
