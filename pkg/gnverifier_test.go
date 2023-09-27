package gnverifier_test

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	gnverifier "github.com/gnames/gnverifier/pkg"
	"github.com/gnames/gnverifier/pkg/config"
	vtest "github.com/gnames/gnverifier/pkg/ent/verifier/verifiertesting"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestDataSources(t *testing.T) {
	dss := dataSources(t)
	vfr := new(vtest.FakeVerifier)
	vfr.DataSourcesReturns(dss, nil)
	cfg := config.New()
	gnv := gnverifier.New(cfg, vfr)
	res, err := gnv.DataSources()
	assert.Nil(t, err)
	assert.Equal(t, dss, res)

	vfr.DataSourcesReturns(nil, errors.New("fake error"))
	_, err = gnv.DataSources()
	assert.NotNil(t, err)
}

func TestDataSource(t *testing.T) {
	ds := dataSources(t)[0]
	vfr := new(vtest.FakeVerifier)
	vfr.DataSourceReturns(ds, nil)
	cfg := config.New()
	gnv := gnverifier.New(cfg, vfr)
	res, err := gnv.DataSource(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, res.ID)
	assert.Equal(t, "Catalogue of Life", res.Title)
}

func TestChangeConfig(t *testing.T) {
	vfr := new(vtest.FakeVerifier)
	cfg := config.New()
	gnv := gnverifier.New(cfg, vfr)
	assert.Equal(t, gnfmt.CSV, gnv.Config().Format)
	assert.Equal(t, 4, gnv.Config().Jobs)
	gnv = gnv.ChangeConfig(
		config.OptFormat(gnfmt.CompactJSON),
		config.OptJobs(10),
	)
	assert.Equal(t, gnfmt.CompactJSON, gnv.Config().Format)
	assert.Equal(t, 10, gnv.Config().Jobs)
}

func TestConfig(t *testing.T) {
	vfr := new(vtest.FakeVerifier)
	cfg := config.New()
	gnv := gnverifier.New(cfg, vfr)
	res := gnv.Config()
	assert.Equal(t, cfg, res)
}

func TestVerifyOne(t *testing.T) {
	verifs := verifications(t)
	vfr := new(vtest.FakeVerifier)
	cfg := config.New()
	gnv := gnverifier.New(cfg, vfr)

	vfr.VerifyReturns(verifs)
	res, err := gnv.VerifyOne("Pomatomus saltatrix (Linnaeus, 1766)")
	assert.Nil(t, err)
	assert.Equal(t, "Pomatomus saltatrix (Linnaeus, 1766)", res.Name)
	assert.NotNil(t, res.BestResult)
	assert.Equal(t, 1, vfr.VerifyCallCount())

	vfr.VerifyReturns(vlib.Output{})
	res, err = gnv.VerifyOne("something")
	assert.NotNil(t, err)
	assert.Equal(t, "", res.Name)
}

func TestVerifyBatch(t *testing.T) {
	verifs := verifications(t)
	vfr := new(vtest.FakeVerifier)
	cfg := config.New()
	gnv := gnverifier.New(cfg, vfr)

	vfr.VerifyReturns(verifs)
	batch := []string{
		"Pomatomus saltatrix (Linnaeus, 1766)",
		"Bubo bubo (Linnaeus, 1782)",
		"NotName",
	}
	res := gnv.VerifyBatch(context.Background(), batch)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, 1, vfr.VerifyCallCount())
}

func TestVerifyStream(t *testing.T) {
	verifs := verifications(t)
	vfr := new(vtest.FakeVerifier)
	cfg := config.New()
	gnv := gnverifier.New(cfg, vfr)

	vfr.VerifyReturns(verifs)
	batch := []string{
		"Pomatomus saltatrix (Linnaeus, 1766)",
		"Bubo bubo (Linnaeus, 1782)",
		"NotName",
	}

	chIn := make(chan []string)
	chOut := make(chan []vlib.Name)
	var wg sync.WaitGroup
	wg.Add(1)
	go gnv.VerifyStream(context.Background(), chIn, chOut)

	go func() {
		defer wg.Done()
		for res := range chOut {
			assert.Equal(t, 3, len(res))
		}
	}()

	var count int
	for count < 3 {
		chIn <- batch
		count++
	}
	close(chIn)
	wg.Wait()
	assert.Equal(t, 3, vfr.VerifyCallCount())
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

func verifications(t *testing.T) vlib.Output {
	c := cassette.New("dss")
	data, err := os.ReadFile("io/verifrest/fixtures/names.yaml")
	assert.Nil(t, err)
	err = yaml.Unmarshal(data, c)
	assert.Nil(t, err)
	dssStr := c.Interactions[0].Response.Body
	enc := gnfmt.GNjson{}
	var res vlib.Output
	err = enc.Decode([]byte(dssStr), &res)
	assert.Nil(t, err)
	return res
}
