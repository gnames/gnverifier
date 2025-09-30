// Package cmd creates a command line interface for gnverifier app.
package cmd

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnquery"
	"github.com/gnames/gnquery/ent/search"
	"github.com/gnames/gnsys"
	gnverifier "github.com/gnames/gnverifier/pkg"
	"github.com/gnames/gnverifier/pkg/config"
	"github.com/gnames/gnverifier/pkg/ent/output"
	"github.com/gnames/gnverifier/pkg/io/verifrest"
	"github.com/gnames/gnverifier/pkg/io/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed gnverifier.yaml
var configText string

var (
	opts, webOpts []config.Option
	msg           []string
)

// cfgData purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type cfgData struct {
	DataSources             []int
	Format                  string
	Jobs                    int
	VerifierURL             string
	WithAllMatches          bool
	WithCapitalization      bool
	WithSpeciesGroup        bool
	WithUninomialFuzzyMatch bool
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
		// if there is version flag, show version and exit
		versionFlag(cmd)

		flags := []funcFlag{
			capitalizeFlag, spGroupFlag, fuzzyRelaxedFlag,
			fuzzyUninomialFlag, formatFlag, jobsFlag, allMatchesFlag,
			sourcesFlag, vernacularsFlag, verifierUrlFlag, quietFlag,
		}

		for _, f := range flags {
			f(cmd)
		}

		// these will only show if quietFlag was false
		for i := range msg {
			slog.Info(msg[i])
		}

		// if port is given run gnverifier web UI instead
		port, _ := cmd.Flags().GetInt("port")
		if port > 0 {
			// Copy opts to webOpts after flags have been processed
			webOpts = append([]config.Option{}, opts...)
			webOpts = append(webOpts, config.OptWithCapitalization(true))
			cfg := config.New(webOpts...)
			vfr := verifrest.New(cfg.VerifierURL)
			gnv := gnverifier.New(cfg, vfr)
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
	initFlags()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var configDir string
	var err error
	configFile := "gnverifier"

	// Find config directory.
	configDir, err = os.UserConfigDir()
	if err != nil {
		slog.Error("Cannot find config directory", "error", err)
		os.Exit(1)
	}

	// Search config in home directory with name ".gnmatcher" (without extension).
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFile)

	// Set environment variables to override
	// config file settings
	_ = viper.BindEnv("DataSources", "GNV_DATA_SOURCES")
	_ = viper.BindEnv("Format", "GNV_FORMAT")
	_ = viper.BindEnv("Jobs", "GNV_JOBS")
	_ = viper.BindEnv("VerifierURL", "GNV_VERIFIER_URL")
	_ = viper.BindEnv("WithAllMatches", "GNV_WITH_ALL_MATCHES")
	_ = viper.BindEnv("WithCapitalization", "GNV_WITH_CAPITALIZATION")
	_ = viper.BindEnv("WithSpeciesGroup", "GNV_WITH_SPECIES_GROUP")

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
		slog.Error("Cannot deserialize config data", "error", err)
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
		slog.Error("Cannot get Stdin info", "error", err)
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
			slog.Error("Cannot open file", "error", err, "file", str)
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

		slog.Info("Verified.",
			"names/sec", humanize.Comma(speed),
			"names", humanize.Comma(int64(total)),
		)
		for _, r := range o {
			if r.Error != "" {
				slog.Error("Error during verification", "error", r.Error)
			}
			fmt.Println(output.NameOutput(r, f))
		}
	}
}

func verifyString(gnv gnverifier.GNverifier, name string) {
	res, err := gnv.VerifyOne(name)
	if err != nil {
		slog.Error("Cannot verify name", "error", err, "name", name)
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
		slog.Error("Cannot run search query", "error", err, "input", inp)
	}

	f := gnv.Config().Format
	if f == gnfmt.CSV || f == gnfmt.TSV {
		fmt.Println(output.CSVHeader(f))
	}

	for _, v := range res {
		if v.Error != "" {
			slog.Error("Error during search", "error", v.Error)
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
	slog.Info("Creating config file", "file", configPath)
	createConfig(configPath)
}

// createConfig creates config file.
func createConfig(path string) {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		slog.Error("Cannot create dir", "dir", path, "error", err)
		os.Exit(1)
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		slog.Error("Cannot write to file", "path", path, "error", err)
		os.Exit(1)
	}
}
