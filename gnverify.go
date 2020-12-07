package gnverify

import (
	"log"
	"sync"

	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnverify/config"
	"github.com/gnames/gnverify/entity/output"
	"github.com/gnames/gnverify/entity/verifier"
	"github.com/gnames/gnverify/io/verifrest"
)

type gnverify struct {
	config   config.Config
	verifier verifier.Verifier
}

// NewGNVerify creates a struct that implements GNVerify interface.
func NewGNVerify(cnf config.Config) GNVerify {
	return &gnverify{
		config:   cnf,
		verifier: verifrest.NewVerifier(cnf.VerifierURL),
	}
}

func (gnv *gnverify) Config() config.Config {
	return gnv.config
}

func (gnv *gnverify) VerifyOne(name string) string {
	params := vlib.VerifyParams{
		NameStrings:      []string{name},
		PreferredSources: gnv.config.PreferredSources,
	}
	verif := gnv.verifier.Verify(params)
	if len(verif) < 1 {
		log.Fatalf("Did not get results from verifier")
	}
	return output.Output(verif[0], gnv.config.Format, gnv.config.PreferredOnly)
}

func (gnv *gnverify) VerifyStream(in <-chan []string,
	out chan []vlib.Verification) {
	vwChan := make(chan vlib.VerifyParams)
	var wg sync.WaitGroup
	wg.Add(gnv.Config().Jobs)

	go func() {
		for names := range in {
			vwChan <- vlib.VerifyParams{
				NameStrings:      names,
				PreferredSources: gnv.config.PreferredSources,
			}
		}
		close(vwChan)
	}()
	for i := 0; i < gnv.Config().Jobs; i++ {
		go gnv.VerifyWorker(vwChan, out, &wg)
	}

	wg.Wait()
	close(out)
}

func (gnv *gnverify) VerifyWorker(in <-chan vlib.VerifyParams,
	out chan<- []vlib.Verification, wg *sync.WaitGroup) {
	defer wg.Done()
	for params := range in {
		if len(params.NameStrings) == 0 {
			continue
		}
		verif := gnv.verifier.Verify(params)
		if len(verif) < 1 {
			log.Fatalf("Did not get results from verifier")
		}
		out <- verif
	}
}
