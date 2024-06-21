package main

import (
	"context"
	"net/http"
	"os"

	"codeberg.org/iklabib/laksana/model"
	"codeberg.org/iklabib/laksana/toolchains"
	"codeberg.org/iklabib/laksana/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	BASE_URL := os.Getenv("BASE_URL")

	e := echo.New()
	e.Validator = util.NewEchoValidator()

	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	e.POST("/run", func(c echo.Context) error {
		var submmission model.Submission
		if err := c.Bind(&submmission); err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}

		if err := c.Validate(submmission); err != nil {
			return err
		}

		ctx := c.Request().Context()
		evaluator := toolchains.NewEvaluator(ctx)
		resultChan := make(chan model.RunResult)

		go func() {
			defer close(resultChan)
			resultChan <- evaluator.Submission(submmission)
		}()

		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != context.Canceled {
				return c.JSON(http.StatusOK, model.RunResult{
					ExitCode: util.GetExitCode(&err),
					Message:  err.Error(),
				})
			} else {
				return c.JSON(http.StatusOK, model.RunResult{
					ExitCode: util.GetExitCode(&err),
					Message:  "canceled",
				})
			}

		case result := <-resultChan:
			return c.JSON(http.StatusOK, result)
		}
	})

	e.Logger.Fatal(e.Start(BASE_URL))
}
