package config_test

import (
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnverify/config"
	"github.com/stretchr/testify/assert"
)

var url = "https://gnames.globalnames.org"

func TestConfigDefault(t *testing.T) {
	cnf := config.NewConfig()
	deflt := config.Config{
		Format:      gnfmt.CSV,
		VerifierURL: "https://verifier.globalnames.org/api/v1/",
	}
	assert.Equal(t, cnf.Format, deflt.Format)
	assert.Equal(t, cnf.VerifierURL, deflt.VerifierURL)
}

func TestConfigOpts(t *testing.T) {
	opts := opts()
	cnf := config.NewConfig(opts...)
	updt := config.Config{
		Format:           gnfmt.PrettyJSON,
		PreferredOnly:    true,
		PreferredSources: []int{1, 2, 3},
		VerifierURL:      url,
	}
	assert.Equal(t, cnf.Format, updt.Format)
	assert.Equal(t, cnf.PreferredOnly, updt.PreferredOnly)
	assert.Equal(t, cnf.PreferredSources, updt.PreferredSources)
	assert.Equal(t, cnf.VerifierURL, updt.VerifierURL)
}

type formatTest struct {
	String string
	gnfmt.Format
}

func opts() []config.Option {
	return []config.Option{
		config.OptFormat(gnfmt.PrettyJSON),
		config.OptPreferredOnly(true),
		config.OptPreferredSources([]int{1, 2, 3}),
		config.OptVerifierURL(url),
	}
}
