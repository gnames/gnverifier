package verifier

import (
	"context"

	vlib "github.com/gnames/gnlib/ent/verifier"
)

//go:generate counterfeiter -o verifiertesting/fake_verifier.go . Verifier

// Verifier takes verification parameters and returns back results
//  of verification of name-strings.
type Verifier interface {
	// Verify takes a slice of strings to verify, optional preferred data-sources
	// and returns results of verification of the strings against known
	// scientific names.
	Verify(ctx context.Context, params vlib.Input) vlib.Output

	// DataSources returns meta-information about aggregated data-sources.
	DataSources(ctx context.Context) ([]vlib.DataSource, error)

	// DataSource returns meta-information about a particular data source.
	DataSource(ctx context.Context, id int) (vlib.DataSource, error)
}
