package routes

import (
	"io"
	"os"
	"path/filepath"

	"github.com/dchest/uniuri"
	"github.com/labstack/echo/v4"
)

func GetFile(c echo.Context) error {
	filename := filepath.Join("public", "shareable", c.Param("id")+".go")
	_, stat := os.Stat(filename)
	if os.IsNotExist(stat) {
		c.JSON(203, "Sharable does not exist")
	}

	return c.File(filename)
}

func ShareFile(c echo.Context) error {
	id := uniuri.New()
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return echo.ErrBadRequest
	}

	file, err := fileHeader.Open()
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer file.Close()

	filename := filepath.Join("public", "shareable", id+".go")
	wr, err := os.Create(filename)
	if err != nil {
		return echo.ErrInternalServerError
	}
	if _, err = io.Copy(wr, file); err != nil {
		return echo.ErrInternalServerError
	}

	return c.String(204, id)
}
