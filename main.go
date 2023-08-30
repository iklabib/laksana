package main

import (
	"fmt"
	"io"
	"markisa/runner/zig"
	"mime/multipart"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var requestTimes = map[string]time.Time {}

func main() {
  if err := godotenv.Load(); err != nil {
    panic(err)
  }

  PORT := os.Getenv("PORT")
  BASE_URL := os.Getenv("BASE_URL")
  URL := BASE_URL + ":" + PORT

  e := echo.New()

  e.Use(middleware.Gzip())
  e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))

  e.GET("/", Root)
  e.POST("/run", Run)

  e.Logger.Fatal(e.Start(URL))
}

func Root(c echo.Context) error {
  getSession(&c)
  return c.NoContent(200)
}

func getSession(c *echo.Context) *sessions.Session {
  sess, _ := session.Get("markisa", *c)
  if !sess.IsNew {
    return sess
  }

  sess.Options = &sessions.Options{
    Path:     "/",
    MaxAge:   3600 * 24 * 7,
    HttpOnly: true,
  }

  sess.Values["visitorId"] = uuid.NewString()
  sess.Save((*c).Request(), (*c).Response())

  return sess
}

func Run(c echo.Context) error {
  visitorId := getSession(&c).Values["visitorId"].(string)
  requestTimes[visitorId] = time.Now()


  reqType := c.FormValue("type")
  reqArchive, err := c.FormFile("archive")
  if err != nil {
    return err
  }

  archive, err := ReadRequestFile(reqArchive)
  if err != nil {
    return err
  }

  switch reqType {
  case "zig":
    start := time.Now()
    ar := string(archive)
    result := zig.Run(ar)
    c.Response().Header().Set("Content-Type", "application/json")

    fmt.Printf("Request running time: %.4f\n", time.Since(start).Seconds())

    return c.JSON(200, result)
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
