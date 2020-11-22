package output

import (
	"strconv"
	"strings"

	gncsv "github.com/gnames/gnlib/csv"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/format"
	jsoniter "github.com/json-iterator/go"
)

func CSVHeader() string {
	return "Kind,MatchType,EditDistance,ScientificName,MatchedName,MatchedCanonical,TaxonId,CurrentName,Synonym,DataSourceId,DataSourceTitle,ClassificationPath"
}

func Output(ver *vlib.Verification, f format.Format, pref_only bool) string {
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

func csvOutput(ver *vlib.Verification, pref_only bool) string {
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

func csvRow(ver *vlib.Verification, prefIndex int) string {
	kind := "BestMatch"
	res := ver.BestResult
	if res == nil {
		res = &vlib.ResultData{}
	}

	if prefIndex > -1 {
		kind = "PreferredMatch"
		res = ver.PreferredResults[prefIndex]
	}

	dsID := ""
	if res.DataSourceID != nil {
		dsID = strconv.Itoa(*res.DataSourceID)
	}

	ed := ""
	if res.EditDistance != nil {
		ed = strconv.Itoa(*res.EditDistance)
	}

	s := []string{kind, ver.MatchType.String(), ed,
		ver.Input, res.MatchedName, res.MatchedCanonicalFull, res.ID,
		res.CurrentName, strconv.FormatBool(res.IsSynonym),
		dsID, res.DataSrouceTitleShort,
		res.ClassificationPath}
	return gncsv.ToCSV(s)
}

func jsonOutput(ver *vlib.Verification, pref_only bool, pretty bool) string {
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
