package gnverify

import (
	"log"
	"sync"

	gne "github.com/gnames/gnames/domain/entity"
	gnusecase "github.com/gnames/gnames/domain/usecase"
	"github.com/gnames/gnverify/config"
	"github.com/gnames/gnverify/output"
	"github.com/gnames/gnverify/verifier"
)

type GNVerify struct {
	config.Config
	gnusecase.Verifier
	Jobs int
}

func NewGNVerify(cnf config.Config) GNVerify {
	return GNVerify{
		Config:   cnf,
		Verifier: verifier.NewVerifierRest(cnf.VerifierURL),
		Jobs:     4,
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

func (gnv GNVerify) VerifyStream(in <-chan []string, out chan []*gne.Verification) {
	vwChan := make(chan gne.VerifyParams)
	var wg sync.WaitGroup
	wg.Add(gnv.Jobs)

	go func() {
		for names := range in {
			vwChan <- gne.VerifyParams{
				NameStrings:      names,
				PreferredSources: gnv.PreferredSources,
			}
		}
		close(vwChan)
	}()
	for i := 0; i < gnv.Jobs; i++ {
		go gnv.VerifyWorker(vwChan, out, &wg)
	}

	wg.Wait()
	close(out)
}

func (gnv GNVerify) VerifyWorker(in <-chan gne.VerifyParams, out chan<- []*gne.Verification,
	wg *sync.WaitGroup) {
	defer wg.Done()
	gnv = NewGNVerify(gnv.Config)
	for params := range in {
		verif := gnv.Verifier.Verify(params)
		if len(verif) < 1 {
			log.Fatalf("Did not get results from verifier")
		}
		out <- verif
	}
}
