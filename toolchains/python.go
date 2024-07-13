package toolchains

import (
	"errors"
	"os"
	"path/filepath"

	"codeberg.org/iklabib/laksana/containers"
	"codeberg.org/iklabib/laksana/model"
	"codeberg.org/iklabib/laksana/util"
)

type Python struct {
	Submission model.Submission
}

func NewPython() *Python {
	return &Python{}
}

// create temporary directory and write source code to file
func (p Python) Prep(submission model.Submission) (string, error) {
	tempDir, err := CreateBox(".")
	if err != nil {
		return tempDir, err
	}

	err = WriteSourceCodes(tempDir, submission.SourceFiles)
	if err != nil {
		return tempDir, err
	}

	err = util.CreateROFile(filepath.Join(tempDir, "test.py"), submission.SourceCodeTest)
	if err != nil {
		return tempDir, errors.New("failed to write test.py")
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
