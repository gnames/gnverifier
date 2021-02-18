package gnverify

import (
	"context"
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

// NewGNVerify constructs an object that implements GNVerify interface
// and can be used for matching strings to scientfic names.
func NewGNVerify(cnf config.Config) GNVerify {
	return &gnverify{
		config:   cnf,
		verifier: verifrest.NewVerifier(cnf.VerifierURL),
	}
}

// Config returns configuration data.
func (gnv *gnverify) Config() config.Config {
	return gnv.config
}

// VerifyOne verifies one input string and returns results
// as a string in JSON or CSV format.
func (gnv *gnverify) VerifyOne(name string) string {
	params := vlib.VerifyParams{
		NameStrings:      []string{name},
		PreferredSources: gnv.config.PreferredSources,
	}
	verif := gnv.verifier.Verify(context.Background(), params)
	if len(verif) < 1 {
		log.Fatalf("Did not get results from verifier")
	}
	return output.Output(verif[0], gnv.config.Format, gnv.config.PreferredOnly)
}

// VerifyStream receives batches of strings through the input
// channel and sends results of verification via output
// channel.
func (gnv *gnverify) VerifyStream(
	in <-chan []string,
	out chan []vlib.Verification,
) {
	var wg sync.WaitGroup
	wg.Add(gnv.Config().Jobs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vwChan := gnv.loadNames(ctx, in)

	for i := 0; i < gnv.Config().Jobs; i++ {
		go gnv.VerifyWorker(ctx, vwChan, out, &wg)
	}

	wg.Wait()
	close(out)
}

func (gnv *gnverify) VerifyWorker(
	ctx context.Context,
	in <-chan vlib.VerifyParams,
	out chan<- []vlib.Verification,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for params := range in {
		if len(params.NameStrings) == 0 {
			continue
		}
		verif := gnv.verifier.Verify(ctx, params)
		if len(verif) < 1 {
			log.Fatalf("Did not get results from verifier")
		}
		out <- verif
	}
}

func (gnv *gnverify) loadNames(
	ctx context.Context,
	inChan <-chan []string,
) <-chan vlib.VerifyParams {
	vwChan := make(chan vlib.VerifyParams)
	go func() {
		defer close(vwChan)
		for names := range inChan {

			params := vlib.VerifyParams{
				NameStrings:      names,
				PreferredSources: gnv.config.PreferredSources,
			}
			select {
			case <-ctx.Done():
				return
			case vwChan <- params:
			}
		}
	}()
	return vwChan
}
