package config

import (
	"github.com/gnames/gnfmt"
)

// Config collects and stores external configuration data.
type Config struct {
	// Batch is the size of the string slices fed into input channel for
	// verification.
	Batch int

	// DataSources are IDs of DataSources that are important for
	// user. Normally only one "the best" reusult returns. If user gives
	// preferred sources, then matches from these sources are also
	// returned.
	DataSources []int

	// Format determins the output. It can be either JSON or CSV.
	Format gnfmt.Format

	// Jobs is the number of verification jobs to run in parallel.
	Jobs int

	// NamesNumThreshold the number of names after which POST gets redirected
	// to GET.
	NamesNumThreshold int

	// VerifierURL URL for gnames verification service. It only needs to
	// be changed if user sets local version of gnames.
	VerifierURL string

	// WithAllMatches flag; if true, results include all matches per source,
	// not only the best match.
	WithAllMatches bool

	// WithCapitalization flag; if true, the first rune of the name-string
	// will be capitalized when appropriate.
	WithCapitalization bool

	// WithSpeciesGroup flag; it is true, verification tries to search not only
	// for the given species name, but also for its species group. It means that
	// searching for "Aus bus" will also search for "Aus bus bus" and vice versa.
	// This function reflects existence of autononyms in botanical code, and
	// coordinated names in zoological code.
	WithSpeciesGroup bool

	// WithRelaxedFuzzyMatch flag; when true, relaxes fuzzy matching rules.
	// This increases recall and decreases precision. It also changes the
	// maximum number of names sent to fuzzy match from 10,000 to 50, because
	// it is assumed that uses has to check every name in the result.
	// It changes fuzzy match edit distance from 1 to 2 and it makes fuzzy
	// matching match slower.
	WithRelaxedFuzzyMatch bool

	// WithUninomialFuzzyMatch flag; when true, uninomial names are not
	// restricted from fuzzy matching. Normally it creates too many false
	// positives and is switched off.
	WithUninomialFuzzyMatch bool
}

// Option is a type of all options for Config.
type Option func(cnf *Config)

// OptDataSources set list of preferred sources.
func OptDataSources(srs []int) Option {
	return func(cnf *Config) {
		cnf.DataSources = srs
	}
}

// OptFormat sets output format
func OptFormat(f gnfmt.Format) Option {
	return func(cnf *Config) {
		cnf.Format = f
	}
}

// OptJobs sets number of jobs to run in parallel.
func OptJobs(i int) Option {
	return func(cnf *Config) {
		cnf.Jobs = i
	}
}

// OptNamesNumThreshold sets number of names after which there is no redirect
// from POST to GET.
func OptNamesNumThreshold(i int) Option {
	return func(cnf *Config) {
		cnf.NamesNumThreshold = i
	}
}

// OptVerifierURL sets URL of the verification resource.
func OptVerifierURL(s string) Option {
	return func(cnf *Config) {
		cnf.VerifierURL = s
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

// OptWithSpeciesGroup sets WithSpeciesGroup field.
func OptWithSpeciesGroup(b bool) Option {
	return func(cnf *Config) {
		cnf.WithSpeciesGroup = b
	}
}

func OptWithRelaxedFuzzyMatch(b bool) Option {
	return func(cnf *Config) {
		cnf.WithRelaxedFuzzyMatch = b
	}
}

func OptWithUninomialFuzzyMatch(b bool) Option {
	return func(cnf *Config) {
		cnf.WithUninomialFuzzyMatch = b
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
