package toolchains

import (
	"log"
	"os"

	"codeberg.org/iklabib/markisa/containers"
	"codeberg.org/iklabib/markisa/model"
	"codeberg.org/iklabib/markisa/util"
)

func EvaluateSubmission(submission model.Submission) model.RunResult {
	minijail := containers.NewMinijail(submission.SandboxConfig)
	defer minijail.Clean()

	switch submission.Type {
	case "python":
		python := NewPython()
		dir, err := python.Prep(submission)
		if err != nil {
			return model.RunResult{
				ExitCode: util.GetExitCode(&err),
				Message:  err.Error(),
			}
		}

		return python.Eval(dir, minijail)

	case "go":
		golang := NewGolang()
		dir, err := golang.Prep(submission)
		if err != nil {
			return model.RunResult{
				ExitCode: util.GetExitCode(&err),
				Message:  err.Error(),
			}
		}

		result := golang.Eval(dir, minijail)
		if err := os.RemoveAll(dir); err != nil {
			return model.RunResult{
				ExitCode: util.GetExitCode(&err),
				Message:  err.Error(),
			}
		}

		return result

	default:
		log.Printf(`"%s is not supported"`, submission.Type)
		return model.RunResult{}
	}
}
