// Package cmd creates a command line interface for gnverifier app.
package cmd

import (
	"bufio"
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
	"github.com/gnames/gnsys"
	"github.com/gnames/gnverifier"
	"github.com/gnames/gnverifier/config"
	"github.com/gnames/gnverifier/ent/output"
	"github.com/gnames/gnverifier/io/verifrest"
	"github.com/gnames/gnverifier/io/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed gnverifier.yaml
var configText string

var (
	opts []config.Option
)

// cfgData purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type cfgData struct {
	Format             string
	PreferredOnly      bool
	PreferredSources   []int
	WithCapitalization bool
	VerifierURL        string
	Jobs               int
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnverifier",
	Short: "Verifies scientific names agains many sources.",
	Long: `gnverifier uses a remote service to verify scientific names against
more than 100 biodiverisity data-sources.`,
	Run: func(cmd *cobra.Command, args []string) {
		webOpts := make([]config.Option, len(opts))
		for i, v := range opts {
			webOpts[i] = v
		}
		webOpts = append(webOpts, config.OptWithCapitalization(true))

		if showVersionFlag(cmd) {
			os.Exit(0)
		}

		caps, _ := cmd.Flags().GetBool("capitalize")
		opts = append(opts, config.OptWithCapitalization(caps))

		pref, _ := cmd.Flags().GetBool("only_preferred")
		opts = append(opts, config.OptPreferredOnly(pref))

		formatString, _ := cmd.Flags().GetString("format")
		frmt, _ := gnfmt.NewFormat(formatString)
		if frmt == gnfmt.FormatNone {
			log.Warnf("Cannot set format from '%s', setting format to csv", formatString)
			frmt = gnfmt.CSV
		}
		opts = append(opts, config.OptFormat(frmt))

		jobs, _ := cmd.Flags().GetInt("jobs")
		if jobs != 4 {
			opts = append(opts, config.OptJobs(jobs))
		}

		sources, _ := cmd.Flags().GetString("sources")
		if sources != "" {
			data_sources := parseDataSources(sources)
			opts = append(opts, config.OptPreferredSources(data_sources))
		}

		url, _ := cmd.Flags().GetString("verifier_url")
		if url != "" {
			opts = append(opts, config.OptVerifierURL(url))
			webOpts = append(webOpts, config.OptVerifierURL(url))
		}

		port, _ := cmd.Flags().GetInt("port")
		if port > 0 {
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
	rootCmd.Flags().BoolP("only_preferred", "o", false, "Ignores best match, returns only preferred results (if any).")
	rootCmd.Flags().StringP("format", "f", "csv", `Format of the output: "compact", "pretty", "csv".
  compact: compact JSON,
  pretty: pretty JSON,
  csv: CSV (DEFAULT)`)
	rootCmd.Flags().IntP("name_field", "n", 1, "Set position of ScientificName field, the first field is 1.")
	rootCmd.Flags().IntP("jobs", "j", 4, "Number of jobs running in parallel.")
	rootCmd.Flags().IntP("port", "p", 0, "Port to run web GUI.")
	rootCmd.Flags().StringP("sources", "s", "", `IDs of important data-sources to verify against (ex "1,11").
  If sources are set and there are matches to their data,
  such matches are returned in "preferred_result" results.
  To find IDs refer to "https://resolver.globalnames.org/data_sources".
  1 - Catalogue of Life
  3 - ITIS
  4 - NCBI
  9 - WoRMS
  11 - GBIF
  12 - Encyclopedia of Life
  167 - IPNI
  170 - Arctos
  172 - PaleoBioDB
  181 - IRMNG`)
	rootCmd.Flags().StringP("verifier_url", "v", "", "URL for verification service")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var configDir string
	var err error
	configFile := "gnverifier"

	// Find config directory.
	configDir, err = os.UserConfigDir()
	if err != nil {
		log.Fatalf("Cannot find config directory: %s.", err)
	}

	// Search config in home directory with name ".gnmatcher" (without extension).
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFile)

	configPath := filepath.Join(configDir, fmt.Sprintf("%s.yaml", configFile))
	touchConfigFile(configPath)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s.", viper.ConfigFileUsed())
	}
	getOpts()
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts() {
	cfg := &cfgData{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("Cannot deserialize config data: %s.", err)
	}

	if cfg.Format != "" {
		cfgFormat, err := gnfmt.NewFormat(cfg.Format)
		if err != nil {
			cfgFormat = gnfmt.CSV
		}
		opts = append(opts, config.OptFormat(cfgFormat))
	}
	if cfg.PreferredOnly {
		opts = append(opts, config.OptPreferredOnly(cfg.PreferredOnly))
	}
	if len(cfg.PreferredSources) > 0 {
		opts = append(opts, config.OptPreferredSources(cfg.PreferredSources))
	}
	if cfg.VerifierURL != "" {
		opts = append(opts, config.OptVerifierURL(cfg.VerifierURL))
	}
	if cfg.Jobs > 0 {
		opts = append(opts, config.OptJobs(cfg.Jobs))
	}
}

// showVersionFlag provides version and the build timestamp. If it returns
// true, it means that version flag was given.
func showVersionFlag(cmd *cobra.Command) bool {
	hasVersionFlag, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Fatalf("Cannot get version flag: %s.", err)
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
			log.Warnf("Cannot convert data-source '%s' to list, skipping", v)
			return nil
		}
		if ds < 1 {
			log.Warnf("Data source ID %d is less than one, skipping", ds)
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
		log.Panic(err)
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

func verify(gnv gnverifier.GNverifier, data string) {
	path := string(data)
	fileExists, _ := gnsys.FileExists(path)
	if fileExists {
		f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		verifyFile(gnv, f)
		f.Close()
	} else {
		verifyString(gnv, data)
	}
}

func verifyFile(gnv gnverifier.GNverifier, f io.Reader) {
	batch := gnv.Config().Batch
	in := make(chan []string)
	out := make(chan []vlib.Verification)
	var wg sync.WaitGroup
	wg.Add(1)
	go gnv.VerifyStream(in, out)
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

func processResults(gnv gnverifier.GNverifier, out <-chan []vlib.Verification,
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

		log.Printf("Verified %s records, %s names/sec\n", humanize.Comma(total),
			humanize.Comma(speed))
		for _, r := range o {
			if r.Error != "" {
				log.Println(r.Error)
			}
			fmt.Println(output.Output(r, f, gnv.Config().PreferredOnly))
		}
	}
}

func verifyString(gnv gnverifier.GNverifier, name string) {
	res, err := gnv.VerifyOne(name)
	if err != nil {
		log.Fatal(err)
	}

	f := gnv.Config().Format
	if f == gnfmt.CSV || f == gnfmt.TSV {
		fmt.Println(output.CSVHeader(f))
	}
	fmt.Println(output.Output(res, f, gnv.Config().PreferredOnly))
}

// touchConfigFile checks if config file exists, and if not, it gets created.
func touchConfigFile(configPath string) {
	fileExists, _ := gnsys.FileExists(configPath)
	if fileExists {
		return
	}

	log.Printf("Creating config file: %s.", configPath)
	createConfig(configPath)
}

// createConfig creates config file.
func createConfig(path string) {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		log.Fatalf("Cannot create dir %s: %s.", path, err)
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		log.Fatalf("Cannot write to file %s: %s", path, err)
	}
}
