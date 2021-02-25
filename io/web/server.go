package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnverify"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const withLogs = true

// Run starts the GNparser web service and servies both RESTful API and
// a website.
func Run(gnv gnverify.GNVerify, port int) {
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

type Data struct {
	Page     string
	Input    string
	Verified []vlib.Verification
}

func home(gnv gnverify.GNVerify) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Page: "home"}
		var names []string

		fmt.Printf("data: %#v\n", data)
		data.Input = c.QueryParam("names")
		if data.Input != "" {
			names = strings.Split(data.Input, "\n")
			fmt.Printf("names: %#v", names)
			data.Verified = gnv.VerifyBatch(names)
		}
		return c.Render(http.StatusOK, "layout", data)
	}
}
