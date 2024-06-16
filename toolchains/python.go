package toolchains

import (
	"errors"
	"os"
	"path/filepath"

	"codeberg.org/iklabib/markisa/containers"
	"codeberg.org/iklabib/markisa/model"
)

type Python struct{}

func NewPython() *Python {
	return &Python{}
}

// create temporary directory and write source code to file
func (p Python) Prep(submission model.Submission) (string, error) {
	tempDir, err := os.MkdirTemp(".", "box")
	if err != nil {
		return tempDir, errors.New("failed to create temp dir")
	}

	if err := os.Chmod(tempDir, 0o0775); err != nil {
		return tempDir, errors.New("failed to set permission for temp dir")
	}

	err = os.WriteFile(filepath.Join(tempDir, "main.py"), []byte(submission.Src), 0o0444)
	if err != nil {
		return tempDir, errors.New("failed to write to file")
	}

	err = os.WriteFile(filepath.Join(tempDir, "test.py"), []byte(submission.SrcTest), 0o0444)
	if err != nil {
		return tempDir, errors.New("failed to write to file")
	}

	return tempDir, nil
}

func (p Python) Eval(dir string, sandbox containers.Sandbox) model.RunResult {
	// executable, err := exec.LookPath("python3")
	// if err != nil {
	//	return model.RunResult{}
	//}

	// commands := []string{executable, "test.py"}
	// execResult := sandbox.ExecConfined(dir, commands)

	if err := os.RemoveAll(dir); err != nil {
		return model.RunResult{}
	}

	return model.RunResult{}
}
