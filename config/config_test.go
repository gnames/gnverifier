package config_test

import (
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnverifier/config"
	"github.com/stretchr/testify/assert"
)

var url = "https://gnames.globalnames.org"

func TestConfigDefault(t *testing.T) {
	cnf := config.New()
	deflt := config.Config{
		Format:      gnfmt.CSV,
		VerifierURL: "https://verifier.globalnames.org/api/v1/",
	}
	assert.Equal(t, deflt.Format, cnf.Format)
	assert.Equal(t, deflt.VerifierURL, cnf.VerifierURL)
}

func TestConfigOpts(t *testing.T) {
	opts := opts()
	cnf := config.New(opts...)
	updt := config.Config{
		Format:      gnfmt.PrettyJSON,
		DataSources: []int{1, 2, 3},
		VerifierURL: url,
	}
	assert.Equal(t, updt.Format, cnf.Format)
	assert.Equal(t, updt.DataSources, cnf.DataSources)
	assert.Equal(t, updt.VerifierURL, cnf.VerifierURL)
}

type formatTest struct {
	String string
	gnfmt.Format
}

func opts() []config.Option {
	return []config.Option{
		config.OptFormat(gnfmt.PrettyJSON),
		config.OptDataSources([]int{1, 2, 3}),
		config.OptVerifierURL(url),
	}
}
