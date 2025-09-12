package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/gnames/gnfmt"
	gnverifier "github.com/gnames/gnverifier/pkg"
	"github.com/gnames/gnverifier/pkg/config"
	"github.com/spf13/cobra"
)

type funcFlag func(cmd *cobra.Command)

func versionFlag(cmd *cobra.Command) {
	hasVersionFlag, _ := cmd.Flags().GetBool("version")
	if hasVersionFlag {
		fmt.Printf("\nversion: %s\nbuild: %s\n\n", gnverifier.Version, gnverifier.Build)
		os.Exit(0)
	}
}

func quietFlag(cmd *cobra.Command) {
	quiet, _ := cmd.Flags().GetBool("quiet")
	if quiet {
		slog.SetLogLoggerLevel(10)
	}
}

func capitalizeFlag(cmd *cobra.Command) {
	caps, _ := cmd.Flags().GetBool("capitalize")
	if caps {
		opts = append(opts, config.OptWithCapitalization(true))
	}
}

func spGroupFlag(cmd *cobra.Command) {
	spGr, _ := cmd.Flags().GetBool("species_group")
	if spGr {
		opts = append(opts, config.OptWithSpeciesGroup(true))
	}
}

func fuzzyRelaxedFlag(cmd *cobra.Command) {
	fuzzyRelaxed, _ := cmd.Flags().GetBool("fuzzy_relaxed")
	if fuzzyRelaxed {
		opts = append(opts, config.OptWithRelaxedFuzzyMatch(true))
	}
}

func fuzzyUninomialFlag(cmd *cobra.Command) {
	fuzzyUni, _ := cmd.Flags().GetBool("fuzzy_uninomial")
	if fuzzyUni {
		opts = append(opts, config.OptWithUninomialFuzzyMatch(true))
	}
}

func formatFlag(cmd *cobra.Command) {
	formatString, _ := cmd.Flags().GetString("format")
	frmt, _ := gnfmt.NewFormat(formatString)
	if frmt == gnfmt.FormatNone {
		slog.Warn("Cannot set format with user inoput, setting format to csv",
			"input", formatString,
		)
		frmt = gnfmt.CSV
	}
	opts = append(opts, config.OptFormat(frmt))
}

func jobsFlag(cmd *cobra.Command) {
	jobs, _ := cmd.Flags().GetInt("jobs")
	if jobs != 4 && jobs > 0 {
		opts = append(opts, config.OptJobs(jobs))
	}
}

func allMatchesFlag(cmd *cobra.Command) {
	allMatches, _ := cmd.Flags().GetBool("all_matches")
	if allMatches {
		opts = append(opts, config.OptWithAllMatches(true))
	}
}

func sourcesFlag(cmd *cobra.Command) {
	sources, _ := cmd.Flags().GetString("sources")
	if sources != "" {
		data_sources := parseDataSources(sources)
		opts = append(opts, config.OptDataSources(data_sources))
	}
}

func verifierUrlFlag(cmd *cobra.Command) {
	url, _ := cmd.Flags().GetString("verifier_url")
	if url != "" {
		opts = append(opts, config.OptVerifierURL(url))
		webOpts = append(webOpts, config.OptVerifierURL(url))
	}
}

func vernacularsFlag(cmd *cobra.Command) {
	vernLangs, _ := cmd.Flags().GetString("vernaculars")
	if vernLangs != "" {
		langs := parseVernacularLanguages(vernLangs)
		opts = append(opts, config.OptVernaculars(langs))
	}
}

func parseDataSources(s string) []int {
	if s == "" {
		return nil
	}
	dss := strings.Split(s, ",")
	res := make([]int, 0, len(dss))
	for _, v := range dss {
		v = strings.Trim(v, " ")
		ds, err := strconv.Atoi(v)
		if err != nil {
			slog.Warn("Cannot convert data-sources to list, skipping", "input", v)
			return nil
		}
		if ds < 0 {
			slog.Warn("Data source ID is less than zero, skipping", "input", ds)
		} else {
			res = append(res, int(ds))
		}
	}
	if len(res) > 0 {
		return res
	}
	return nil
}

func parseVernacularLanguages(langs string) []string {
	var res []string
	for v := range strings.SplitSeq(langs, ",") {
		v := strings.TrimSpace(v)
		if len(v) != 3 {
			slog.Warn("languages for vernacular names has to be 3 letters long", "input", v)
			continue
		}
		res = append(res, v)
	}
	return res
}
