package web

import (
	"embed"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier"
	"github.com/gnames/gnverifier/config"
	"github.com/gnames/gnverifier/ent/output"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const withLogs = false

type formInput struct {
	Names         string `query:"names" form:"names"`
	Format        string `query:"format" form:"format"`
	PreferredOnly string `query:"preferred_only" form:"preferred_only"`
	AllSources    string `query:"all_sources" form:"all_sources"`
	AllMatches    string `query:"all_matches" form:"all_matches"`
	Capitalize    string `query:"capitalize" form:"capitalize"`
	DS            []int  `query:"ds" form:"ds"`
}

//go:embed static
var static embed.FS

// Run starts the GNparser web service and servies both RESTful API and
// a website.
func Run(gnv gnverifier.GNverifier, port int) {
	var err error
	e := echo.New()

	e.Use(middleware.Gzip())
	if withLogs {
		e.Use(middleware.Logger())
	}

	e.Renderer, err = NewTemplate()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/", homeGET(gnv))
	e.POST("/", homePOST(gnv))
	e.GET("/data_sources", dataSources(gnv))
	e.GET("/data_sources/:id", dataSource(gnv))
	e.GET("/about", about(gnv))
	e.GET("/api", api(gnv))

	fs := http.FileServer(http.FS(static))
	e.GET("/static/*", echo.WrapHandler(fs))

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
	AllMatches  bool
	Verified    []vlib.Name
	DataSources []vlib.DataSource
	DataSource  vlib.DataSource
	Version     string
}

func about(gnv gnverifier.GNverifier) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Page: "about", Version: gnv.GetVersion().Version}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func api(gnv gnverifier.GNverifier) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Page: "api", Version: gnv.GetVersion().Version}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func dataSources(gnv gnverifier.GNverifier) func(echo.Context) error {
	return func(c echo.Context) error {
		var err error
		data := Data{Page: "data_sources", Version: gnv.GetVersion().Version}
		data.DataSources, err = gnv.DataSources()
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func dataSource(gnv gnverifier.GNverifier) func(echo.Context) error {
	return func(c echo.Context) error {
		var err error
		data := Data{Page: "data_source", Version: gnv.GetVersion().Version}
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}
		data.DataSource, err = gnv.DataSource(id)
		if err != nil {
			return fmt.Errorf("cannot find DataSource for id '%s'", idStr)
		}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func homeGET(gnv gnverifier.GNverifier) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Page: "home", Format: "html", Version: gnv.GetVersion().Version}

		inp := new(formInput)
		err := c.Bind(inp)
		if err != nil {
			return err
		}

		if strings.TrimSpace(inp.Names) == "" {
			return c.Render(http.StatusOK, "layout", data)
		}

		return verificationResults(c, gnv, inp, data)
	}
}

func homePOST(gnv gnverifier.GNverifier) func(echo.Context) error {
	return func(c echo.Context) error {
		inp := new(formInput)
		data := Data{Page: "home", Format: "html", Version: gnv.GetVersion().Version}

		err := c.Bind(inp)
		if err != nil {
			return err
		}

		if strings.TrimSpace(inp.Names) == "" {
			return c.Redirect(http.StatusFound, "")
		}

		split := strings.Split(inp.Names, "\n")
		if len(split) > 5_000 {
			split = split[0:5_000]
		}

		if len(split) < gnv.Config().NamesNumThreshold {
			return redirectToHomeGET(c, inp)
		}

		return verificationResults(c, gnv, inp, data)
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

func redirectToHomeGET(c echo.Context, inp *formInput) error {
	prefOnly := inp.PreferredOnly == "on"
	caps := inp.Capitalize == "on"
	q := make(url.Values)
	q.Set("names", inp.Names)
	q.Set("format", inp.Format)
	q.Set("all_sources", inp.AllSources)
	q.Set("all_matches", inp.AllMatches)
	if prefOnly {
		q.Set("preferred_only", inp.PreferredOnly)
	}
	if caps {
		q.Set("capitalize", inp.Capitalize)
	}
	for i := range inp.DS {
		q.Add("ds", strconv.Itoa(inp.DS[i]))
	}
	url := fmt.Sprintf("/?%s", q.Encode())
	return c.Redirect(http.StatusFound, url)
}

func verificationResults(
	c echo.Context,
	gnv gnverifier.GNverifier,
	inp *formInput,
	data Data,
) error {
	var names []string
	prefOnly := inp.PreferredOnly == "on"
	caps := inp.Capitalize == "on"
	data.AllMatches = inp.AllMatches == "on"

	data.Input = inp.Names

	data.Preferred = inp.DS
	if inp.AllSources == "on" {
		data.Preferred = []int{0}
	}

	format := inp.Format
	if format == "csv" || format == "json" || format == "tsv" {
		data.Format = format
	}

	if data.Input != "" {
		split := strings.Split(data.Input, "\n")
		if len(split) > 5_000 {
			split = split[0:5_000]
		}

		names = make([]string, 0, len(split))
		for i := range split {
			name := strings.TrimSpace(split[i])
			if name != "" {
				names = append(names, name)
			}
		}

		opts := []config.Option{
			config.OptDataSources(data.Preferred),
			config.OptWithCapitalization(caps),
			config.OptWithAllMatches(data.AllMatches),
		}
		gnv = gnv.ChangeConfig(opts...)

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
		res := formatRows(data, prefOnly, gnfmt.CSV)
		return c.String(http.StatusOK, strings.Join(res, "\n"))
	case "tsv":
		res := formatRows(data, prefOnly, gnfmt.TSV)
		return c.String(http.StatusOK, strings.Join(res, "\n"))
	default:
		return c.Render(http.StatusOK, "layout", data)
	}
}

func formatRows(data Data, prefOnly bool, f gnfmt.Format) []string {
	res := make([]string, len(data.Verified)+1)
	res[0] = output.CSVHeader(f)
	for i, v := range data.Verified {
		res[i+1] = output.NameOutput(v, f, prefOnly)
	}
	return res
}
