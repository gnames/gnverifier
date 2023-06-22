// Package cmd creates a command line interface for gnverifier app.
package cmd

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnquery"
	"github.com/gnames/gnquery/ent/search"
	"github.com/gnames/gnsys"
	"github.com/gnames/gnverifier/internal/ent/output"
	"github.com/gnames/gnverifier/internal/io/verifrest"
	"github.com/gnames/gnverifier/internal/io/web"
	gnverifier "github.com/gnames/gnverifier/pkg"
	"github.com/gnames/gnverifier/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed gnverifier.yaml
var configText string

var (
	opts []config.Option
	msg  []string
)

// cfgData purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type cfgData struct {
	DataSources             []int
	Format                  string
	Jobs                    int
	NsqdContainsFilter      string
	NsqdRegexFilter         string
	NsqdTCPAddress          string
	VerifierURL             string
	WithAllMatches          bool
	WithCapitalization      bool
	WithSpeciesGroup        bool
	WithUninomialFuzzyMatch bool
	WithWebLogs             bool
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnverifier",
	Short: "Verifies scientific names agains many sources.",
	Long: `gnverifier uses a remote service to verify scientific names against
more than 100 biodiverisity data-sources. See more info at
https://github.com/gnames/gnverifier

  examples:
    gnverifier "Pardosa moesta"
    gnverifier file_with_names.txt
    gnverifier "g:M. sp:galloprovincialis au:Oliv."
`,
	Run: func(cmd *cobra.Command, args []string) {
		webOpts := make([]config.Option, len(opts))
		copy(webOpts, opts)
		webOpts = append(webOpts, config.OptWithCapitalization(true))

		if showVersionFlag(cmd) {
			os.Exit(0)
		}

		quiet, _ := cmd.Flags().GetBool("quiet")
		if quiet {
			zerolog.SetGlobalLevel(zerolog.Disabled)
		}

		for i := range msg {
			log.Info().Msg(msg[i])
		}

		caps, _ := cmd.Flags().GetBool("capitalize")
		if caps {
			opts = append(opts, config.OptWithCapitalization(true))
		}

		spGr, _ := cmd.Flags().GetBool("species_group")
		if spGr {
			opts = append(opts, config.OptWithSpeciesGroup(true))
		}

		fuzzyUni, _ := cmd.Flags().GetBool("fuzzy_uninomial")
		if fuzzyUni {
			opts = append(opts, config.OptWithUninomialFuzzyMatch(true))
		}

		formatString, _ := cmd.Flags().GetString("format")
		frmt, _ := gnfmt.NewFormat(formatString)
		if frmt == gnfmt.FormatNone {
			log.Warn().
				Msgf(
					"Cannot set format from '%s', setting format to csv",
					formatString,
				)
			frmt = gnfmt.CSV
		}
		opts = append(opts, config.OptFormat(frmt))

		jobs, _ := cmd.Flags().GetInt("jobs")
		if jobs != 4 {
			opts = append(opts, config.OptJobs(jobs))
		}

		allMatches, _ := cmd.Flags().GetBool("all_matches")
		if allMatches {
			opts = append(opts, config.OptWithAllMatches(true))
		}

		sources, _ := cmd.Flags().GetString("sources")
		if sources != "" {
			data_sources := parseDataSources(sources)
			opts = append(opts, config.OptDataSources(data_sources))
		}

		url, _ := cmd.Flags().GetString("verifier_url")
		if url != "" {
			opts = append(opts, config.OptVerifierURL(url))
			webOpts = append(webOpts, config.OptVerifierURL(url))
		}

		port, _ := cmd.Flags().GetInt("port")
		if port > 0 {
			weblogs, _ := cmd.Flags().GetBool("web-logs")
			if weblogs {
				webOpts = append(webOpts, config.OptWithWebLogs(true))
			}
			nsqAddr, _ := cmd.Flags().GetString("nsqd-tcp")
			if nsqAddr != "" {
				webOpts = append(webOpts, config.OptNsqdTCPAddress(nsqAddr))
			}

			log.Logger = zerolog.New(os.Stderr).With().
				Str("gnApp", "gnverifier").Logger()
			cnf := config.New(webOpts...)
			vfr := verifrest.New(cnf.VerifierURL)
			gnv := gnverifier.New(cnf, vfr)
			web.Run(gnv, port)
			os.Exit(0)
		}

		cfg := config.New(opts...)
		vfr := verifrest.New(cfg.VerifierURL)

		if len(args) == 0 {
			processStdin(cmd, cfg)
			os.Exit(0)
		}
		data := getInput(cmd, args)
		gnv := gnverifier.New(cfg, vfr)
		verify(gnv, data)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "V", false, "Prints version information")
	rootCmd.Flags().BoolP("capitalize", "c", false, "capitalizes first character")
	rootCmd.Flags().StringP("format", "f", "csv", `Format of the output: "compact", "pretty", "csv", "tsv".
  compact: compact JSON,
  pretty: pretty JSON,
  csv: CSV (DEFAULT)`)
	rootCmd.Flags().IntP("name_field", "n", 1, "Set position of ScientificName field, the first field is 1.")
	rootCmd.Flags().IntP("jobs", "j", 4, "Number of jobs running in parallel.")
	rootCmd.Flags().IntP("port", "p", 0, "Port to run web GUI.")
	rootCmd.Flags().BoolP("all_matches", "M", false, "return all matched results per source, not just the best one.")
	rootCmd.Flags().BoolP("species_group", "g", false, "searching for species names also searches their species groups.")
	rootCmd.Flags().BoolP("fuzzy_uninomial", "z", false,
		"allows fuzzy matching for uninomial names.")
	rootCmd.Flags().BoolP("quiet", "q", false, "do not show progress")
	rootCmd.Flags().StringP("sources", "s", "", `IDs of important data-sources to verify against (ex "1,11").
  If sources are set and there are matches to their data,
  such matches are returned in "preferred_result" results.
  If the option is set to "0" all matched sources are returned.

  To find IDs refer to "https://verifier.globalnames.org/data_sources".
  1 - Catalogue of Life
  3 - ITIS
  4 - NCBI
  9 - WoRMS
  11 - GBIF
  12 - Encyclopedia of Life
  167 - IPNI
  170 - Arctos
  172 - PaleoBioDB
  181 - IRMNG
  194 - PLAZI
  195 - AlgaeBase`)
	rootCmd.Flags().StringP("verifier_url", "v", "",
		`URL for verification service.
  Default: https://verifier.globalnames.org/api/v1`)
	rootCmd.Flags().BoolP("web-logs", "", false, "enable logs for the web service")
	rootCmd.Flags().StringP("nsqd-tcp", "", "", "an addresss pointing to NSQ TCP service for logs redirection (e.g. 127.0.0.1:4150)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var configDir string
	var err error
	configFile := "gnverifier"

	// Find config directory.
	configDir, err = os.UserConfigDir()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot find config directory")
	}

	// Search config in home directory with name ".gnmatcher" (without extension).
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFile)

	// Set environment variables to override
	// config file settings
	_ = viper.BindEnv("DataSources", "GNV_DATA_SOURCES")
	_ = viper.BindEnv("Format", "GNV_FORMAT")
	_ = viper.BindEnv("Jobs", "GNV_JOBS")
	_ = viper.BindEnv("NsqdContainsFilter", "GNV_NSQD_CONTAINS_FILTER")
	_ = viper.BindEnv("NsqdRegexFilter", "GNV_NSQD_REGEX_FILTER")
	_ = viper.BindEnv("NsqdTCPAddress", "GNV_NSQD_TCP_ADDRESS")
	_ = viper.BindEnv("VerifierURL", "GNV_VERIFIER_URL")
	_ = viper.BindEnv("WithAllMatches", "GNV_WITH_ALL_MATCHES")
	_ = viper.BindEnv("WithCapitalization", "GNV_WITH_CAPITALIZATION")
	_ = viper.BindEnv("WithSpeciesGroup", "GNV_WITH_SPECIES_GROUP")
	_ = viper.BindEnv(
		"WithUninomialFuzzyMatch",
		"GNV_WITH_UNINOMIAL_FUZZY_MATCH",
	)
	_ = viper.BindEnv("WithWebLogs", "GNV_WITH_WEB_LOGS")

	viper.AutomaticEnv() // read in environment variables that match

	configPath := filepath.Join(configDir, fmt.Sprintf("%s.yaml", configFile))
	touchConfigFile(configPath)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		msg = append(msg,
			fmt.Sprintf("Using config file: %s.", viper.ConfigFileUsed()))
	}
	getOpts()
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts() {
	cfg := &cfgData{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot deserialize config data")
	}

	if len(cfg.DataSources) > 0 {
		opts = append(opts, config.OptDataSources(cfg.DataSources))
	}
	if cfg.Format != "" {
		cfgFormat, err := gnfmt.NewFormat(cfg.Format)
		if err != nil {
			cfgFormat = gnfmt.CSV
		}
		opts = append(opts, config.OptFormat(cfgFormat))
	}
	if cfg.Jobs > 0 {
		opts = append(opts, config.OptJobs(cfg.Jobs))
	}
	if cfg.NsqdContainsFilter != "" {
		opts = append(opts, config.OptNsqdContainsFilter(cfg.NsqdContainsFilter))
	}
	if cfg.NsqdRegexFilter != "" {
		opts = append(opts, config.OptNsqdRegexFilter(cfg.NsqdRegexFilter))
	}
	if cfg.NsqdTCPAddress != "" {
		opts = append(opts, config.OptNsqdTCPAddress(cfg.NsqdTCPAddress))
	}
	if cfg.VerifierURL != "" {
		opts = append(opts, config.OptVerifierURL(cfg.VerifierURL))
	}
	if cfg.WithAllMatches {
		opts = append(opts, config.OptWithAllMatches(true))
	}
	if cfg.WithCapitalization {
		opts = append(opts, config.OptWithCapitalization(true))
	}
	if cfg.WithSpeciesGroup {
		opts = append(opts, config.OptWithSpeciesGroup(true))
	}
	if cfg.WithUninomialFuzzyMatch {
		opts = append(opts, config.OptWithUninomialFuzzyMatch(true))
	}
	if cfg.WithWebLogs {
		opts = append(opts, config.OptWithWebLogs(true))
	}
}

