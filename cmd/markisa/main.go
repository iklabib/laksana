package main

import (
	"encoding/json"
	"net/http"
	"os"

	"codeberg.org/iklabib/markisa/containers"
	"codeberg.org/iklabib/markisa/db"
	"codeberg.org/iklabib/markisa/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	BASE_URL := os.Getenv("BASE_URL")

	e := echo.New()

	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	client := containers.NewDefault()
	record := db.NewRecord()

	record.CleanRecord()
	client.CleanContainers()

	exercise := db.NewExerciseDefault()

	e.POST("/run", func(c echo.Context) error {
		var req model.SubmissionRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}

		clientId := req.User
		
		containerId := record.Retrieve(clientId); 
		if containerId == "" {
			containerId, err := client.SpawnTenant()
			if err != nil {
				return c.String(http.StatusInternalServerError, "Internal Error")
			}

			record.Insert(clientId, containerId)
		}

		submmission := model.Submission{
			Type: req.Type,
			Src: req.Src,
			SrcTest: exercise.RetrieveTestCase(req.ExerciseId),
		}

		marshaled, err := json.Marshal(submmission)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Error")
		}

		output, err := client.ExecTenant(containerId[:12], marshaled)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		unmarshal := model.RunResult{}
		if err := json.Unmarshal(output, &unmarshal); err != nil {
			return c.String(http.StatusInternalServerError, "Internal Error")
		}
		return c.JSON(http.StatusOK, unmarshal)
	})

	e.Logger.Fatal(e.Start(BASE_URL))
}
