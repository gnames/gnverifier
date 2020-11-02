/*
Copyright © 2020 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	gne "github.com/gnames/gnames/domain/entity"
	"github.com/gnames/gnlib/format"
	"github.com/gnames/gnlib/sys"
	"github.com/gnames/gnverify"
	"github.com/gnames/gnverify/config"
	"github.com/gnames/gnverify/output"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	opts []config.Option
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnverify",
	Short: "Verifies scientific names agains many sources.",
	Long: `gnverify uses a remote service to verify scientific names against
more than 100 biodiverisity data-sources.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersionFlag(cmd) {
			os.Exit(0)
		}
		pref, _ := cmd.Flags().GetBool("preferred_only")
		opts = append(opts, config.OptPreferredOnly(pref))

		formatString, _ := cmd.Flags().GetString("format")
		frmt, _ := format.NewFormat(formatString)
		if frmt == format.FormatNone {
			log.Warnf("Cannot set format from '%s', setting format to csv", formatString)
			frmt = format.CSV
		}
		opts = append(opts, config.OptFormat(frmt))

		name_field, err := cmd.Flags().GetInt("name_field")
		if err != nil {
			log.Warnf("Cannot set position of the name_field: %s", err)
			name_field = 1
		}
		if name_field < 1 {
			log.Warnf("Cannot set name_field index because %d is less than 1", name_field)
			name_field = 1
		}
		opts = append(opts, config.OptNameField(uint(name_field-1)))

		sources, _ := cmd.Flags().GetString("sources")
		data_sources := parseDataSources(sources)
		opts = append(opts, config.OptPreferredSources(data_sources))

		url, _ := cmd.Flags().GetString("verifier_url")
		if len(url) > 0 {
			opts = append(opts, config.OptVerifierURL(url))
		}

		cnf := config.NewConfig(opts...)

		if len(args) == 0 {
			processStdin(cmd, cnf)
			os.Exit(0)
		}
		data := getInput(cmd, args)
		gnv := gnverify.NewGNVerify(cnf)
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
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "V", false, "Prints version information")
	rootCmd.Flags().BoolP("preferred_only", "p", false, "Ignores best match, returns only preferred results (if any).")
	rootCmd.Flags().StringP("format", "f", "csv", `Format of the output: "compact", "pretty", "csv".
  compact: compact JSON,
  pretty: pretty JSON,
  csv: CSV (DEFAULT)`)
	rootCmd.Flags().IntP("name_field", "n", 1, "Set position of ScientificName field, the first field is 1.")
	rootCmd.Flags().StringP("sources", "s", "", `IDs of important data-sources to verify against (ex "1,11").
  If sources are set and there are matches to their data,
  such matches are returned in "preferred_result" results.
  To find IDs refer to "https://resolver.globalnames.org/resources".
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

// showVersionFlag provides version and the build timestamp. If it returns
// true, it means that version flag was given.
func showVersionFlag(cmd *cobra.Command) bool {
	hasVersionFlag, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Fatalf("Cannot get version flag: %s.", err)
	}

	if hasVersionFlag {
		fmt.Printf("\nversion: %s\nbuild: %s\n\n", gnverify.Version, gnverify.Build)
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

func processStdin(cmd *cobra.Command, cnf config.Config) {
	if !checkStdin() {
		_ = cmd.Help()
		return
	}
	gnv := gnverify.NewGNVerify(cnf)
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

func verify(gnv gnverify.GNVerify, data string) {
	path := string(data)
	if sys.FileExists(path) {
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

func verifyFile(gnv gnverify.GNVerify, f io.Reader) {
	batch := 5000
	in := make(chan []string)
	out := make(chan []*gne.Verification)
	var wg sync.WaitGroup
	wg.Add(1)

	go gnv.VerifyStream(in, out)
	go processResults(gnv, out, &wg)
	sc := bufio.NewScanner(f)
	count := 0
	names := make([]string, 0, batch)
	for sc.Scan() {
		count++
		if count%50000 == 0 {
			log.Printf("Verifying %d-th line\n", count)
		}

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

func processResults(gnv gnverify.GNVerify, out <-chan []*gne.Verification,
	wg *sync.WaitGroup) {
	defer wg.Done()
	if gnv.Format == format.CSV {
		fmt.Println(output.CSVHeader())
	}
	for o := range out {
		for _, r := range o {
			if r.Error != "" {
				log.Println(r.Error)
			}
			fmt.Println(output.Output(r, gnv.Format, gnv.PreferredOnly))

		}
	}
}

func verifyString(gnv gnverify.GNVerify, name string) {
	res := gnv.Verify(name)
	if gnv.Format == format.CSV {
		fmt.Println(output.CSVHeader())
	}
	fmt.Println(res)
}
