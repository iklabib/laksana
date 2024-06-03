package toolchains

import (
	"errors"
	"os"
	"path/filepath"

	"codeberg.org/iklabib/markisa/containers"
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
		return tempDir, errors.New("failed to create temp dir")
	}

	submmission := filepath.Join(tempDir, "main.py")
	if file, err := os.Create(submmission); err != nil {
		return tempDir, errors.New("failed to write to file")
	} else {
		file.WriteString(src)
		file.Close()
	}

	testCases := filepath.Join(tempDir, "test.py")
	if file, err := os.Create(testCases); err != nil {
		return tempDir, errors.New("failed to write to file")
	} else {
		file.WriteString(srcTest)
		file.Close()
	}

	return tempDir, nil
}

func (p Python) Eval(dir string) model.RunResult {
	commands := []string{"$(which python3)", "test.py"}
	minijail := containers.NewMinijail()
	execResult := minijail.ExecConfined(dir, commands)
	exitCode := util.GetExitCode(&execResult.Error)

	// we assumed that non-zero exit is error
	status := "SUCCESS"
	if exitCode != 0 {
		status = "FAILED"
	}

	return model.RunResult{
		ExitCode: exitCode,
		Stdout:   execResult.Stdout,
		Stderr:   execResult.Stderr,
		Status:   status,
	}
}
