package config

import "github.com/gnames/gnverify/output"

// Config collects and stores external configuration data.
type Config struct {
	output.Format
	PreferredOnly    bool
	NameField        uint
	PreferredSources []int
	VerifierURL      string
}

// NewConfig is a Config constructor that takes external options to
// update default values to external ones.
func NewConfig(opts ...Option) Config {
	cnf := Config{
		Format:      output.CSV,
		NameField:   0,
		VerifierURL: "http://:8888",
	}
	for _, opt := range opts {
		opt(&cnf)
	}
	return cnf
}

// Option is a type of all options for Config.
type Option func(cnf *Config)

// OptFormat sets output format
func OptFormat(f output.Format) Option {
	return func(cnf *Config) {
		cnf.Format = f
	}
}

// OptPreferredOnly sets PreferredOnly field. If it is true output only
// contains results from preferred data-sources.
func OptPreferredOnly(b bool) Option {
	return func(cnf *Config) {
		cnf.PreferredOnly = b
	}
}

// OptNameField sets index of name in CSV file.
func OptNameField(i uint) Option {
	return func(cnf *Config) {
		cnf.NameField = i
	}
}

// OptPreferredSources set list of preferred sources.
func OptPreferredSources(srs []int) Option {
	return func(cnf *Config) {
		cnf.PreferredSources = srs
	}
}

// OptVerifierURL sets URL of the verification resource.
func OptVerifierURL(s string) Option {
	return func(cnf *Config) {
		cnf.VerifierURL = s
	}
}
