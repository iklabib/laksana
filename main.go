package main

import (
	"io"
	"markisa/runner/zig"
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

func main() {
  e := echo.New()
  e.POST("/run", Run)

  e.Logger.Fatal(e.Start("127.0.0.1:5176"))
}

func Run(c echo.Context) error {
  reqType := c.FormValue("Type")
  reqArchive, err := c.FormFile("Archive")
  if err != nil {
    panic(err)
  }

  archive, err := ReadRequestFile(reqArchive)
  if err != nil {
    return err
  }

  switch reqType {
  case "zig":
    ar := string(archive)
    // TODO: return compile error & system error
    stdout, _, _ := zig.Run(ar)
    c.Response().Header().Set("Content-Type", "application/json")
    return c.String(200, stdout)
  }
  return nil
}

func ReadRequestFile(fileHeader *multipart.FileHeader) ([]byte, error) {
  file, err := fileHeader.Open()
  if err != nil {
    return nil, err
  }
  defer file.Close()
  
  if file, err := io.ReadAll(file); err != nil {
    return nil, err
  } else {
    return file, nil
  }
}
