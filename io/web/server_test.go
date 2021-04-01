package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier"
	"github.com/gnames/gnverifier/config"
	vtest "github.com/gnames/gnverifier/ent/verifier/verifiertesting"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func handlerGET(path string, t *testing.T) (echo.Context, *httptest.ResponseRecorder) {
	var err error
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	e.Renderer, err = NewTemplate()
	assert.Nil(t, err)
	c := e.NewContext(req, rec)
	return c, rec
}

func TestAbout(t *testing.T) {
	c, rec := handlerGET("/about", t)

	assert.Nil(t, about()(c))
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Contains(t, rec.Body.String(), "Matching Process")
}

func TestAPI(t *testing.T) {
	c, rec := handlerGET("/api", t)

	assert.Nil(t, api()(c))
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Contains(t, rec.Body.String(), "OpenAPI Schema")
}

func TestHomeGET(t *testing.T) {
	c, rec := handlerGET("/", t)

	verifs := verifications(t)
	cfg := config.New()
	vfr := new(vtest.FakeVerifier)
	vfr.VerifyReturns(verifs)
	gnv := gnverifier.New(cfg, vfr)

	assert.Nil(t, homeGET(gnv)(c))
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Contains(t, rec.Body.String(), "Global Names Verifier")
	assert.Contains(t, rec.Body.String(), "Advanced Options")
}

func TestHomePOST(t *testing.T) {
	var err error
	verifs := verifications(t)
	f := make(url.Values)
	f.Set("names", "Bubo bubo\nPomatomus saltator\nNotName")
	f.Set("format", "html")

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(f.Encode()),
	)
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	e := echo.New()
	e.Renderer, err = NewTemplate()
	assert.Nil(t, err)
	c := e.NewContext(req, rec)

	cfg := config.New(config.OptNamesNumThreshold(2))
	vfr := new(vtest.FakeVerifier)
	vfr.VerifyReturns(verifs)
	gnv := gnverifier.New(cfg, vfr)
	assert.Nil(t, homePOST(gnv)(c))
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Contains(t, rec.Body.String(), "Bubo (genus)")
}

func TestHomePostGet(t *testing.T) {
	var err error
	verifs := verifications(t)
	f := make(url.Values)
	f.Set("names", "Bubo bubo\nPomatomus saltator\nNotName")
	f.Set("format", "html")

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(f.Encode()),
	)
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	e := echo.New()
	e.Renderer, err = NewTemplate()
	assert.Nil(t, err)
	c := e.NewContext(req, rec)

	cfg := config.New(config.OptNamesNumThreshold(20))
	vfr := new(vtest.FakeVerifier)
	vfr.VerifyReturns(verifs)
	gnv := gnverifier.New(cfg, vfr)
	assert.Nil(t, homePOST(gnv)(c))
	// redirect to GET
	assert.Equal(t, rec.Code, http.StatusFound)
	assert.NotContains(t, rec.Body.String(), "Bubo (genus)")
}

func verifications(t *testing.T) []vlib.Verification {
	c := cassette.New("dss")
	data, err := os.ReadFile("../verifrest/fixtures/names.yaml")
	assert.Nil(t, err)
	err = yaml.Unmarshal(data, c)
	assert.Nil(t, err)
	dssStr := c.Interactions[0].Response.Body
	enc := gnfmt.GNjson{}
	res := make([]vlib.Verification, 0)
	err = enc.Decode([]byte(dssStr), &res)
	assert.Nil(t, err)
	return res
}
