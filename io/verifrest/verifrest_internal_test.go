package verifrest

import (
	"context"
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnquery"
	"github.com/stretchr/testify/assert"
)

var urlAPI = "https://verifier.globalnames.org/api/v0/"

func TestDataSources(t *testing.T) {
	r, err := recorder.New("fixtures/dss")
	defer r.Stop()
	assert.Nil(t, err)

	client := &http.Client{Transport: r}
	verif := &verifrest{
		verifierURL: urlAPI,
		client:      client,
	}
	ds, err := verif.DataSources(context.Background())
	assert.Nil(t, err)
	assert.Greater(t, len(ds), 50)
}

func TestDataSource(t *testing.T) {
	r, err := recorder.New("fixtures/ds4")
	defer r.Stop()
	assert.Nil(t, err)

	client := &http.Client{Transport: r}
	verif := &verifrest{
		verifierURL: urlAPI,
		client:      client,
	}

	ds, err := verif.DataSource(context.Background(), 4)
	assert.Nil(t, err)
	assert.Equal(t, "NCBI", ds.TitleShort)
}

func TestSearch(t *testing.T) {
	tests := []struct {
		fixture    string
		query      string
		hasResults bool
	}{
		{"fixtures/search", "n:Bubo bubo all:t", true},
	}
	for i, v := range tests {
		r, err := recorder.New(tests[i].fixture)
		assert.Nil(t, err)

		client := &http.Client{Transport: r}
		verif := &verifrest{
			verifierURL: urlAPI,
			client:      client,
		}

		gnq := gnquery.New()
		inp := gnq.Parse(v.query)
		vs, err := verif.Search(context.Background(), inp)
		assert.Nil(t, err)
		for j := range vs.Names {
			assert.Equal(t, v.hasResults, len(vs.Names[j].Results) > 0)
		}

		r.Stop()
	}
}

func TestVerify(t *testing.T) {
	tests := []struct {
		fixture         string
		params          vlib.Input
		matchTypes      []vlib.MatchTypeValue
		matchCanonicals []string
		hasResults      []bool
	}{
		{
			fixture: "fixtures/capitalize",
			params: vlib.Input{
				NameStrings:        []string{"plantago major"},
				WithCapitalization: true,
			},
			matchTypes: []vlib.MatchTypeValue{
				vlib.Exact,
			},
			matchCanonicals: []string{"Plantago major"},
			hasResults:      []bool{false},
		},
		{
			fixture: "fixtures/name",
			params: vlib.Input{
				NameStrings: []string{"Plantago major L."},
			},
			matchTypes: []vlib.MatchTypeValue{
				vlib.Exact,
			},
			matchCanonicals: []string{"Plantago major"},
			hasResults:      []bool{false},
		},
		{
			fixture: "fixtures/name_pref",
			params: vlib.Input{
				NameStrings: []string{"Pomatomus saltatrix (Linnaeus, 1766)"},
				DataSources: []int{1, 12},
			},
			matchTypes: []vlib.MatchTypeValue{
				vlib.Exact,
			},
			matchCanonicals: []string{"Pomatomus saltatrix"},
			hasResults:      []bool{true},
		},
		{
			fixture: "fixtures/names",
			params: vlib.Input{
				NameStrings: []string{
					"Pomatomus saltatrix (Linnaeus, 1766)",
					"Bubo bubo (Linnaeus, 1758)",
					"NotAName",
				},
				DataSources: []int{4, 12},
			},
			matchTypes: []vlib.MatchTypeValue{
				vlib.Exact,
				vlib.Exact,
				vlib.NoMatch,
			},
			matchCanonicals: []string{"Pomatomus saltatrix", "Bubo bubo", ""},
			hasResults:      []bool{true, true, false},
		},
	}
	for i := range tests {
		r, err := recorder.New(tests[i].fixture)
		assert.Nil(t, err)

		client := &http.Client{Transport: r}
		verif := &verifrest{
			verifierURL: urlAPI,
			client:      client,
		}
		vs := verif.Verify(context.Background(), tests[i].params)
		for j := range vs.Names {
			assert.Equal(t, tests[i].matchTypes[j], vs.Names[j].MatchType)
			if tests[i].matchCanonicals[j] != "" {
				assert.Equal(
					t,
					tests[i].matchCanonicals[j],
					vs.Names[j].BestResult.MatchedCanonicalSimple,
				)
			}
			assert.Equal(t, tests[i].hasResults[j], len(vs.Names[j].Results) > 0)
		}

		r.Stop()
	}
}
