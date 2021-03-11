package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverify"
	"github.com/gnames/gnverify/config"
	"github.com/gnames/gnverify/entity/output"
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
	e.GET("/data_sources", dataSources(gnv))
	e.GET("/data_sources/:id", dataSource(gnv))
	e.GET("/about", about())
	e.GET("/api", api())

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
	Page        string
	Input       string
	Format      string
	Preferred   []int
	Verified    []vlib.Verification
	DataSources []vlib.DataSource
	DataSource  vlib.DataSource
}

func about() func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Page: "about"}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func api() func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Page: "api"}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func dataSources(gnv gnverify.GNVerify) func(echo.Context) error {
	return func(c echo.Context) error {
		var err error
		data := Data{Page: "data_sources"}
		data.DataSources, err = gnv.DataSources()
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func dataSource(gnv gnverify.GNVerify) func(echo.Context) error {
	return func(c echo.Context) error {
		var err error
		data := Data{Page: "data_source"}
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}
		data.DataSource, err = gnv.DataSource(id)
		if err != nil {
			return fmt.Errorf("Cannot find DataSource for id '%s'", idStr)
		}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func home(gnv gnverify.GNVerify) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Page: "home", Format: "html"}
		var names []string

		params := c.QueryParams()
		data.Input = params.Get("names")
		data.Preferred = getPreferredSources(params["ds"])
		prefOnly := params.Get("preferred_only") == "on"

		format := params.Get("format")
		if format == "csv" || format == "json" {
			data.Format = format
		}

		if data.Input != "" {
			split := strings.Split(data.Input, "\n")
			names = make([]string, len(split))
			for i := range split {
				names[i] = strings.TrimSpace(split[i])
			}

			opts := []config.Option{config.OptPreferredSources(data.Preferred)}
			gnv.ChangeConfig(opts...)

			data.Verified = gnv.VerifyBatch(names)
			if prefOnly {
				for i := range data.Verified {
					data.Verified[i].BestResult = nil
				}
			}
		}

		switch data.Format {
		case "json":
			return c.JSON(http.StatusOK, data.Verified)
		case "csv":
			res := make([]string, len(data.Verified)+1)
			res[0] = output.CSVHeader()
			for i, v := range data.Verified {
				res[i+1] = output.Output(v, gnfmt.CSV, prefOnly)
			}
			return c.String(http.StatusOK, strings.Join(res, "\n"))
		default:
			return c.Render(http.StatusOK, "layout", data)
		}
	}
}

func getPreferredSources(ds []string) []int {
	var res []int
	if len(ds) == 0 {
		return res
	}

	for _, v := range ds {
		id, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		res = append(res, id)
	}
	return res
}
