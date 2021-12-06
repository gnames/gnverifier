package verifrest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnuuid"
	"github.com/gnames/gnverifier/ent/verifier"
	log "github.com/sirupsen/logrus"
)

type verifrest struct {
	verifierURL string
	client      *http.Client
}

// New returns object that implements Verifier interface.
func New(url string) verifier.Verifier {
	if url[len(url)-1] != '/' {
		url = url + "/"
	}
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Timeout: 4 * time.Minute, Transport: tr}
	return &verifrest{verifierURL: url, client: client}
}

// DataSources returns meta-data about aggregated data-sources.
func (vr verifrest) DataSources(
	ctx context.Context,
) ([]vlib.DataSource, error) {
	enc := gnfmt.GNjson{}
	url := vr.verifierURL + "data_sources"
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Warnf("Cannot create request: %v", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := vr.client.Do(request)
	if err != nil {
		log.Warn("Cannot get data-sources information.")
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("Body reading is failing for data-sources.")
		return nil, err
	}
	response := make([]vlib.DataSource, 0)
	err = enc.Decode(respBytes, &response)
	if err != nil {
		log.Warnf("Cannot decode data-sources")
		return nil, err
	}
	return response, nil
}

// DataSource returns meta-data about a data-source found by ID.
func (vr verifrest) DataSource(
	ctx context.Context,
	id int,
) (vlib.DataSource, error) {
	response := vlib.DataSource{}
	enc := gnfmt.GNjson{}
	url := fmt.Sprintf("%sdata_sources/%d", vr.verifierURL, id)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Warnf("Cannot create request: %v", err)
		return response, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := vr.client.Do(request)
	if err != nil {
		log.Warn("Cannot get data-sources information.")
		return response, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("Body reading is failing for data-sources.")
		return response, err
	}
	err = enc.Decode(respBytes, &response)
	if err != nil {
		log.Warnf("Cannot decode data-sources")
		return response, err
	}
	return response, nil
}

// Verify takes names-strings and options and returns verification result.
func (vr verifrest) Verify(
	ctx context.Context,
	input vlib.Input,
) vlib.Output {
	var attempts int
	var response vlib.Output
	enc := gnfmt.GNjson{}
	paramsData, err := enc.Encode(input)
	if err != nil {
		log.Printf("Cannot encode names for verification: %s.", err)
	}

	attempts, err = try(func(int) (bool, error) {
		var cancel func()
		ctx, cancel = context.WithCancel(ctx)
		// client has Timeout, meaning cancel will propagate to the server after
		// the set time.
		defer cancel()

		d := bytes.NewReader(paramsData)
		namesRange := fmt.Sprintf(
			"%s-%s",
			input.NameStrings[0],
			input.NameStrings[len(input.NameStrings)-1],
		)

		url := vr.verifierURL + "verifications"
		var request *http.Request
		var resp *http.Response
		var respBytes []byte
		request, err = http.NewRequestWithContext(ctx, http.MethodPost, url, d)
		if err != nil {
			log.Fatalf("Cannot create request: %v", err)
		}
		request.Header.Set("Content-Type", "application/json")

		resp, err = vr.client.Do(request)
		if err != nil {
			log.Warnf("Request is failing for %s.", namesRange)
			return true, err
		}
		defer resp.Body.Close()

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Warnf("Body reading is failing for %s.", namesRange)
			return true, err
		}
		err = enc.Decode(respBytes, &response)
		if err != nil {
			log.Warnf("Response decoding is failing for %s.", namesRange)
			return true, err
		}
		return false, nil
	})

	if err != nil {
		log.Printf("Verification failed for %s-%s after %d attempts.", input.NameStrings[0],
			input.NameStrings[len(input.NameStrings)-1], attempts)
		res := vlib.Output{Names: make([]vlib.Name, len(input.NameStrings))}
		for i := range input.NameStrings {
			name := input.NameStrings[i]
			res.Names[i] = vlib.Name{
				ID:    gnuuid.New(name).String(),
				Name:  name,
				Error: err.Error(),
			}
		}
		return res
	}
	return response
}

func try(fn func(int) (bool, error)) (int, error) {
	var (
		err        error
		tryAgain   bool
		maxRetries = 3
		attempt    = 1
	)
	for {
		tryAgain, err = fn(attempt)
		if !tryAgain || err == nil {
			break
		}
		attempt++
		if attempt > maxRetries {
			return maxRetries, err
		}
		time.Sleep(200 * time.Millisecond)
	}
	return attempt, err
}
