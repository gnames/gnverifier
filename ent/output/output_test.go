package output_test

import (
	"os"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier/ent/output"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestOutput(t *testing.T) {
	verifs := verifications(t).Names
	tests := []struct {
		msg      string
		input    vlib.Name
		format   gnfmt.Format
		prefOnly bool
		test     func(*testing.T, string)
		linesNum int
	}{
		{
			msg:      "csv_prefOnly_false",
			input:    verifs[0],
			format:   gnfmt.CSV,
			prefOnly: false,
			test: func(t *testing.T, res string) {
				assert.NotContains(t, res, "inputID", "csv,false 1")
				assert.Contains(t, res, "BestMatch", "csv, false 2")
				assert.Contains(t, res, "PreferredMatch", "csv, false 3")
				assert.True(t, strings.HasPrefix(res, "BestMatch"), "csv, false 4")
			},
			linesNum: 3,
		},
		{
			msg:      "not_name_csv_prefOnly_true",
			input:    verifs[2],
			format:   gnfmt.CSV,
			prefOnly: true,
			test: func(t *testing.T, res string) {
				assert.NotContains(t, res, "inputID", "noname, csv, true 1")
				assert.NotContains(t, res, "BestMatch", "noname, csv, true 2")
				assert.Contains(t, res, "PreferredMatch", "noname, csv, true 3")
				assert.Contains(t, res, "NoMatch", "noname, csv, true 4")
			},
			linesNum: 1,
		},
		{
			msg:      "pretty",
			input:    verifs[0],
			format:   gnfmt.PrettyJSON,
			prefOnly: false,
			test: func(t *testing.T, res string) {
				assert.Contains(t, res, "id", "pretty 1")
				assert.Contains(t, res, "bestResult", "pretty 2")
				assert.Contains(t, res, "results", "pretty 3")
			},
			linesNum: 104,
		},
		{
			msg:      "compact",
			input:    verifs[0],
			format:   gnfmt.CompactJSON,
			prefOnly: false,
			test: func(t *testing.T, res string) {
				assert.Contains(t, res, "id", "compact 1")
				assert.Contains(t, res, "bestResult", "compact 2")
				assert.Contains(t, res, "results", "compact 3")
			},
			linesNum: 1,
		},
	}

	for i := range tests {
		t.Run(tests[i].msg, func(t *testing.T) {
			res := output.NameOutput(tests[i].input, tests[i].format, tests[i].prefOnly)
			lines := strings.Split(res, "\n")
			assert.Equal(t, len(lines), tests[i].linesNum)
			tests[i].test(t, res)
		})
	}

}

func verifications(t *testing.T) vlib.Output {
	c := cassette.New("dss")
	data, err := os.ReadFile("../../io/verifrest/fixtures/names.yaml")
	assert.Nil(t, err)
	err = yaml.Unmarshal(data, c)
	assert.Nil(t, err)
	dssStr := c.Interactions[0].Response.Body
	enc := gnfmt.GNjson{}
	var res vlib.Output
	err = enc.Decode([]byte(dssStr), &res)
	assert.Nil(t, err)
	return res
}
