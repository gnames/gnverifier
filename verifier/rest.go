package verifier

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gnames/gnames/domain/entity"
	"github.com/gnames/gnames/lib/encode"
)

type VerifierRest struct {
	VerifierURL string
	encode.Encoder
}

func NewVerifierRest(url string) VerifierRest {
	if url[len(url)-1] != '/' {
		url = url + "/"
	}
	return VerifierRest{VerifierURL: url, Encoder: encode.GNjson{}}
}

// Verify takes names-strings and options and returns verification result.
func (vr VerifierRest) Verify(params entity.VerifyParams) []*entity.Verification {
	req, err := encode.GNjson{}.Encode(params)
	if err != nil {
		log.Printf("Cannot encode names for verification: %s.", err)
	}
	r := bytes.NewReader(req)
	resp, err := http.Post(vr.VerifierURL+"verification", "application/x-binary", r)
	if err != nil {
		log.Printf("Cannot send verification request: %s.", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Cannot read body of response: %s.", err)
	}
	var response []*entity.Verification
	err = encode.GNjson{}.Decode(respBytes, &response)
	if err != nil {
		log.Printf("Cannot decode verification results: %s.", err)
	}
	return response
}

// DataSources takes data-source id and opts and returns the data-source
// metadata.  If no id is provided, it returns metadata for all data-sources.
func (vr VerifierRest) DataSources(entity.DataSourcesOpts) []*entity.DataSource {
	var res []*entity.DataSource
	return res
}
