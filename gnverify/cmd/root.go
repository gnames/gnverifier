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
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/gogna/gnverify"
	gncnf "gitlab.com/gogna/gnverify/config"
)

var (
	opts []gncnf.Option
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
		opts = append(opts, gncnf.OptPreferredOnly(pref))

		formatString, _ := cmd.Flags().GetString("format")
		format := gncnf.NewFormat(formatString)
		if format == gncnf.InvalidFormat {
			log.Warnf("Cannot set format from '%s', setting format to csv")
			format = gncnf.CSV
		}
		opts = append(opts, gncnf.OptFormat(format))

		name_field, err := cmd.Flags().GetInt("name_field")
		if err != nil {
			log.Warnf("Cannot set position of the name_field: %s", err)
			name_field = 1
		}
		if name_field < 1 {
			log.Warnf("Cannot set name_field index because %d is less than 1", name_field)
			name_field = 1
		}
		opts = append(opts, gncnf.OptNameField(uint(name_field-1)))

		sources, _ := cmd.Flags().GetString("sources")
		data_sources := parseDataSources(sources)
		opts = append(opts, gncnf.OptPreferredSources(data_sources))

		url, _ := cmd.Flags().GetString("verifier_url")
		if len(url) > 0 {
			opts = append(opts, gncnf.OptVerifierURL(url))
		}

		cnf := gncnf.NewConfig(opts...)

		_ = gnverify.NewGNVerify(cnf)
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

func parseDataSources(s string) []uint {
	dss := strings.Split(s, ",")
	res := make([]uint, 0, len(dss))
	for _, v := range dss {
		v = strings.Trim(v, " ")
		ds, err := strconv.Atoi(v)
		if err != nil {
			log.Warnf("Cannot convert data-source '%s' to list, skipping")
			return nil
		}
		if ds < 1 {
			log.Warnf("Data source ID %d is less than one, skipping", ds)
		} else {
			res = append(res, uint(ds))
		}
	}
	if len(res) > 0 {
		return res
	}
	return nil
}
