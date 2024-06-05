package db

import (
	"log"
	"os"
	"path"
)

type Exercise struct{
	path string
}

func NewExerciseDefault() *Exercise {
	exercise := Exercise{}
	cwd, _ := os.Getwd()
	exercise.path = path.Join(cwd, "..", "..", "exercises")
	return &exercise
}

func (ex *Exercise) RetrieveTestCase(exerciseId string) string{
	file := path.Join(ex.path, exerciseId, "test.py")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Fatalf("exercise %s does not exist\n", exerciseId)
	}

	buff, err := os.ReadFile(file)
	if (err != nil) {
		log.Fatalf("failed to read %s test\n", exerciseId)
	}
	return string(buff)
}