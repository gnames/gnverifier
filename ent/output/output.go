package output

import (
	"strconv"
	"strings"

	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
)

type csvField int

const (
	kind csvField = iota
	matchType
	editDistance
	input
	matchedName
	matchedCanonical
	taxonID
	currentName
	synonym
	dataSourceID
	dataSourceTitle
	classificationPath
	error
)

const prefMatch = "PreferredMatch"

// Output takes result of verification for one string and converts it into
// required format (CSV or JSON).
func Output(ver vlib.Verification, f gnfmt.Format, prefOnly bool) string {
	switch f {
	case gnfmt.CSV:
		return csvOutput(ver, prefOnly)
	case gnfmt.CompactJSON:
		return jsonOutput(ver, prefOnly, false)
	case gnfmt.PrettyJSON:
		return jsonOutput(ver, prefOnly, true)
	}
	return "N/A"
}

// CSVHeader returns the header string for CSV output format.
func CSVHeader() string {
	return "Kind,MatchType,EditDistance,ScientificName,MatchedName,MatchedCanonical,TaxonId,CurrentName,Synonym,DataSourceId,DataSourceTitle,ClassificationPath,Error"
}

func csvOutput(ver vlib.Verification, prefOnly bool) string {
	var res []string
	if !prefOnly {
		best := csvRow(ver, -1)
		res = append(res, best)
	}
	if prefOnly && len(ver.PreferredResults) == 0 {
		res = append(res, csvNoPrefRow(ver))
	}
	for i := range ver.PreferredResults {
		pref := csvRow(ver, i)
		res = append(res, pref)
	}

	return strings.Join(res, "\n")
}

func csvNoPrefRow(ver vlib.Verification) string {
	s := []string{
		prefMatch, vlib.NoMatch.String(), "", ver.Input,
		"", "", "", "", "", "", "", "", ver.Error,
	}
	return gnfmt.ToCSV(s)
}

func csvRow(ver vlib.Verification, prefIndex int) string {
	kind := "BestMatch"
	res := ver.BestResult

	if prefIndex > -1 {
		kind = prefMatch
		res = ver.PreferredResults[prefIndex]
	}

	s := []string{
		kind, ver.MatchType.String(), "", ver.Input,
		"", "", "", "", "", "", "", "", ver.Error,
	}

	if res != nil {
		s[editDistance] = strconv.Itoa(res.EditDistance)
		s[matchedName] = res.MatchedName
		s[matchedCanonical] = res.MatchedCanonicalFull
		s[taxonID] = res.RecordID
		s[currentName] = res.CurrentName
		s[synonym] = strconv.FormatBool(res.IsSynonym)
		s[dataSourceID] = strconv.Itoa(res.DataSourceID)
		s[dataSourceTitle] = res.DataSourceTitleShort
		s[classificationPath] = res.ClassificationPath
	}

	return gnfmt.ToCSV(s)
}

func jsonOutput(ver vlib.Verification, prefOnly bool, pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	if prefOnly {
		ver.BestResult = nil
	}
	res, _ := enc.Encode(ver)
	return string(res)
}
