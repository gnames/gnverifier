package gnverifier

import (
	"context"

	"github.com/gnames/gnlib/ent/gnvers"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnquery/ent/search"
	"github.com/gnames/gnverifier/pkg/config"
)

// GNverifier is the use-case interface of the gnverifier app. It determines
// methods needed to verify (reconcile/resolve) strings to scientific
// names.
type GNverifier interface {
	// VerifyOne takes a name-string and returns the result of verification.
	VerifyOne(name string) (vlib.Name, error)

	// VerifyBatch takes a slice of names and verifies them all at once
	VerifyBatch(ctx context.Context, names []string) []vlib.Name

	// VerifyStream receves batches of strings via one channel, verifies
	// the strings and sends results to another channel.
	VerifyStream(ctx context.Context, in <-chan []string, out chan []vlib.Name)

	// Search provides faceted search functionality.
	Search(ctx context.Context, srch search.Input) ([]vlib.Name, error)

	// ChangeConfig modifies configuration of GNverifier.
	ChangeConfig(opts ...config.Option) GNverifier

	// Config returns  configuration data.
	Config() config.Config

	// DataSources returns information about Data Sources harvested for
	// verification.
	DataSources() ([]vlib.DataSource, error)

	// DataSource uses ID input to return meta-information about a particular
	// data-source.
	DataSource(id int) (vlib.DataSource, error)

	// GetVersion returns version of the gnverifier
	GetVersion() gnvers.Version
}
