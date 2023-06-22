package web

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnquery"
	"github.com/gnames/gnquery/ent/search"
	"github.com/gnames/gnuuid"
	"github.com/gnames/gnverifier/internal/ent/output"
	gnverifier "github.com/gnames/gnverifier/pkg"
	"github.com/gnames/gnverifier/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	zlog "github.com/rs/zerolog/log"
	nsqcfg "github.com/sfgrp/lognsq/config"
	"github.com/sfgrp/lognsq/ent/nsq"
	"github.com/sfgrp/lognsq/io/nsqio"
)

type formInput struct {
	Names          string `query:"names" form:"names"`
	Format         string `query:"format" form:"format"`
	AllMatches     string `query:"all_matches" form:"all_matches"`
	Capitalize     string `query:"capitalize" form:"capitalize"`
	SpeciesGroup   string `query:"species_group" form:"species_group"`
	FuzzyUninomial string `query:"fuzzy_uninomial" form:"fuzzy_uninomial"`
	DataSources    []int  `query:"ds" form:"ds"`
}

//go:embed static
var static embed.FS

// Run starts the GNparser web service and servies both RESTful API and
// a website.
func Run(gnv gnverifier.GNverifier, port int) {
	var err error
	e := echo.New()

	e.Use(middleware.Gzip())

	loggerNSQ := setLogger(e, gnv)
	if loggerNSQ != nil {
		defer loggerNSQ.Stop()
	}

	e.Renderer, err = NewTemplate()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/", homeGET(gnv))
	e.POST("/", homePOST(gnv))
	e.GET("/data_sources", dataSources(gnv))
	e.GET("/data_sources/:id", dataSource(gnv))
	e.GET("/name_strings/:id", nameString(gnv))
	e.GET("/name_strings/widget/:id", nameStringWidget(gnv))
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
	Page          string
	Input         string
	Format        string
	DataSourceIDs []int
	AllMatches    bool
	Verified      []vlib.Name
	DataSources   []vlib.DataSource
	DataSource    vlib.DataSource
	Version       string
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

func nameString(gnv gnverifier.GNverifier) func(echo.Context) error {
	return func(c echo.Context) error {
		data, err := getNameString(c, gnv)
		if err != nil {
			return err
		}

		switch data.Format {
		case "json":
			return c.JSON(http.StatusOK, data.Verified)
		case "csv":
			res := formatRows(data, gnfmt.CSV)
			return c.String(http.StatusOK, strings.Join(res, "\n"))
		case "tsv":
			res := formatRows(data, gnfmt.TSV)
			return c.String(http.StatusOK, strings.Join(res, "\n"))
		default:
			return c.Render(http.StatusOK, "layout", data)
		}
	}
}

func nameStringWidget(gnv gnverifier.GNverifier) func(echo.Context) error {
	return func(c echo.Context) error {
		data, err := getNameString(c, gnv)
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "name_string_widget", data)
	}
}

func getNameString(
	c echo.Context,
	gnv gnverifier.GNverifier,
) (Data, error) {
	var res Data
	id, _ := url.QueryUnescape(c.Param("id"))
	var ds []int
	var allMatches bool
	dsStr := c.QueryParam("data_sources")
	if dsStr != "" {
		dss := strings.Split(dsStr, ",")
		for i := range dss {
			num, err := strconv.Atoi(dss[i])
			if err == nil {
				ds = append(ds, num)
			}
		}
	}
	allMatches = c.QueryParam("all_matches") == "true"
	inp := vlib.NameStringInput{
		ID:             id,
		DataSources:    ds,
		WithAllMatches: allMatches,
	}
	out, err := gnv.NameString(inp)
	if err != nil {
		return res, err
	}
	var names []vlib.Name
	if out.Name != nil {
		names = []vlib.Name{*out.Name}
	}

	res = Data{
		Input:         inp.ID,
		Format:        "html",
		DataSourceIDs: inp.DataSources,
		AllMatches:    inp.WithAllMatches,
		Page:          "home",
		Verified:      names,
		Version:       gnv.GetVersion().Version,
	}
	format := c.QueryParam("format")
	if format == "csv" || format == "json" || format == "tsv" {
		res.Format = format
	}

	return res, nil
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

		return verificationResults(c, gnv, inp, data, "GET")
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

		return verificationResults(c, gnv, inp, data, "POST")
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
	caps := inp.Capitalize == "on"
	spGr := inp.SpeciesGroup == "on"
	fuzzyUni := inp.FuzzyUninomial == "on"
	all := inp.AllMatches == "on"
	q := make(url.Values)
	q.Set("names", inp.Names)
	q.Set("format", inp.Format)
	if caps {
		q.Set("capitalize", inp.Capitalize)
	}
	if fuzzyUni {
		q.Set("fuzzy_uninomial", inp.FuzzyUninomial)
	}
	if all {
		q.Set("all_matches", inp.AllMatches)
	}
	if spGr {
		q.Set("species_group", inp.SpeciesGroup)
	}
	for i := range inp.DataSources {
		q.Add("ds", strconv.Itoa(inp.DataSources[i]))
	}
	url := fmt.Sprintf("/?%s", q.Encode())
	return c.Redirect(http.StatusFound, url)
}

