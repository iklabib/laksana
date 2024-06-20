package main

import (
	"net/http"
	"os"

	"codeberg.org/iklabib/markisa/model"
	"codeberg.org/iklabib/markisa/toolchains"
	"codeberg.org/iklabib/markisa/util"
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

		evaluator := toolchains.NewEvaluator(c.Request().Context())
		evaluationResult := evaluator.Submission(submmission)

		return c.JSON(http.StatusOK, evaluationResult)
	})

	e.Logger.Fatal(e.Start(BASE_URL))
}
