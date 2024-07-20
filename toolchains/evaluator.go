package toolchains

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"codeberg.org/iklabib/laksana/containers"
	"codeberg.org/iklabib/laksana/model"
)

type Evaluator struct {
	Workdir string
}

func NewEvaluator(workdir string) *Evaluator {
	return &Evaluator{
		Workdir: workdir,
	}
}

func (ev Evaluator) Eval(ctx context.Context, submission model.Submission) model.RunResult {
	resultChan := make(chan model.RunResult)

	go func() {
		defer close(resultChan)
		resultChan <- ev.Submission(ctx, submission)
	}()

	select {
	case <-ctx.Done():
		var message string
		if err := ctx.Err(); errors.Is(err, context.Canceled) {
			message = "canceled"
		} else {
			message = err.Error()
		}

		return model.RunResult{Message: message}

	case res := <-resultChan:
		return res
	}
}

func (ev Evaluator) Submission(ctx context.Context, submission model.Submission) model.RunResult {
	switch submission.Type {
	case "go":
		return ev.EvalGo(ctx, submission)

	case "csharp":
		return ev.EvalCSharp(ctx, submission)

	default:
		return model.RunResult{
			Message: fmt.Sprintf(`"%s is not supported"`, submission.Type),
		}
	}
}

func (ev Evaluator) EvalGo(ctx context.Context, submission model.Submission) model.RunResult {
	configPath, _ := filepath.Abs("configs/minijail/go.cfg")
	minijail := containers.NewMinijail(ctx, configPath)
	golang := NewGolang(ctx, ev.Workdir)
	result := golang.Run(minijail)
	return result
}

func (ev Evaluator) EvalCSharp(ctx context.Context, submission model.Submission) model.RunResult {
	configPath, _ := filepath.Abs("configs/minijail/csharp.cfg")
	minijail := containers.NewMinijail(ctx, configPath)
	csharp := NewCSharp(ev.Workdir, submission)
	return csharp.Run(minijail)
}
