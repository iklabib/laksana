package toolchains

import (
	"errors"
	"os"
	"os/exec"
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

	if err := os.Chmod(tempDir, 0775); err != nil {
		return tempDir, errors.New("failed to set permission for temp dir")
	}

	submmission := filepath.Join(tempDir, "main.py")
	if file, err := os.Create(submmission); err != nil {
		return tempDir, errors.New("failed to write to file")
	} else {
		file.WriteString(src)
		file.Close()
		os.Chmod(submmission, 0444)
	}

	testCases := filepath.Join(tempDir, "test.py")
	if file, err := os.Create(testCases); err != nil {
		return tempDir, errors.New("failed to write to file")
	} else {
		file.WriteString(srcTest)
		file.Close()
		os.Chmod(submmission, 0444)
	}

	return tempDir, nil
}

func (p Python) Eval(dir string) model.RunResult {
	executable, err := exec.LookPath("python3")
	if err != nil {
		return model.RunResult{
			ExitCode: -2,
			Status:   "FAILED",
			Stderr:   err.Error(),
		}
	}
	commands := []string{executable, "test.py"}
	minijail := containers.NewMinijail()
	execResult := minijail.ExecConfined(dir, commands)
	exitCode := util.GetExitCode(&execResult.Error)

	// we assumed that non-zero exit is error
	status := "SUCCESS"
	if exitCode != 0 {
		status = "FAILED"
	}

	// TODO: add fallback mechanism
	if err := os.RemoveAll(dir); err != nil {
		return model.RunResult{
			ExitCode: -1,
			Status:   "Exit clean failed",
			Stderr:   err.Error(),
		}
	}

	return model.RunResult{
		ExitCode: exitCode,
		Stdout:   execResult.Stdout,
		Stderr:   execResult.Stderr,
		Status:   status,
	}
}
