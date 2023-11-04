package routes

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	CRunner "gitlab.com/iklabib/markisa/runner/c"

	"github.com/labstack/echo/v4"
)

func Run(c echo.Context) error {
	reqType := strings.ToLower(c.FormValue("type"))
	file, err := c.FormFile("src")
	if err != nil {
		c.JSON(500, "Error occured when handling source file")
	}

	if file == nil {
		c.JSON(400, "No file provided")
	}

	src, err := ReadRequestFile(file)
	if err != nil {
		c.JSON(500, "Error occured when handling source file")
	}

	switch reqType {
	case "zig":
		return c.JSON(200, "")

	case "c":
		start := time.Now()
		build := CRunner.Build(src)
		result := CRunner.Run(build.EncodedBinary)
		c.Response().Header().Set("Content-Type", "application/json")

		fmt.Printf("Request running time: %.4f\n", time.Since(start).Seconds())

		return c.JSON(200, result)
	}

	return nil
}

func ReadRequestFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	if file, err := io.ReadAll(file); err != nil {
		return "", err
	} else {
		return string(file), nil
	}
}
