package config

import (
	"regexp"

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

	// NsqdTCPAddress provides an address to the NSQ messenger TCP service. If
	// this value is set and valid, the web logs will be published to the NSQ.
	// The option is ignored if `Port` is not set.
	//
	// If WithWebLogs option is set to `false`, but `NsqdTCPAddress` is set to a
	// valid URL, the logs will be sent to the NSQ messanging service, but they
	// wil not appear as STRERR output.
	// Example: `127.0.0.1:4150`
	NsqdTCPAddress string

	// NsqdContainsFilter logs should match the filter to be sent to NSQ
	// service.
	// Examples:
	// "api" - logs should contain "api"
	// "!api" - logs should not contain "api"
	NsqdContainsFilter string

	// NsqdRegexFilter logs should match the regular expression to be sent to
	// NSQ service.
	// Example: `api\/v(0|1)`
	NsqdRegexFilter *regexp.Regexp

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

	// WithUninomialFuzzyMatch flag; when true, uninomial names are not
	// restricted from fuzzy matching. Normally it creates too many false
	// positives and is switched off.
	WithUninomialFuzzyMatch bool

	// WithWebLogs flag enables logs when running web-service. This flag is
	// ignored if `Port` value is not set.
	WithWebLogs bool
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

// OptNsqdContainsFilter provides a filter for logs sent to NSQ service.
func OptNsqdContainsFilter(s string) Option {
	return func(cfg *Config) {
		cfg.NsqdContainsFilter = s
	}
}

// OptNsqdRegexFilter provides a regular expression filter for
// logs sent to NSQ service.
func OptNsqdRegexFilter(s string) Option {
	return func(cfg *Config) {
		r := regexp.MustCompile(s)
		cfg.NsqdRegexFilter = r
	}
}

// OptNsqdTCPAddress provides a URL to NSQ messanging service.
func OptNsqdTCPAddress(s string) Option {
	return func(cfg *Config) {
		cfg.NsqdTCPAddress = s
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

func OptWithUninomialFuzzyMatch(b bool) Option {
	return func(cnf *Config) {
		cnf.WithUninomialFuzzyMatch = b
	}
}

// OptWithWebLogs sets the WithWebLogs field.
func OptWithWebLogs(b bool) Option {
	return func(cfg *Config) {
		cfg.WithWebLogs = b
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
