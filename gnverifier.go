package gnverifier

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gnames/gnlib/ent/gnvers"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnquery/ent/search"
	"github.com/gnames/gnverifier/config"
	"github.com/gnames/gnverifier/ent/verifier"
)

type gnverifier struct {
	cfg      config.Config
	verifier verifier.Verifier
}

// New constructs an object that implements GNVerifier interface
// and can be used for matching strings to scientfic names.
func New(
	cfg config.Config,
	vfr verifier.Verifier,
) GNverifier {
	return gnverifier{
		cfg:      cfg,
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
		opts[i](&gnv.cfg)
	}
	return gnv
}

// Config returns configuration data.
func (gnv gnverifier) Config() config.Config {
	return gnv.cfg
}

// VerifyOne verifies one input string and returns results
// as a string in JSON or CSV format.
func (gnv gnverifier) VerifyOne(name string) (vlib.Name, error) {
	params := gnv.setParams([]string{name})
	verif := gnv.verifier.Verify(context.Background(), params)
	if len(verif.Names) < 1 {
		return vlib.Name{}, errors.New("no verification results")
	}
	return verif.Names[0], nil
}

// VerifyBatch takes a list of name-strings, verifies them and returns
// a batch of results back.
func (gnv gnverifier) VerifyBatch(
	ctx context.Context,
	nameStrings []string,
) []vlib.Name {
	params := gnv.setParams(nameStrings)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	return gnv.verifier.Verify(ctx, params).Names
}

// VerifyStream receives batches of strings through the input
// channel and sends results of verification via output
// channel.
func (gnv gnverifier) VerifyStream(
	ctx context.Context,
	in <-chan []string,
	out chan []vlib.Name,
) {
	var wg sync.WaitGroup
	wg.Add(gnv.cfg.Jobs)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	vwChan := gnv.loadNames(ctx, in)

	for i := 0; i < gnv.cfg.Jobs; i++ {
		go gnv.verifyWorker(ctx, vwChan, out, &wg)
	}

	wg.Wait()
	close(out)
}

func (gnv gnverifier) verifyWorker(
	ctx context.Context,
	in <-chan vlib.Input,
	out chan<- []vlib.Name,
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
		if len(verif.Names) < 1 {
			log.Fatalf("Did not get results from verifier")
		}
		out <- verif.Names
	}
}

func (gnv gnverifier) Search(
	ctx context.Context,
	inp search.Input,
) ([]vlib.Name, error) {
	res, err := gnv.verifier.Search(ctx, inp)
	return res.Names, err
}

func (gnv gnverifier) loadNames(
	ctx context.Context,
	inChan <-chan []string,
) <-chan vlib.Input {
	vwChan := make(chan vlib.Input)
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

func (gnv gnverifier) setParams(names []string) vlib.Input {
	res := vlib.Input{
		NameStrings:        names,
		DataSources:        gnv.cfg.DataSources,
		WithCapitalization: gnv.cfg.WithCapitalization,
		WithAllMatches:     gnv.cfg.WithAllMatches,
	}
	return res
}
