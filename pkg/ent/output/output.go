package output

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
)

type csvField int

const (
	kind csvField = iota
	sortScore
	matchType
	editDistance
	input
	matchedName
	matchedCanonical
	taxonID
	currentName
	taxonomicStatus
	dataSourceID
	dataSourceTitle
	classificationPath
	error
)

const sortedMatch = "SortedMatch"

// NameOutput takes result of verification for one string and converts it into
// required format (CSV or JSON).
func NameOutput(ver vlib.Name, f gnfmt.Format) string {
	switch f {
	case gnfmt.CSV:
		return csvOutput(ver, ',')
	case gnfmt.TSV:
		return csvOutput(ver, '\t')
	case gnfmt.CompactJSON:
		return jsonOutput(ver, false)
	case gnfmt.PrettyJSON:
		return jsonOutput(ver, true)
	}
	return "N/A"
}

// CSVHeader returns the header string for CSV output format.
func CSVHeader(f gnfmt.Format) string {
	header := []string{"Kind", "SortScore", "MatchType", "EditDistance", "ScientificName",
		"MatchedName", "MatchedCanonical", "TaxonId", "CurrentName", "TaxonomicStatus",
		"DataSourceId", "DataSourceTitle", "ClassificationPath", "Error"}
	switch f {
	case gnfmt.CSV:
		return gnfmt.ToCSV(header, ',')
	case gnfmt.TSV:
		return gnfmt.ToCSV(header, '\t')
	default:
		return ""
	}
}

func csvOutput(ver vlib.Name, sep rune) string {
	var res []string
	if ver.BestResult != nil {
		best := csvRow(ver, -1, sep)
		res = append(res, best)
	} else if len(ver.Results) == 0 {
		res = append(res, csvEmptyRow(ver, sep))
	}
	for i := range ver.Results {
		pref := csvRow(ver, i, sep)
		res = append(res, pref)
	}

	return strings.Join(res, "\n")
}

func csvEmptyRow(ver vlib.Name, sep rune) string {
	s := []string{
		sortedMatch, "0.0", vlib.NoMatch.String(), "", ver.Name,
		"", "", "", "", "", "", "", "", ver.Error,
	}
	return gnfmt.ToCSV(s, sep)
}

func csvRow(ver vlib.Name, prefIndex int, sep rune) string {
	kind := "BestMatch"
	res := ver.BestResult

	if prefIndex > -1 {
		if prefIndex > 0 {
			kind = sortedMatch
		}
		res = ver.Results[prefIndex]
	}

	s := []string{
		kind, "0.0", vlib.NoMatch.String(), "", ver.Name,
		"", "", "", "", "", "", "", "", ver.Error,
	}

	if res != nil {
		s[editDistance] = strconv.Itoa(res.EditDistance)
		s[sortScore] = fmt.Sprintf("%0.5f", res.SortScore)
		s[matchType] = res.MatchType.String()
		s[matchedName] = res.MatchedName
		s[matchedCanonical] = res.MatchedCanonicalFull
		s[taxonID] = res.RecordID
		s[currentName] = res.CurrentName
		s[taxonomicStatus] = res.TaxonomicStatus
		s[dataSourceID] = strconv.Itoa(res.DataSourceID)
		s[dataSourceTitle] = res.DataSourceTitleShort
		s[classificationPath] = res.ClassificationPath
	}

	return gnfmt.ToCSV(s, sep)
}

func jsonOutput(ver vlib.Name, pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(ver)
	return string(res)
}
