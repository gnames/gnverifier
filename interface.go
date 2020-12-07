package gnverify

import (
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnverify/config"
)

type GNVerify interface {
	VerifyOne(name string) string
	VerifyStream(in <-chan []string, out chan []vlib.Verification)
	Config() config.Config
}
