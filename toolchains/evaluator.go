package toolchains

import (
	"fmt"

	"codeberg.org/iklabib/markisa/model"
)

func EvaluateSubmission(submission model.Submission) model.RunResult {
	switch submission.Type {
	case "python":
		python := NewPython()
		dir, err := python.Prep(submission.Src, submission.SrcTest)

		if err != nil {
			return model.RunResult{
				ExitCode: -1,
				Status:   "Preparation failed",
				Stderr:   err.Error(),
			}
		}

		return python.Eval(dir)

	default:
		return model.RunResult{
			ExitCode: -1,
			Status:   fmt.Sprintf("Unsupported type \"%s\"\n", submission.Type),
		}
	}
}
