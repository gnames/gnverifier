package verifrest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnquery/ent/search"
	"github.com/gnames/gnuuid"
	"github.com/gnames/gnverifier/pkg/ent/verifier"
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

func (vr verifrest) Search(
	ctx context.Context,
	inp search.Input,
) (search.Output, error) {
	enc := gnfmt.GNjson{}
	urlQ := vr.verifierURL + "search/" + url.PathEscape(inp.ToQuery())
	var res search.Output
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, urlQ, nil)
	if err != nil {
		slog.Error("Cannot create request", "error", err)
		return res, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := vr.client.Do(request)
	if err != nil {
		slog.Error("Cannot get data-sources information", "error", err)
		return res, err
	}
	defer resp.Body.Close()

	var respBytes []byte
	respBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Body reading is failing for a search", "error", err)
		return res, err
	}
	err = enc.Decode(respBytes, &res)
	if err != nil {
		slog.Error("Cannot decode search result", "error", err)
		return res, err
	}

	return res, nil
}

// DataSources returns meta-data about aggregated data-sources.
func (vr verifrest) DataSources(
	ctx context.Context,
) ([]vlib.DataSource, error) {
	enc := gnfmt.GNjson{}
	url := vr.verifierURL + "data_sources"
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Cannot create request", "error", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := vr.client.Do(request)
	if err != nil {
		slog.Error("Cannot get data-sources information", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Body reading is failing for data-sources", "error", err)
		return nil, err
	}
	response := make([]vlib.DataSource, 0)
	err = enc.Decode(respBytes, &response)
	if err != nil {
		slog.Error("Cannot decode data-sources", "error", err)
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
		slog.Error("Cannot create request", "error", err)
		return response, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := vr.client.Do(request)
	if err != nil {
		slog.Error("Cannot get data-sources information", "error", err)
		return response, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Body reading is failing for data-sources", "error", err)
		return response, err
	}
	err = enc.Decode(respBytes, &response)
	if err != nil {
		slog.Error("Cannot decode data-sources", "error", err)
		return response, err
	}
	return response, nil
}

func (vr verifrest) NameString(
	ctx context.Context,
	input vlib.NameStringInput,
) (vlib.NameStringOutput, error) {
	var res vlib.NameStringOutput
	enc := gnfmt.GNjson{}
	url := fmt.Sprintf("%sname_strings/%s", vr.verifierURL, input.ID)
	var params []string
	if len(input.DataSources) > 0 {
		nums := make([]string, len(input.DataSources))
		for i, v := range input.DataSources {
			nums[i] = strconv.Itoa(v)
		}
		ds := "data_sources=" + strings.Join(nums, ",")
		params = []string{ds}
	}
	if input.WithAllMatches {
		params = append(params, "all_matches=true")
	}
	if len(params) > 0 {
		url = url + "?" + strings.Join(params, "&")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Cannot create request", "error", err)
		return res, err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := vr.client.Do(request)
	if err != nil {
		slog.Error("Cannot get name-string information", "error", err)
		return res, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(
			"Body reading is failing for a name-string response",
			"error", err,
		)
		return res, err
	}
	err = enc.Decode(respBytes, &res)
	if err != nil {
		slog.Error("Cannot decode name-string response", "error", err)
		return res, err
	}
	return res, nil
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
		slog.Error("Cannot encode names for verification", "error", err)
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
			slog.Error("Cannot create request", "error", err)
		}
		request.Header.Set("Content-Type", "application/json")

		resp, err = vr.client.Do(request)
		if err != nil {
			slog.Error(
				"Request is failing for %s",
				"names-range", namesRange,
				"error", err,
			)
			return true, err
		}
		defer resp.Body.Close()

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			slog.Error(
				"Body reading is failing for",
				"names-range", namesRange,
				"error", err,
			)
			return true, err
		}
		err = enc.Decode(respBytes, &response)
		if err != nil {
			slog.Error(
				"Response decoding is failing for",
				"names-range", namesRange,
				"error", err,
			)
			return true, err
		}
		return false, nil
	})

	if err != nil {
		rng := fmt.Sprintf(
			"%s-%s",
			input.NameStrings[0],
			input.NameStrings[len(input.NameStrings)-1])
		slog.Warn(
			"Verification failed.",
			"names-range", rng,
			"attempts", attempts,
			"error", err,
		)
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