func verificationResults(
	c echo.Context,
	gnv gnverifier.GNverifier,
	inp *formInput,
	data Data,
	method string,
) error {
	var names []string
	caps := inp.Capitalize == "on"
	spGr := inp.SpeciesGroup == "on"
	fuzzyUni := inp.FuzzyUninomial == "on"
	data.AllMatches = inp.AllMatches == "on"

	data.Input = inp.Names

	data.DataSourceIDs = inp.DataSources

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
			config.OptDataSources(data.DataSourceIDs),
			config.OptWithCapitalization(caps),
			config.OptWithSpeciesGroup(spGr),
			config.OptWithUninomialFuzzyMatch(fuzzyUni),
			config.OptWithAllMatches(data.AllMatches),
		}
		gnv = gnv.ChangeConfig(opts...)

		if search.IsQuery(names[0]) {
			var err error
			inp := gnquery.New().Parse(names[0])
			if dss := gnv.Config().DataSources; len(dss) > 0 {
				inp.DataSources = dss
			}
			if all := gnv.Config().WithAllMatches; all {
				inp.WithAllMatches = all
			}
			data.Verified, err = gnv.Search(context.Background(), inp)
			if err != nil {
				log.Warn(err)
			}
			zlog.Info().
				Str("query", names[0]).
				Str("method", method).
				Int("verified", len(data.Verified)).
				Msg("Search")
			if len(data.Verified) == 0 {
				data.Verified = []vlib.Name{
					{
						ID:   gnuuid.New(inp.Query).String(),
						Name: inp.Query,
					},
				}
			}
		} else {
			data.Verified = gnv.VerifyBatch(context.Background(), names)

			if l := len(names); l > 0 {
				zlog.Info().
					Int("namesNum", len(names)).
					Str("example", names[0]).
					Str("method", method).
					Msg("Verification")
			}
		}
	}

	switch data.Format {
	case "json":
		return c.JSON(http.StatusOK, data.Verified)
	case "csv":
		res := formatRows(data, gnfmt.CSV)
		return c.String(http.StatusOK, strings.Join(res, "\n"))
	case "tsv":
		res := formatRows(data, gnfmt.TSV)
		return c.String(http.StatusOK, strings.Join(res, "\n"))
	default:
		return c.Render(http.StatusOK, "layout", data)
	}
}

func formatRows(data Data, f gnfmt.Format) []string {
	res := make([]string, len(data.Verified)+1)
	res[0] = output.CSVHeader(f)
	for i, v := range data.Verified {
		res[i+1] = output.NameOutput(v, f)
	}
	return res
}

func setLogger(e *echo.Echo, m gnverifier.GNverifier) nsq.NSQ {
	cfg := m.Config()
	nsqAddr := cfg.NsqdTCPAddress
	withLogs := cfg.WithWebLogs
	contains := cfg.NsqdContainsFilter
	regex := cfg.NsqdRegexFilter

	if nsqAddr != "" {
		cfg := nsqcfg.Config{
			StderrLogs: withLogs,
			Topic:      "gnverifier",
			Address:    nsqAddr,
			Contains:   contains,
			Regex:      regex,
		}
		remote, err := nsqio.New(cfg)
		logCfg := middleware.DefaultLoggerConfig
		if err == nil {
			logCfg.Output = remote
			zlog.Logger = zlog.Output(remote)
		}
		e.Use(middleware.LoggerWithConfig(logCfg))
		if err != nil {
			log.Warn(err)
		}
		return remote
	} else if withLogs {
		e.Use(middleware.Logger())
		return nil
	}
	return nil
}
