package output

import (
	"strconv"
	"strings"

	gncsv "github.com/gnames/gnlib/csv"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/encode"
	"github.com/gnames/gnlib/format"
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
)

func Output(ver vlib.Verification, f format.Format, pref_only bool) string {
	switch f {
	case format.CSV:
		return csvOutput(ver, pref_only)
	case format.CompactJSON:
		return jsonOutput(ver, pref_only, false)
	case format.PrettyJSON:
		return jsonOutput(ver, pref_only, true)
	}
	return "N/A"
}

func CSVHeader() string {
	return "Kind,MatchType,EditDistance,ScientificName,MatchedName,MatchedCanonical,TaxonId,CurrentName,Synonym,DataSourceId,DataSourceTitle,ClassificationPath"
}

func csvOutput(ver vlib.Verification, pref_only bool) string {
	var res []string
	if !pref_only {
		best := csvRow(ver, -1)
		res = append(res, best)
	}
	for i := range ver.PreferredResults {
		pref := csvRow(ver, i)
		res = append(res, pref)
	}

	return strings.Join(res, "\n")
}

func csvRow(ver vlib.Verification, prefIndex int) string {
	kind := "BestMatch"
	res := ver.BestResult

	if prefIndex > -1 {
		kind = "PreferredMatch"
		res = ver.PreferredResults[prefIndex]
	}

	s := []string{
		kind, ver.MatchType.String(), "", ver.Input,
		"", "", "", "", "", "", "", "",
	}

	if res != nil {
		s[editDistance] = strconv.Itoa(res.EditDistance)
		s[matchedName] = res.MatchedName
		s[matchedCanonical] = res.MatchedCanonicalFull
		s[taxonID] = res.ID
		s[currentName] = res.CurrentName
		s[synonym] = strconv.FormatBool(res.IsSynonym)
		s[dataSourceID] = strconv.Itoa(res.DataSourceID)
		s[dataSourceTitle] = res.DataSourceTitleShort
		s[classificationPath] = res.ClassificationPath
	}

	return gncsv.ToCSV(s)
}

func jsonOutput(ver vlib.Verification, pref_only bool, pretty bool) string {
	enc := encode.GNjson{Pretty: pretty}
	if pref_only {
		ver.BestResult = nil
	}
	res, _ := enc.Encode(ver)
	return string(res)
}
