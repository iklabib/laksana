package toolchains

import (
	"log"
	"os"

	"codeberg.org/iklabib/markisa/containers"
	"codeberg.org/iklabib/markisa/model"
)

func EvaluateSubmission(submission model.Submission) model.RunResult {
	minijail := containers.NewMinijail(submission.SandboxConfig)
	switch submission.Type {
	case "python":
		python := NewPython()
		dir, err := python.Prep(submission)
		if err != nil {
			log.Println(err.Error())
			return model.RunResult{}
		}

		return python.Eval(dir, minijail)

	case "go":
		golang := NewGolang()
		dir, err := golang.Prep(submission)
		if err != nil {
			log.Println(err.Error())
			return model.RunResult{}
		}

		result := golang.Eval(dir, minijail)
		if err := os.RemoveAll(dir); err != nil {
			log.Println(err.Error())
			return model.RunResult{}
		}

		return result

	default:
		log.Printf(`"%s is not supported"`, submission.Type)
		return model.RunResult{}
	}
}
