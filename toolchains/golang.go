package toolchains

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"codeberg.org/iklabib/laksana/containers"
	"codeberg.org/iklabib/laksana/model"
	"codeberg.org/iklabib/laksana/util"
)

type Golang struct {
	Ctx     context.Context
	Workdir string
}

func NewGolang(ctx context.Context, workdir string) *Golang {
	return &Golang{
		Ctx:     ctx,
		Workdir: workdir,
	}
}

func (g Golang) Prep(submission model.Submission) (string, error) {
	cwd, _ := os.Getwd()
	tempDir, err := os.MkdirTemp(g.Workdir, "box_*")
	if err != nil {
		return tempDir, errors.New("failed to create temp dir")
	}

	if err := os.Chown(tempDir, 1000, 1000); err != nil {
		return tempDir, errors.New("failed to set directory permission")
	}

	err = util.CreateReadOnlyFile(filepath.Join(tempDir, "main.go"), []byte(submission.Src))
	if err != nil {
		return tempDir, errors.New("failed to write to main file")
	}

	err = util.CreateReadOnlyFile(filepath.Join(tempDir, "main_test.go"), []byte(submission.SrcTest))
	if err != nil {
		return tempDir, errors.New("failed to write to main_test file")
	}

	target := filepath.Join(cwd, "runner", "go")
	err1 := util.Copy(filepath.Join(target, "go.mod"), filepath.Join(tempDir, "go.mod"))
	err2 := util.Copy(filepath.Join(target, "run.bash"), filepath.Join(tempDir, "run.bash"))
	err3 := util.Copy(filepath.Join(target, "goenv.bash"), filepath.Join(tempDir, "goenv.bash"))

	if err1 != nil || err2 != nil || err3 != nil {
		return tempDir, errors.New("failed to copy runner files")
	}

	return tempDir, nil
}

func (g Golang) buildTest(executable, dir string) (bytes.Buffer, error) {
	stderr := bytes.Buffer{}
	cmd := exec.Command(executable, "run.bash", "build-test")
	cmd.Stderr = &stderr
	cmd.Dir = dir
	err := cmd.Run()
	return stderr, err
}

func (g Golang) Eval(dir string, sandbox containers.Sandbox) model.RunResult {
	executable := "/bin/bash"

	if stderr, err := g.buildTest(executable, dir); err != nil {
		return model.RunResult{
			ExitCode: util.GetExitCode(&err),
			Message:  err.Error(),
			Builds:   g.ParseCompileErrors(stderr),
		}
	}

	commands := []string{executable, "run.bash", "run"}
	execResult := sandbox.ExecConfined(dir, commands)

	// when exit code is 1 we can ignore it
	// it is likely because of test fail, not actual error
	exitCode := util.GetExitCode(&execResult.Error)
	if exitCode > 1 {
		return model.RunResult{
			ExitCode: exitCode,
			Message:  execResult.Error.Error(),
		}
	}

	testResult := g.ParseTestEvent(execResult.Stdout)
	return model.RunResult{
		ExitCode: exitCode,
		Tests:    testResult,
	}
}

func (g Golang) ParseTestEvent(out bytes.Buffer) []model.TestResult {
	// skip first line
	out.ReadString('\n')

	// test case order
	// action: run
	// action: output -> there are multiple of this
	// action: pass or action: fail

	var results []model.TestResult
	for idx := 1; ; {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}

		var testEvent goTestEvent
		if err := json.Unmarshal([]byte(line), &testEvent); err != nil {
			log.Printf("error unmarshalling JSON: %v", err)
			continue
		}

		// line action must be "run", else we done
		if testEvent.Action != "run" {
			break
		}

		testCase := model.TestResult{
			Name:  testEvent.Test,
			Order: idx,
		}

		var output []string
	loop:
		for {
			line, err := out.ReadString('\n')
			if err != nil {
				break
			}

			line = strings.Trim(line, "\n")

			var event goTestEvent
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				log.Printf("error unmarshalling JSON: %v", err)
				continue
			}

			switch event.Action {
			case "output":
				output = append(output, event.Output)
			case "fail", "pass":
				testCase.Status = event.Action
				joinedOutput := strings.Join(output[1:len(output)-1], "")
				testCase.Output = strings.TrimSpace(joinedOutput)
				break loop
			}
		}

		idx++
		results = append(results, testCase)
	}

	return results
}

func (g Golang) ParseCompileErrors(out bytes.Buffer) []model.CompileError {
	var compileErrors []model.CompileError
	// first line should be module name, skip it
	out.ReadString('\n')

	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}

		compileError, err := parseError(line)
		if err != nil {
			log.Printf("failed to parse: %s\n", line)
			continue
		}

		compileErrors = append(compileErrors, compileError)
	}

	return compileErrors
}

func parseError(out string) (model.CompileError, error) {
	var compileError model.CompileError

	parts := strings.Split(out, ":")
	if len(parts) < 4 {
		return compileError, fmt.Errorf("invalid error format")
	}

	compileError.Filename = parts[0]

	line, err := strconv.Atoi(parts[1])
	if err != nil {
		return compileError, fmt.Errorf("failed to parse line number")
	}
	compileError.Line = line

	column, err := strconv.Atoi(parts[2])
	if err != nil {
		return compileError, fmt.Errorf("failed to parse column number")
	}
	compileError.Column = column

	return compileError, nil
}

type goTestEvent struct {
	Time    time.Time `json:"time"`
	Action  string    `json:"action"`
	Package string    `json:"package"`
	Test    string    `json:"test,omitempty"`
	Output  string    `json:"output,omitempty"`
	Elapsed float64   `json:"elapsed,omitempty"`
}
