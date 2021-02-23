package web

import (
	"fmt"
	"net/http"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gnames/gnverify"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const withLogs = false

// Run starts the GNparser web service and servies both RESTful API and
// a website.
func Run(gnv *gnverify.GNVerify, port int) {
	var err error
	e := echo.New()

	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	e.Use(middleware.CSRF())
	if withLogs {
		e.Use(middleware.Logger())
	}

	e.Renderer, err = NewTemplate()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/", home(gnv))

	assetHandler := http.FileServer(rice.MustFindBox("assets").HTTPBox())
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}

func home(gnv *gnverify.GNVerify) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "layout.html", map[string]interface{}{
			"name": "Dolly!",
		})
	}
}
