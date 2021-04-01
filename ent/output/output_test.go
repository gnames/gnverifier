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
	verifs := verifications(t)
	tests := []struct {
		msg      string
		input    vlib.Verification
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
				assert.Contains(t, res, "inputId", "pretty 1")
				assert.Contains(t, res, "bestResult", "pretty 2")
				assert.Contains(t, res, "preferredResults", "pretty 3")
			},
			linesNum: 78,
		},
		{
			msg:      "compact",
			input:    verifs[0],
			format:   gnfmt.CompactJSON,
			prefOnly: false,
			test: func(t *testing.T, res string) {
				assert.Contains(t, res, "inputId", "compact 1")
				assert.Contains(t, res, "bestResult", "compact 2")
				assert.Contains(t, res, "preferredResults", "compact 3")
			},
			linesNum: 1,
		},
	}

	for i := range tests {
		t.Run(tests[i].msg, func(t *testing.T) {
			res := output.Output(tests[i].input, tests[i].format, tests[i].prefOnly)
			lines := strings.Split(res, "\n")
			assert.Equal(t, len(lines), tests[i].linesNum)
			tests[i].test(t, res)
		})
	}

}

func verifications(t *testing.T) []vlib.Verification {
	c := cassette.New("dss")
	data, err := os.ReadFile("../../io/verifrest/fixtures/names.yaml")
	assert.Nil(t, err)
	err = yaml.Unmarshal(data, c)
	assert.Nil(t, err)
	dssStr := c.Interactions[0].Response.Body
	enc := gnfmt.GNjson{}
	res := make([]vlib.Verification, 0)
	err = enc.Decode([]byte(dssStr), &res)
	assert.Nil(t, err)
	return res
}
