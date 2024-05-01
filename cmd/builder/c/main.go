package main

// "instance" received a compiled program and run it.
// The program must be complete within 10 seconds, else be killed.
// "instance" supposed to be the entry point for container

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/iklabib/markisa/model"
	"github.com/iklabib/markisa/util"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.POST("/", func(c echo.Context) error {
		if c.Request().Body == nil {
			return c.JSON(400, "No file provided")
		}

		file, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(500, "Error occured when handling source file")
		}

		build := Build(string(file))
		return c.JSON(200, build)
	})

	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}

func Build(source string) model.BuildResult {
	dir := "/tmp/c"
	srcPath := filepath.Join(dir, "prog.c")
	src, err := os.Create(srcPath)
	if err != nil {
		return internalError()
	}
	src.WriteString(source)
	src.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("clang", "prog.c", "-o", "prog")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = dir

	err = cmd.Run()

	buildResult := model.BuildResult{
		ExitCode: util.GetExitCode(&err),
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Status:   "BUILD_SUCCESS",
	}

	if err != nil || buildResult.ExitCode != 0 {
		buildResult.Status = "BUILD_ERROR"
		return buildResult
	}

	prog, err := os.ReadFile(filepath.Join(dir, "prog"))
	if err != nil {
		return internalError()
	}

	// encode binary as ascii85 before get jsonified
	buildResult.Executable = prog

	return buildResult
}

func internalError() model.BuildResult {
	return model.BuildResult{
		ExitCode: -1,
		Status:   "INTERNAL_ERROR",
	}
}
