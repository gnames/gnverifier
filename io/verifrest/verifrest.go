package verifrest

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/encode"
	"github.com/gnames/gnverify/entity/verifier"
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
	enc := encode.GNjson{}
	req, err := enc.Encode(params)
	if err != nil {
		log.Printf("Cannot encode names for verification: %s.", err)
	}
	r := bytes.NewReader(req)
	resp, err := vr.client.Post(vr.verifierURL+"verifications", "application/json", r)
	if err != nil {
		log.Printf("Cannot send verification request: %s.", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Cannot read body of response: %s.", err)
	}
	var response []vlib.Verification
	err = enc.Decode(respBytes, &response)
	if err != nil {
		log.Printf("Cannot decode verification results: %s.", err)
	}
	return response
}