// showVersionFlag provides version and the build timestamp. If it returns
// true, it means that version flag was given.
func showVersionFlag(cmd *cobra.Command) bool {
	hasVersionFlag, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot get version flag")
	}

	if hasVersionFlag {
		fmt.Printf("\nversion: %s\nbuild: %s\n\n", gnverifier.Version, gnverifier.Build)
	}
	return hasVersionFlag
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
			log.Warn().Msgf("Cannot convert data-source '%s' to list, skipping", v)
			return nil
		}
		if ds < 0 {
			log.Warn().Msgf("Data source ID %d is less than zero, skipping", ds)
		} else {
			res = append(res, int(ds))
		}
	}
	if len(res) > 0 {
		return res
	}
	return nil
}

func processStdin(
	cmd *cobra.Command,
	cfg config.Config,
) {
	if !checkStdin() {
		_ = cmd.Help()
		return
	}
	vfr := verifrest.New(cfg.VerifierURL)
	gnv := gnverifier.New(cfg, vfr)
	verifyFile(gnv, os.Stdin)
}

func checkStdin() bool {
	stdInFile := os.Stdin
	stat, err := stdInFile.Stat()
	if err != nil {
		log.Fatal().Err(err).Msg("checkStdin")
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func getInput(cmd *cobra.Command, args []string) string {
	var data string
	switch len(args) {
	case 1:
		data = args[0]
	default:
		_ = cmd.Help()
		os.Exit(0)
	}
	return data
}

func verify(gnv gnverifier.GNverifier, str string) {
	fileExists, _ := gnsys.FileExists(str)
	if fileExists {
		f, err := os.OpenFile(str, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal().Err(err).Msg("verify")
		}
		verifyFile(gnv, f)
		f.Close()
	} else if search.IsQuery(str) {
		searchQuery(gnv, str)
	} else {
		verifyString(gnv, str)
	}
}

func verifyFile(gnv gnverifier.GNverifier, f io.Reader) {
	batch := gnv.Config().Batch
	in := make(chan []string)
	out := make(chan []vlib.Name)
	var wg sync.WaitGroup
	wg.Add(1)
	go gnv.VerifyStream(context.Background(), in, out)
	go processResults(gnv, out, &wg)
	sc := bufio.NewScanner(f)
	names := make([]string, 0, batch)
	for sc.Scan() {
		name := sc.Text()
		names = append(names, strings.Trim(name, " "))
		if len(names) == batch {
			in <- names
			names = make([]string, 0, batch)
		}
	}
	in <- names
	close(in)
	wg.Wait()
}

func processResults(gnv gnverifier.GNverifier, out <-chan []vlib.Name,
	wg *sync.WaitGroup) {
	defer wg.Done()
	timeStart := time.Now().UnixNano()
	f := gnv.Config().Format
	if f == gnfmt.CSV || f == gnfmt.TSV {
		fmt.Println(output.CSVHeader(f))
	}
	var count int
	for o := range out {
		count++
		total := int64(count * len(o))
		timeSpent := float64(time.Now().UnixNano()-timeStart) / 1_000_000_000
		speed := int64(float64(total) / timeSpent)

		log.Info().
			Str("names/sec", humanize.Comma(speed)).
			Str("names", humanize.Comma(int64(total))).
			Msg("Verified")
		for _, r := range o {
			if r.Error != "" {
				log.Warn().Msg(r.Error)
			}
			fmt.Println(output.NameOutput(r, f))
		}
	}
}

func verifyString(gnv gnverifier.GNverifier, name string) {
	res, err := gnv.VerifyOne(name)
	if err != nil {
		log.Fatal().Err(err).Msg("verifyString")
	}

	f := gnv.Config().Format
	if f == gnfmt.CSV || f == gnfmt.TSV {
		fmt.Println(output.CSVHeader(f))
	}
	fmt.Println(output.NameOutput(res, f))
}

func searchQuery(gnv gnverifier.GNverifier, s string) {
	gnq := gnquery.New()
	inp := gnq.Parse(s)
	if ds := gnv.Config().DataSources; len(ds) > 0 {
		inp.DataSources = ds
	}
	if all := gnv.Config().WithAllMatches; all {
		inp.WithAllMatches = all
	}
	res, err := gnv.Search(context.Background(), inp)
	if err != nil {
		log.Fatal().Err(err).Msg("searchQuery")
	}

	f := gnv.Config().Format
	if f == gnfmt.CSV || f == gnfmt.TSV {
		fmt.Println(output.CSVHeader(f))
	}

	for _, v := range res {
		if v.Error != "" {
			log.Warn().Msg(v.Error)
		}
		fmt.Println(output.NameOutput(v, f))
	}
}

// touchConfigFile checks if config file exists, and if not, it gets created.
func touchConfigFile(configPath string) {
	fileExists, _ := gnsys.FileExists(configPath)
	if fileExists {
		return
	}
	log.Info().Msgf("Creating config file '%s'", configPath)
	createConfig(configPath)
}

// createConfig creates config file.
func createConfig(path string) {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot create dir %s", path)
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot write to file %s", path)
	}
}
