package main

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	routes "gitlab.com/iklabib/markisa/routes"
)

func main() {
	PORT := os.Getenv("PORT")
	BASE_URL := os.Getenv("BASE_URL")
	URL := BASE_URL + ":" + PORT

	e := echo.New()

	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	var sesskey = []byte(os.Getenv("SESSION_KEY"))
	e.Use(session.Middleware(sessions.NewCookieStore(sesskey)))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))

	e.POST("/run", routes.Run)

	e.GET("/share/:id", routes.GetFile)
	e.POST("/share", routes.ShareFile)

	e.Logger.Fatal(e.Start(URL))
}
