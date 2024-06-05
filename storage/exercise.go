package storage

import (
	"fmt"
	"os"
	"path"
)

type Exercise struct {
	path string
}

func NewExerciseDefault() *Exercise {
	exercise := Exercise{}
	cwd, _ := os.Getwd()
	exercise.path = path.Join(cwd, "exercises")
	return &exercise
}

func (ex *Exercise) RetrieveTestCase(exerciseId string) (string, error) {
	file := path.Join(ex.path, exerciseId, "test.py")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return "", fmt.Errorf("exercise %s does not exist", exerciseId)
	}

	buff, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read %s test", exerciseId)
	}
	return string(buff), nil
}
