package gnverify_test

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverify"
	"github.com/gnames/gnverify/config"
	vtest "github.com/gnames/gnverify/ent/verifier/verifiertesting"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestDataSources(t *testing.T) {
	dss := dataSources(t)
	vfr := new(vtest.FakeVerifier)
	vfr.DataSourcesReturns(dss, nil)
	cfg := config.New()
	gnv := gnverify.New(cfg, vfr)
	res, err := gnv.DataSources()
	assert.Nil(t, err)
	assert.Equal(t, res, dss)

	vfr.DataSourcesReturns(nil, errors.New("fake error"))
	res, err = gnv.DataSources()
	assert.NotNil(t, err)
}

func TestDataSource(t *testing.T) {
	ds := dataSources(t)[0]
	vfr := new(vtest.FakeVerifier)
	vfr.DataSourceReturns(ds, nil)
	cfg := config.New()
	gnv := gnverify.New(cfg, vfr)
	res, err := gnv.DataSource(1)
	assert.Nil(t, err)
	assert.Equal(t, res.ID, 1)
	assert.Equal(t, res.Title, "Catalogue of Life")
}

func TestChangeConfig(t *testing.T) {
	vfr := new(vtest.FakeVerifier)
	cfg := config.New()
	gnv := gnverify.New(cfg, vfr)
	assert.Equal(t, gnv.Config().Format, gnfmt.CSV)
	assert.Equal(t, gnv.Config().Jobs, 4)
	gnv.ChangeConfig(
		config.OptFormat(gnfmt.CompactJSON),
		config.OptJobs(10),
	)
	assert.Equal(t, gnv.Config().Format, gnfmt.CompactJSON)
	assert.Equal(t, gnv.Config().Jobs, 10)
}

func TestConfig(t *testing.T) {
	vfr := new(vtest.FakeVerifier)
	cfg := config.New()
	gnv := gnverify.New(cfg, vfr)
	res := gnv.Config()
	assert.Equal(t, res, cfg)
}

func TestVerifyOne(t *testing.T) {
	verifs := verifications(t)
	vfr := new(vtest.FakeVerifier)
	cfg := config.New()
	gnv := gnverify.New(cfg, vfr)

	vfr.VerifyReturns(verifs[0:1])
	res, err := gnv.VerifyOne("Pomatomus saltatrix (Linnaeus, 1766)")
	assert.Nil(t, err)
	assert.Equal(t, res.Input, "Pomatomus saltatrix (Linnaeus, 1766)")
	assert.NotNil(t, res.BestResult)
	assert.Equal(t, vfr.VerifyCallCount(), 1)

	vfr.VerifyReturns([]vlib.Verification{})
	res, err = gnv.VerifyOne("something")
	assert.NotNil(t, err)
	assert.Equal(t, res.Input, "")
}

func TestVerifyBatch(t *testing.T) {
	verifs := verifications(t)
	vfr := new(vtest.FakeVerifier)
	cfg := config.New()
	gnv := gnverify.New(cfg, vfr)

	vfr.VerifyReturns(verifs)
	batch := []string{
		"Pomatomus saltatrix (Linnaeus, 1766)",
		"Bubo bubo (Linnaeus, 1782)",
		"NotName",
	}
	res := gnv.VerifyBatch(batch)
	assert.Equal(t, len(res), 3)
	assert.Equal(t, vfr.VerifyCallCount(), 1)
}

func TestVerifyStream(t *testing.T) {
	verifs := verifications(t)
	vfr := new(vtest.FakeVerifier)
	cfg := config.New()
	gnv := gnverify.New(cfg, vfr)

	vfr.VerifyReturns(verifs)
	batch := []string{
		"Pomatomus saltatrix (Linnaeus, 1766)",
		"Bubo bubo (Linnaeus, 1782)",
		"NotName",
	}

	chIn := make(chan []string)
	chOut := make(chan []vlib.Verification)
	var wg sync.WaitGroup
	wg.Add(1)
	go gnv.VerifyStream(chIn, chOut)

	go func() {
		defer wg.Done()
		for res := range chOut {
			assert.Equal(t, len(res), 3)
		}
	}()

	var count int
	for count < 3 {
		chIn <- batch
		count++
	}
	close(chIn)
	wg.Wait()
	assert.Equal(t, vfr.VerifyCallCount(), 3)
}

func dataSources(t *testing.T) []vlib.DataSource {
	c := cassette.New("dss")
	data, err := os.ReadFile("io/verifrest/fixtures/dss.yaml")
	assert.Nil(t, err)
	err = yaml.Unmarshal(data, c)
	assert.Nil(t, err)
	dssStr := c.Interactions[0].Response.Body
	enc := gnfmt.GNjson{}
	res := make([]vlib.DataSource, 0)
	err = enc.Decode([]byte(dssStr), &res)
	assert.Nil(t, err)
	return res
}

func verifications(t *testing.T) []vlib.Verification {
	c := cassette.New("dss")
	data, err := os.ReadFile("io/verifrest/fixtures/names.yaml")
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
