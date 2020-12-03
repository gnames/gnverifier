package gnverify

import (
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/format"
)

type GNVerify interface {
	VerifyOne(name string) string
	VerifyStream(in <-chan []string, out chan []vlib.Verification)
	Format() format.Format
	PreferredOnly() bool
}
