package toolchains

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"codeberg.org/iklabib/markisa/model"
	"codeberg.org/iklabib/markisa/util"
)

type Python struct{}

func NewPython() *Python {
	return &Python{}
}

// create temporary directory and write source code to file
func (p Python) Prep(src string, srcTest string) (string, error) {
	tempDir, err := os.MkdirTemp(".", "box")
	if err != nil {
		return "", errors.New("failed to create temp dir")
	}

	submmission := filepath.Join(tempDir, "main.py")
	if file, err := os.Create(submmission); err != nil {
		return "", errors.New("failed to write to file")
	} else {
		file.WriteString(src)
		file.Close()
	}

	testCases := filepath.Join(tempDir, "test.py")
	if file, err := os.Create(testCases); err != nil {
		return "", errors.New("failed to write to file")
	} else {
		file.WriteString(srcTest)
		file.Close()
	}

	return tempDir, nil
}

func (p Python) Eval(dir string) model.RunResult {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("python3", "test.py")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = dir

	err := cmd.Run()

	// we assumed that non-zero exit is error
	exitCode := util.GetExitCode(&err)
	status := "SUCCESS"
	if exitCode != 0 {
		status = "FAILED"
	}

	return model.RunResult{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Status:   status,
	}
}
