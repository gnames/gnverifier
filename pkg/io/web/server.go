package web

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnquery"
	"github.com/gnames/gnquery/ent/search"
	"github.com/gnames/gnuuid"
	gnverifier "github.com/gnames/gnverifier/pkg"
	"github.com/gnames/gnverifier/pkg/config"
	"github.com/gnames/gnverifier/pkg/ent/output"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type formInput struct {
	Names          string `query:"names" form:"names"`
	Vernaculars    string `query:"vernaculars" form:"vernaculars"`
	Format         string `query:"format" form:"format"`
	AllMatches     string `query:"all_matches" form:"all_matches"`
	Capitalize     string `query:"capitalize" form:"capitalize"`
	SpeciesGroup   string `query:"species_group" form:"species_group"`
	FuzzyUninomial string `query:"fuzzy_uninomial" form:"fuzzy_uninomial"`
	FuzzyRelaxed   string `query:"fuzzy_relaxed" form:"fuzzy_relaxed"`
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

	e.Renderer, err = NewTemplate()
	if err != nil {
		e.Logger.Fatal(err)
	}

	handle := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(handle.With(
		slog.String("gnApp", "gnmatcher"),
	))

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
	Vernaculars   []string
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
			split = split[:5_000]
		}
		if inp.Vernaculars != "" && len(split) > 50 {
			split = split[:50]
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
	fuzzyRel := inp.FuzzyRelaxed == "on"
	fuzzyUni := inp.FuzzyUninomial == "on"
	all := inp.AllMatches == "on"
	q := make(url.Values)
	q.Set("names", inp.Names)
	q.Set("format", inp.Format)
	if inp.Vernaculars != "" {
		q.Set("vernaculars", inp.Vernaculars)
	}
	if caps {
		q.Set("capitalize", inp.Capitalize)
	}
	if fuzzyRel {
		q.Set("fuzzy_relaxed", inp.FuzzyRelaxed)
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

// verificationResults processes name verification requests and returns results
// in the requested format (HTML, JSON, CSV, or TSV).
//
// The function handles both search queries and batch name verification:
//   - Search queries: Detected via search.IsQuery() for advanced search operations
//   - Batch verification: Processes multiple scientific names for verification
//
// Parameters:
//   - c: Echo context for HTTP request/response handling
//   - gnv: GNverifier instance for performing verifications
//   - inp: Form input containing user-submitted data and configuration
//   - data: Data struct that will be populated with verification results
//   - method: HTTP method used ("GET" or "POST") for logging purposes
//
// Returns an error if processing or response rendering fails.
func verificationResults(
	c echo.Context,
	gnv gnverifier.GNverifier,
	inp *formInput,
	data Data,
	method string,
) error {
	// Parse configuration options from form input
	opts := parseFormOptions(inp, &data)
	gnv = gnv.ChangeConfig(opts...)

	// Process names if input is provided
	if data.Input != "" {
		names := splitAndLimitNames(data.Input)

		// Handle search query vs. batch verification
		if len(names) > 0 && search.IsQuery(names[0]) {
			data.Verified = processSearchQuery(c.Request().Context(), gnv, names[0], method)
		} else {
			data.Verified = processBatchVerification(c.Request().Context(), gnv, names, method)
		}
	}

	// Return results in requested format
	return renderResults(c, data)
}

// parseFormOptions extracts configuration options from form input and updates data.
func parseFormOptions(inp *formInput, data *Data) []config.Option {
	// Parse boolean flags
	caps := inp.Capitalize == "on"
	spGr := inp.SpeciesGroup == "on"
	fuzzyUni := inp.FuzzyUninomial == "on"
	fuzzyRel := inp.FuzzyRelaxed == "on"
	data.AllMatches = inp.AllMatches == "on"

	// Set input data
	data.Input = inp.Names
	data.DataSourceIDs = inp.DataSources

	// Parse vernaculars (3-letter language codes)
	if inp.Vernaculars != "" {
		for v := range strings.SplitSeq(inp.Vernaculars, ",") {
			if len(v) == 3 {
				data.Vernaculars = append(data.Vernaculars, v)
			}
		}
	}

	// Set output format if valid
	if inp.Format == "csv" || inp.Format == "json" || inp.Format == "tsv" {
		data.Format = inp.Format
	}

	// Build configuration options
	return []config.Option{
		config.OptDataSources(data.DataSourceIDs),
		config.OptVernaculars(data.Vernaculars),
		config.OptWithCapitalization(caps),
		config.OptWithSpeciesGroup(spGr),
		config.OptWithRelaxedFuzzyMatch(fuzzyRel),
		config.OptWithUninomialFuzzyMatch(fuzzyUni),
		config.OptWithAllMatches(data.AllMatches),
	}
}

// splitAndLimitNames splits input by newlines, trims whitespace, and limits to 5,000 names.
func splitAndLimitNames(input string) []string {
	const maxNames = 5_000

	split := strings.Split(input, "\n")
	if len(split) > maxNames {
		split = split[:maxNames]
	}

	names := make([]string, 0, len(split))
	for i := range split {
		if name := strings.TrimSpace(split[i]); name != "" {
			names = append(names, name)
		}
	}
	return names
}

// processSearchQuery handles advanced search queries using gnquery syntax.
func processSearchQuery(
	ctx context.Context,
	gnv gnverifier.GNverifier,
	query string,
	method string,
) []vlib.Name {
	inp := gnquery.New().Parse(query)

	// Apply configuration from gnv
	if dss := gnv.Config().DataSources; len(dss) > 0 {
		inp.DataSources = dss
	}
	if all := gnv.Config().WithAllMatches; all {
		inp.WithAllMatches = all
	}

	verified, err := gnv.Search(ctx, inp)
	if err != nil {
		log.Warn(err)
	}

	slog.Info(
		"Search",
		"query", query,
		"method", method,
		"verified", len(verified),
	)

	// Return placeholder result if no matches found
	if len(verified) == 0 {
		return []vlib.Name{
			{
				ID:   gnuuid.New(inp.Query).String(),
				Name: inp.Query,
			},
		}
	}

	return verified
}

// processBatchVerification verifies a batch of scientific names.
func processBatchVerification(
	ctx context.Context,
	gnv gnverifier.GNverifier,
	names []string,
	method string,
) []vlib.Name {
	if len(names) == 0 {
		return nil
	}

	verified := gnv.VerifyBatch(ctx, names)

	slog.Info(
		"Verification",
		"namesNum", len(names),
		"example", names[0],
		"method", method,
	)

	return verified
}

// renderResults returns verification results in the requested format.
func renderResults(c echo.Context, data Data) error {
	switch data.Format {
	case "json":
		return c.JSON(http.StatusOK, data.Verified)
	case "csv":
		rows := formatRows(data, gnfmt.CSV)
		return c.String(http.StatusOK, strings.Join(rows, "\n"))
	case "tsv":
		rows := formatRows(data, gnfmt.TSV)
		return c.String(http.StatusOK, strings.Join(rows, "\n"))
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
