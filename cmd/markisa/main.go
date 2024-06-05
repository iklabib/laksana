package main

import (
	"net/http"
	"os"

	"codeberg.org/iklabib/markisa/model"
	"codeberg.org/iklabib/markisa/storage"
	"codeberg.org/iklabib/markisa/toolchains"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	BASE_URL := os.Getenv("BASE_URL")

	e := echo.New()

	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	exercise := storage.NewExerciseDefault()

	e.POST("/run", func(c echo.Context) error {
		var req model.SubmissionRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}

		if err := c.Validate(req); err != nil {
			return err
		}

		testCase, err := exercise.RetrieveTestCase(req.ExerciseId)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request: "+err.Error())
		}

		submmission := model.Submission{
			Type:    req.Type,
			Src:     req.Src,
			SrcTest: testCase,
		}

		evaluationResult := toolchains.EvaluateSubmission(submmission)

		// 0 success
		// -1 is internal error
		// -2 is evaluation error
		if evaluationResult.ExitCode == -1 {
			return c.String(http.StatusInternalServerError, "Internal Error")
		}
		return c.JSON(http.StatusOK, evaluationResult)
	})

	e.Logger.Fatal(e.Start(BASE_URL))
}
