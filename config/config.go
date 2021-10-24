package config

import (
	"github.com/gnames/gnfmt"
)

// Config collects and stores external configuration data.
type Config struct {
	// Format determins the output. It can be either JSON or CSV.
	Format gnfmt.Format

	// PreferredOnly hides BestResult if the user wants to see only
	// preferred results.
	PreferredOnly bool

	// PreferredSources are IDs of DataSources that are important for
	// user. Normally only one "the best" reusult returns. If user gives
	// preferred sources, then matches from these sources are also
	// returned.
	PreferredSources []int

	// WithAllSources flag; if true, best results from all matched sources are
	// returned. If this flag is true, PreferredSources are ignored.
	WithAllSources bool

	// WithAllMatches flag; if true, results include all matches per source,
	// not only the best match.
	WithAllMatches bool

	// WithCapitalization flag; if true, the first rune of the name-string
	// will be capitalized when appropriate.
	WithCapitalization bool

	// VerifierURL URL for gnames verification service. It only needs to
	// be changed if user sets local version of gnames.
	VerifierURL string

	// Jobs is the number of verification jobs to run in parallel.
	Jobs int

	// Batch is the size of the string slices fed into input channel for
	// verification.
	Batch int

	// NamesNumThreshold the number of names after which POST gets redirected
	// to GET.
	NamesNumThreshold int
}

// Option is a type of all options for Config.
type Option func(cnf *Config)

// OptFormat sets output format
func OptFormat(f gnfmt.Format) Option {
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

// OptWithAllSources sets WithAllSources flag.
func OptWithAllSources(b bool) Option {
	return func(cnf *Config) {
		cnf.WithAllSources = b
	}
}

// OptWithAllMatches sets WithAllMatches flag.
func OptWithAllMatches(b bool) Option {
	return func(cnf *Config) {
		cnf.WithAllMatches = b
	}
}

// OptWithCapitalization sets WithCapitalization field.
func OptWithCapitalization(b bool) Option {
	return func(cnf *Config) {
		cnf.WithCapitalization = b
	}
}

//OptJobs sets number of jobs to run in parallel.
func OptJobs(i int) Option {
	return func(cnf *Config) {
		cnf.Jobs = i
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

// OptNamesNumThreshold sets number of names after which there is a redirect
// from POST to GET.
func OptNamesNumThreshold(i int) Option {
	return func(cnf *Config) {
		cnf.NamesNumThreshold = i
	}
}

// New is a Config constructor that takes external options to
// update default values to external ones.
func New(opts ...Option) Config {
	cnf := Config{
		Format:            gnfmt.CSV,
		VerifierURL:       "https://verifier.globalnames.org/api/v1/",
		Batch:             5000,
		Jobs:              4,
		NamesNumThreshold: 20,
	}
	for _, opt := range opts {
		opt(&cnf)
	}
	return cnf
}
