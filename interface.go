package gnverify

import (
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverify/config"
)

// GNVerify is the use-case interface of the gnverify app. It determines
// methods needed to verify (reconcile/resolve) strings to scientific
// names.
type GNVerify interface {
	// VerifyOne takes a name-string and returns result of verification as a
	// JSON of CSV string.
	VerifyOne(name string) string
	// VerifyBatch takes a slice of names and verifies them all at once
	VerifyBatch(names []string) []vlib.Verification
	// VerifyStream receves batches of strings via one channel, verifies
	// the strings and sends results to another channel.
	VerifyStream(in <-chan []string, out chan []vlib.Verification)
	// ChangeConfig modifies configuration of GNVerify.
	ChangeConfig(opts ...config.Option)
	// Config returns  configuration data.
	Config() config.Config
}
