package gnverifier

import (
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier/config"
)

// GNverifier is the use-case interface of the gnverifier app. It determines
// methods needed to verify (reconcile/resolve) strings to scientific
// names.
type GNverifier interface {
	// VerifyOne takes a name-string and returns the result of verification.
	VerifyOne(name string) (vlib.Verification, error)

	// VerifyBatch takes a slice of names and verifies them all at once
	VerifyBatch(names []string) []vlib.Verification

	// VerifyStream receves batches of strings via one channel, verifies
	// the strings and sends results to another channel.
	VerifyStream(in <-chan []string, out chan []vlib.Verification)

	// ChangeConfig modifies configuration of GNVerifier.
	ChangeConfig(opts ...config.Option)

	// Config returns  configuration data.
	Config() config.Config

	// DataSources returns information about Data Sources harvested for
	// verification.
	DataSources() ([]vlib.DataSource, error)

	// DataSource uses ID input to return meta-information about a particular
	// data-source.
	DataSource(id int) (vlib.DataSource, error)
}
