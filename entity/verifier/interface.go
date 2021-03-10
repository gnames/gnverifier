package verifier

import (
	"context"

	vlib "github.com/gnames/gnlib/ent/verifier"
)

// Verifier takes verification parameters and returns back results
//  of verification of name-strings.
type Verifier interface {
	// Verify takes a slice of strings to verify, optional preferred data-sources
	// and returns results of verification of the strings against known
	// scientific names.
	Verify(ctx context.Context, params vlib.VerifyParams) []vlib.Verification
}
