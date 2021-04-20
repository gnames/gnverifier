package verifrest

import (
	"context"
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/stretchr/testify/assert"
)

var urlVerif = "https://verifier.globalnames.org/api/v1/"

func TestDataSources(t *testing.T) {
	r, err := recorder.New("fixtures/dss")
	defer r.Stop()
	assert.Nil(t, err)

	client := &http.Client{Transport: r}
	verif := &verifrest{
		verifierURL: urlVerif,
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
		verifierURL: urlVerif,
		client:      client,
	}

	ds, err := verif.DataSource(context.Background(), 4)
	assert.Nil(t, err)
	assert.Equal(t, ds.TitleShort, "NCBI")
}

func TestVerify(t *testing.T) {
	tests := []struct {
		fixture         string
		params          vlib.VerifyParams
		matchTypes      []vlib.MatchTypeValue
		matchCanonicals []string
		hasPreferred    []bool
	}{
		{
			fixture: "fixtures/capitalize",
			params: vlib.VerifyParams{
				NameStrings:        []string{"plantago major"},
				WithCapitalization: true,
			},
			matchTypes: []vlib.MatchTypeValue{
				vlib.Exact,
			},
			matchCanonicals: []string{"Plantago major"},
			hasPreferred:    []bool{false},
		},
		{
			fixture: "fixtures/name",
			params: vlib.VerifyParams{
				NameStrings: []string{"Plantago major L."},
			},
			matchTypes: []vlib.MatchTypeValue{
				vlib.Exact,
			},
			matchCanonicals: []string{"Plantago major"},
			hasPreferred:    []bool{false},
		},
		{
			fixture: "fixtures/name_pref",
			params: vlib.VerifyParams{
				NameStrings:      []string{"Pomatomus saltatrix (Linnaeus, 1766)"},
				PreferredSources: []int{1, 12},
			},
			matchTypes: []vlib.MatchTypeValue{
				vlib.Exact,
			},
			matchCanonicals: []string{"Pomatomus saltatrix"},
			hasPreferred:    []bool{true},
		},
		{
			fixture: "fixtures/names",
			params: vlib.VerifyParams{
				NameStrings: []string{
					"Pomatomus saltatrix (Linnaeus, 1766)",
					"Bubo bubo (Linnaeus, 1758)",
					"NotAName",
				},
				PreferredSources: []int{4, 12},
			},
			matchTypes: []vlib.MatchTypeValue{
				vlib.Exact,
				vlib.Exact,
				vlib.NoMatch,
			},
			matchCanonicals: []string{"Pomatomus saltatrix", "Bubo bubo", ""},
			hasPreferred:    []bool{true, true, false},
		},
	}
	for i := range tests {
		r, err := recorder.New(tests[i].fixture)
		assert.Nil(t, err)

		client := &http.Client{Transport: r}
		verif := &verifrest{
			verifierURL: urlVerif,
			client:      client,
		}
		vs := verif.Verify(context.Background(), tests[i].params)
		for j := range vs {
			assert.Equal(t, vs[j].MatchType, tests[i].matchTypes[j])
			if tests[i].matchCanonicals[j] != "" {
				assert.Equal(
					t,
					vs[j].BestResult.MatchedCanonicalSimple,
					tests[i].matchCanonicals[j],
				)
			}
			assert.Equal(t, len(vs[j].PreferredResults) > 0, tests[i].hasPreferred[j])
		}

		r.Stop()
	}
}
