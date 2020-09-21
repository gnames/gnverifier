package entity

import (
	gne "github.com/gnames/gnames/domain/entity"
)

type Output struct {
	/// Name-string supplied by user for verification.
	Name string
	/// Match type of the best result after verification attempt.
	gne.MatchType
	/// The number of Data Sources that could be matched to the name-string.
	DataSourcesNum int
	/// Indicates if the name was matched to Data Sources with human or
	/// automatic curation of the data.
	gne.CurationLevel
	/// How many retries were needed to send the name-string to gnindex
	/// server.
	Retries int
	/// Contains an error string (if any) after verification attempt.
	Error string
	/// The apparent best match of the name-string to gnindex data sets.
	/// The best match is determined by a score that takes in account if
	/// the match was exact, partial, or fuzzy, if it was a match of uninomial,
	/// binomial, or multinomial, if there authors matched in the name-string
	/// and gnindex data.
	BestResult *gne.ResultData
	/// Contains all matches found in the user-specified Data Sources.
	PreferredResults []*gne.ResultData
}

