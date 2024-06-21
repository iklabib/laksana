package toolchains

import (
	"context"
	"fmt"
	"os"

	"codeberg.org/iklabib/laksana/containers"
	"codeberg.org/iklabib/laksana/model"
	"codeberg.org/iklabib/laksana/util"
)

type Evaluator struct {
	Ctx context.Context
}

func NewEvaluator(ctx context.Context) *Evaluator {
	return &Evaluator{
		Ctx: ctx,
	}
}

func (ev Evaluator) Submission(submission model.Submission) model.RunResult {
	minijail := containers.NewMinijail(ev.Ctx, submission.SandboxConfig)
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
		golang := NewGolang(ev.Ctx)
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
		return model.RunResult{
			ExitCode: -1,
			Message:  fmt.Sprintf(`"%s is not supported"`, submission.Type),
		}
	}
}
