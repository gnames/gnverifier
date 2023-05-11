package verifier

import (
	"context"

	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnquery/ent/search"
)

//go:generate counterfeiter -o verifiertesting/fake_verifier.go . Verifier

// Verifier takes verification parameters and returns back results
//
//	of verification of name-strings.
type Verifier interface {
	// Verify takes a slice of strings to verify, optional preferred data-sources
	// and returns results of verification of the strings against known
	// scientific names.
	Verify(ctx context.Context, params vlib.Input) vlib.Output

	// NameString takes a name-string or its ID, as well as query parameters.
	// It returns results for this particular name-string.
	NameString(
		ctx context.Context,
		inp vlib.NameStringInput,
	) (vlib.NameStringOutput, error)

	DataSourcer
	Searcher
}

// DataSourcer provides information about available data-sources.
type DataSourcer interface {
	// DataSources returns meta-information about aggregated data-sources.
	DataSources(ctx context.Context) ([]vlib.DataSource, error)

	// DataSource returns meta-information about a particular data source.
	DataSource(ctx context.Context, id int) (vlib.DataSource, error)
}

// Searcher provides methods for faceted search.
type Searcher interface {
	// Search takes facets data (information about genus, species, author, year,
	// data-sources). And returns back names that match these components.
	Search(context.Context, search.Input) (search.Output, error)
}
