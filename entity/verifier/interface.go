package verifier

import vlib "github.com/gnames/gnlib/domain/entity/verifier"

// Verifier takes verification parameters and returns back results
//  of verification of name-strings.
type Verifier interface {
	// Verify takes a slice of strings to verify, optional preferred data-sources
	// and returns results of verification of the strings against known
	// scientific names.
	Verify(params vlib.VerifyParams) []vlib.Verification
}
