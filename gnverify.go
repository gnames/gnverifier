package gnverify

import (
	"log"

	gne "github.com/gnames/gnames/domain/entity"
	gnusecase "github.com/gnames/gnames/domain/usecase"
	"github.com/gnames/gnverify/config"
	"github.com/gnames/gnverify/output"
	"github.com/gnames/gnverify/verifier"
)

type GNVerify struct {
	config.Config
	gnusecase.Verifier
}

func NewGNVerify(cnf config.Config) GNVerify {
	return GNVerify{
		Config:   cnf,
		Verifier: verifier.NewVerifierRest(cnf.VerifierURL),
	}
}

func (gnv GNVerify) Verify(name string) string {
	params := gne.VerifyParams{
		NameStrings:      []string{name},
		PreferredSources: gnv.Config.PreferredSources,
	}
	verif := gnv.Verifier.Verify(params)
	if len(verif) < 1 {
		log.Fatalf("Did not get results from verifier")
	}
	return output.Output(verif[0], gnv.Format, gnv.PreferredOnly)
}

func (gnv GNVerify) VerifyStream(in chan<- []string, out <-chan []*gne.Verification) {
}
