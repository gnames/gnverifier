package gnverifier

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gnames/gnlib/ent/gnvers"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier/config"
	"github.com/gnames/gnverifier/ent/verifier"
)

type gnverifier struct {
	config   config.Config
	verifier verifier.Verifier
}

// New constructs an object that implements GNVerifier interface
// and can be used for matching strings to scientfic names.
func New(cnf config.Config, vfr verifier.Verifier) GNverifier {
	return gnverifier{
		config:   cnf,
		verifier: vfr,
	}
}

// GetVersion returns version and build of GNverifier
func (gnv gnverifier) GetVersion() gnvers.Version {
	return gnvers.Version{Version: Version, Build: Build}
}

// DataSources returns meta-information about aggregated data-sources.
func (gnv gnverifier) DataSources() ([]vlib.DataSource, error) {
	return gnv.verifier.DataSources(context.Background())
}

// DataSource returns meta-information about a data-source found by its ID.
func (gnv gnverifier) DataSource(id int) (vlib.DataSource, error) {
	return gnv.verifier.DataSource(context.Background(), id)
}

// ChangeConfig modifies configuration.
func (gnv gnverifier) ChangeConfig(opts ...config.Option) GNverifier {
	for i := range opts {
		opts[i](&gnv.config)
	}
	return gnv
}

// Config returns configuration data.
func (gnv gnverifier) Config() config.Config {
	return gnv.config
}

// VerifyOne verifies one input string and returns results
// as a string in JSON or CSV format.
func (gnv gnverifier) VerifyOne(name string) (vlib.Verification, error) {
	params := gnv.setParams([]string{name})
	verif := gnv.verifier.Verify(context.Background(), params)
	if len(verif) < 1 {
		return vlib.Verification{}, errors.New("no verification results")
	}
	return verif[0], nil
}

// VerifyBatch takes a list of name-strings, verifies them and returns
// a batch of results back.
func (gnv gnverifier) VerifyBatch(nameStrings []string) []vlib.Verification {
	params := gnv.setParams(nameStrings)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	return gnv.verifier.Verify(ctx, params)
}

// VerifyStream receives batches of strings through the input
// channel and sends results of verification via output
// channel.
func (gnv gnverifier) VerifyStream(
	in <-chan []string,
	out chan []vlib.Verification,
) {
	var wg sync.WaitGroup
	wg.Add(gnv.config.Jobs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vwChan := gnv.loadNames(ctx, in)

	for i := 0; i < gnv.config.Jobs; i++ {
		go gnv.VerifyWorker(ctx, vwChan, out, &wg)
	}

	wg.Wait()
	close(out)
}

func (gnv gnverifier) VerifyWorker(
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

func (gnv gnverifier) loadNames(
	ctx context.Context,
	inChan <-chan []string,
) <-chan vlib.VerifyParams {
	vwChan := make(chan vlib.VerifyParams)
	go func() {
		defer close(vwChan)
		for names := range inChan {

			params := gnv.setParams(names)
			select {
			case <-ctx.Done():
				return
			case vwChan <- params:
			}
		}
	}()
	return vwChan
}

func (gnv gnverifier) setParams(names []string) vlib.VerifyParams {
	res := vlib.VerifyParams{
		NameStrings:        names,
		PreferredSources:   gnv.config.PreferredSources,
		WithCapitalization: gnv.config.WithCapitalization,
		WithAllMatches:     gnv.config.WithAllMatches,
	}
	return res
}
