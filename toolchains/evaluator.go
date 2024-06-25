package toolchains

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/iklabib/laksana/containers"
	"codeberg.org/iklabib/laksana/model"
	"codeberg.org/iklabib/laksana/util"
)

const (
	INTERNAL_ERROR = -1
	RUNTIME_ERROR  = -2
)

type Evaluator struct {
	Workdir string
}

func NewEvaluator(workdir string) *Evaluator {
	return &Evaluator{
		Workdir: workdir,
	}
}

func (ev Evaluator) Submission(ctx context.Context, submission model.Submission) model.RunResult {
	switch submission.Type {
	case "python":
		configPath, _ := filepath.Abs("configs/minijail/minijail.cfg")
		minijail := containers.NewMinijail(ctx, configPath)
		python := NewPython()
		dir, err := python.Prep(submission)
		if err != nil {
			return model.RunResult{
				ExitCode: INTERNAL_ERROR,
				Message:  err.Error(),
			}
		}

		return python.Eval(dir, minijail)

	case "go":
		configPath, _ := filepath.Abs("configs/minijail/go.cfg")
		minijail := containers.NewMinijail(ctx, configPath)
		golang := NewGolang(ctx, ev.Workdir)
		dir, err := golang.Prep(submission)
		if err != nil {
			return model.RunResult{
				ExitCode: INTERNAL_ERROR,
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
		return model.RunResult{
			ExitCode: -1,
			Message:  fmt.Sprintf(`"%s is not supported"`, submission.Type),
		}
	}
}
