package verifrest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/encode"
	"github.com/gnames/gnverify/entity/verifier"
	log "github.com/sirupsen/logrus"
)

type verifrest struct {
	verifierURL string
	client      *http.Client
}

// NewVerifier returns object that implements Verifier interface.
func NewVerifier(url string) verifier.Verifier {
	if url[len(url)-1] != '/' {
		url = url + "/"
	}
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: tr}
	return &verifrest{verifierURL: url, client: client}
}

// Verify takes names-strings and options and returns verification result.
func (vr *verifrest) Verify(params vlib.VerifyParams) []vlib.Verification {
	var attempts int
	var response []vlib.Verification
	enc := encode.GNjson{}
	req, err := enc.Encode(params)
	if err != nil {
		log.Printf("Cannot encode names for verification: %s.", err)
	}

	attempts, err = try(func(int) (bool, error) {
		r := bytes.NewReader(req)
		namesRange := fmt.Sprintf("%s-%s", params.NameStrings[0], params.NameStrings[len(params.NameStrings)-1])

		resp, err := vr.client.Post(vr.verifierURL+"verifications", "application/json", r)
		if err != nil {
			log.Warnf("Request is failing for %s.", namesRange)
			return true, err
		}
		defer resp.Body.Close()

		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warnf("Body reading is failing for %s.", namesRange)
			return true, err
		}
		response = make([]vlib.Verification, 0)
		err = enc.Decode(respBytes, &response)
		if err != nil {
			log.Warnf("Response decoding is failing for %s.", namesRange)
			return true, err
		}
		return false, nil
	})

	if err != nil {
		log.Printf("Verification failed for %s-%s after %d attempts.", params.NameStrings[0],
			params.NameStrings[len(params.NameStrings)-1], attempts)
		log.Fatal(err)
	}
	return response
}

func try(fn func(int) (bool, error)) (int, error) {
	var (
		err        error
		tryAgain   bool
		maxRetries int = 5
		attempt    int = 1
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
