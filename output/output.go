package output

import (
	"strconv"
	"strings"

	gne "github.com/gnames/gnames/domain/entity"
	gncsv "github.com/gnames/gnames/lib/csv"
	jsoniter "github.com/json-iterator/go"
)

func CSVHeader() string {
	return "Kind,MatchType,EditDistance,ScientificName,MatchedName,MatchedCanonical,TaxonId,CurrentName,Synonym,DataSourceId,DataSourceTitle,ClassificationPath"
}

func Output(ver *gne.Verification, f Format, pref_only bool) string {
	switch f {
	case CSV:
		return csvOutput(ver, pref_only)
	case CompactJSON:
		return jsonOutput(ver, pref_only, false)
	case PrettyJSON:
		return jsonOutput(ver, pref_only, true)
	}
	return "N/A"
}

func csvOutput(ver *gne.Verification, pref_only bool) string {
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

func csvRow(ver *gne.Verification, prefIndex int) string {
	kind := "BestMatch"
	res := ver.BestResult
	if res == nil {
		res = &gne.ResultData{}
	}

	if prefIndex > -1 {
		kind = "PreferredMatch"
		res = ver.PreferredResults[prefIndex]
	}

	s := []string{kind, ver.MatchType.String(), strconv.Itoa(res.EditDistance),
		ver.Input, res.MatchedName, res.MatchedCanonicalFull, res.ID,
		res.CurrentName, strconv.FormatBool(res.IsSynonym),
		strconv.Itoa(res.DataSourceID), res.DataSrouceTitleShort,
		res.ClassificationPath}
	return gncsv.ToCSV(s)
}

func jsonOutput(ver *gne.Verification, pref_only bool, pretty bool) string {
	res := []byte("bad JSON")
	if pref_only {
		ver.BestResult = nil
	}
	if pretty {
		res, _ = jsoniter.MarshalIndent(ver, "", "  ")
	} else {
		res, _ = jsoniter.Marshal(ver)
	}
	return string(res)
}
