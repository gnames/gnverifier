package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitFlags(t *testing.T) {
	// Store the original rootCmd
	originalRootCmd := rootCmd
	defer func() {
		rootCmd = originalRootCmd
	}()

	// Create a new command to test flag initialization
	cmd := &cobra.Command{
		Use: "test",
	}

	// Set our test command as rootCmd
	rootCmd = cmd

	// Debug: check rootCmd is set correctly
	assert.Equal(t, cmd, rootCmd, "rootCmd should be set to our test command")

	// Initialize flags - this will add flags to the rootCmd (which is now our test cmd)
	initFlags()

	// Verify all expected flags are present
	expectedFlags := map[string]struct{}{
		"version":         {},
		"port":            {},
		"verifier_url":    {},
		"name_field":      {},
		"all_matches":     {},
		"species_group":   {},
		"fuzzy_relaxed":   {},
		"fuzzy_uninomial": {},
		"vernaculars":     {},
		"quiet":           {},
		"capitalize":      {},
		"format":          {},
		"jobs":            {},
		"sources":         {},
	}

	// Check that all expected flags exist
	foundFlags := 0
	for flagName := range expectedFlags {
		flag := cmd.Flags().Lookup(flagName)
		if assert.NotNil(t, flag, "Flag %s should exist", flagName) {
			foundFlags++
		}
	}

	// Verify total number of flags (should match expected count)
	assert.Equal(t, len(expectedFlags), foundFlags, "Should have exactly %d flags, but found %d", len(expectedFlags), foundFlags)
}

func TestBaseFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()
	rootCmd = cmd

	baseFlags()

	// Check version flag
	versionFlag := cmd.Flags().Lookup("version")
	require.NotNil(t, versionFlag)
	assert.Equal(t, "V", versionFlag.Shorthand)
	assert.Equal(t, "false", versionFlag.DefValue)
	assert.Equal(t, "Prints version information", versionFlag.Usage)
}

func TestWebFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()
	rootCmd = cmd

	webFlags()

	// Check port flag
	portFlag := cmd.Flags().Lookup("port")
	require.NotNil(t, portFlag)
	assert.Equal(t, "p", portFlag.Shorthand)
	assert.Equal(t, "0", portFlag.DefValue)
	assert.Equal(t, "Port to run web GUI.", portFlag.Usage)
}

func TestVerificationFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()
	rootCmd = cmd

	verificationFlags()

	tests := []struct {
		name      string
		shorthand string
		defValue  string
		usage     string
	}{
		{
			name:      "verifier_url",
			shorthand: "v",
			defValue:  "",
			usage:     "URL for verification service.\n  Default: https://verifier.globalnames.org/api/v1",
		},
		{
			name:      "name_field",
			shorthand: "n",
			defValue:  "1",
			usage:     "Set position of ScientificName field, the first field is 1.",
		},
		{
			name:      "all_matches",
			shorthand: "M",
			defValue:  "false",
			usage:     "return all matched results per source, not just the best one.",
		},
		{
			name:      "species_group",
			shorthand: "G",
			defValue:  "false",
			usage:     "searching for species names also searches their species groups.",
		},
		{
			name:      "fuzzy_relaxed",
			shorthand: "R",
			defValue:  "false",
			usage:     "relaxes fuzzy matching rules, decreses max names to 50.",
		},
		{
			name:      "fuzzy_uninomial",
			shorthand: "U",
			defValue:  "false",
			usage:     "allows fuzzy matching for uninomial names.",
		},
		{
			name:      "vernaculars",
			shorthand: "r",
			defValue:  "",
			usage:     "sets languages for vernacular names search (e.g., \"eng,deu,rus\")\nlimited to 50 scientific names, try it with iNaturalist (id 180)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := cmd.Flags().Lookup(tt.name)
			require.NotNil(t, flag, "Flag %s should exist", tt.name)
			assert.Equal(t, tt.shorthand, flag.Shorthand, "Shorthand for %s", tt.name)
			assert.Equal(t, tt.defValue, flag.DefValue, "Default value for %s", tt.name)
			assert.Equal(t, tt.usage, flag.Usage, "Usage for %s", tt.name)
		})
	}
}

func TestFormatFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()
	rootCmd = cmd

	formatFlags()

	tests := []struct {
		name      string
		shorthand string
		defValue  string
		usage     string
	}{
		{
			name:      "quiet",
			shorthand: "q",
			defValue:  "false",
			usage:     "do not show progress",
		},
		{
			name:      "capitalize",
			shorthand: "c",
			defValue:  "false",
			usage:     "capitalizes first character",
		},
		{
			name:      "format",
			shorthand: "f",
			defValue:  "csv",
			usage:     "Format of the output: \"compact\", \"pretty\", \"csv\", \"tsv\".\n  compact: compact JSON,\n  pretty: pretty JSON,\n  csv: CSV (DEFAULT)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := cmd.Flags().Lookup(tt.name)
			require.NotNil(t, flag, "Flag %s should exist", tt.name)
			assert.Equal(t, tt.shorthand, flag.Shorthand, "Shorthand for %s", tt.name)
			assert.Equal(t, tt.defValue, flag.DefValue, "Default value for %s", tt.name)
			assert.Equal(t, tt.usage, flag.Usage, "Usage for %s", tt.name)
		})
	}
}

func TestPerformanceFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()
	rootCmd = cmd

	performanceFlags()

	// Check jobs flag
	jobsFlag := cmd.Flags().Lookup("jobs")
	require.NotNil(t, jobsFlag)
	assert.Equal(t, "j", jobsFlag.Shorthand)
	assert.Equal(t, "4", jobsFlag.DefValue)
	assert.Equal(t, "Number of jobs running in parallel.", jobsFlag.Usage)
}

func TestDataSourcesFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()
	rootCmd = cmd

	dataSourcesFlags()

	// Check sources flag
	sourcesFlag := cmd.Flags().Lookup("sources")
	require.NotNil(t, sourcesFlag)
	assert.Equal(t, "s", sourcesFlag.Shorthand)
	assert.Equal(t, "", sourcesFlag.DefValue)

	expectedUsage := `IDs of important data-sources to verify against (ex "1,11").
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
  180 - iNaturalist
  181 - IRMNG
  194 - PLAZI
  195 - AlgaeBase`

	assert.Equal(t, expectedUsage, sourcesFlag.Usage)
}

func TestFlagTypes(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()
	rootCmd = cmd

	initFlags()

	// Test flag types
	flagTypes := map[string]string{
		"version":         "bool",
		"port":            "int",
		"verifier_url":    "string",
		"name_field":      "int",
		"all_matches":     "bool",
		"species_group":   "bool",
		"fuzzy_relaxed":   "bool",
		"fuzzy_uninomial": "bool",
		"vernaculars":     "string",
		"quiet":           "bool",
		"capitalize":      "bool",
		"format":          "string",
		"jobs":            "int",
		"sources":         "string",
	}

	for flagName, expectedType := range flagTypes {
		t.Run(flagName, func(t *testing.T) {
			flag := cmd.Flags().Lookup(flagName)
			require.NotNil(t, flag, "Flag %s should exist", flagName)
			assert.Equal(t, expectedType, flag.Value.Type(), "Flag %s should be of type %s", flagName, expectedType)
		})
	}
}

func TestFlagShorthands(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()
	rootCmd = cmd

	initFlags()

	// Test that all expected shorthands are unique and present
	expectedShorthands := map[string]string{
		"V": "version",
		"p": "port",
		"v": "verifier_url",
		"n": "name_field",
		"M": "all_matches",
		"G": "species_group",
		"R": "fuzzy_relaxed",
		"U": "fuzzy_uninomial",
		"r": "vernaculars",
		"q": "quiet",
		"c": "capitalize",
		"f": "format",
		"j": "jobs",
		"s": "sources",
	}

	for shorthand, flagName := range expectedShorthands {
		t.Run(shorthand, func(t *testing.T) {
			flag := cmd.Flags().ShorthandLookup(shorthand)
			require.NotNil(t, flag, "Shorthand %s should exist", shorthand)
			assert.Equal(t, flagName, flag.Name, "Shorthand %s should map to flag %s", shorthand, flagName)
		})
	}
}

func TestFlagDefaults(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	originalRootCmd := rootCmd
	defer func() { rootCmd = originalRootCmd }()
	rootCmd = cmd

	initFlags()

	// Test default values
	defaults := map[string]interface{}{
		"version":         false,
		"port":            0,
		"verifier_url":    "",
		"name_field":      1,
		"all_matches":     false,
		"species_group":   false,
		"fuzzy_relaxed":   false,
		"fuzzy_uninomial": false,
		"vernaculars":     "",
		"quiet":           false,
		"capitalize":      false,
		"format":          "csv",
		"jobs":            4,
		"sources":         "",
	}

	for flagName, expectedDefault := range defaults {
		t.Run(flagName, func(t *testing.T) {
			flag := cmd.Flags().Lookup(flagName)
			require.NotNil(t, flag, "Flag %s should exist", flagName)

			switch v := expectedDefault.(type) {
			case bool:
				actual, err := cmd.Flags().GetBool(flagName)
				require.NoError(t, err)
				assert.Equal(t, v, actual, "Default value for %s", flagName)
			case int:
				actual, err := cmd.Flags().GetInt(flagName)
				require.NoError(t, err)
				assert.Equal(t, v, actual, "Default value for %s", flagName)
			case string:
				actual, err := cmd.Flags().GetString(flagName)
				require.NoError(t, err)
				assert.Equal(t, v, actual, "Default value for %s", flagName)
			}
		})
	}
}
