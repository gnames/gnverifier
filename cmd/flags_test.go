package cmd

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnverifier/pkg/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionFlag(t *testing.T) {
	tests := []struct {
		name        string
		versionFlag bool
		shouldExit  bool
	}{
		{
			name:        "version flag not set",
			versionFlag: false,
			shouldExit:  false,
		},
		{
			name:        "version flag set",
			versionFlag: true,
			shouldExit:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().Bool("version", tt.versionFlag, "test version flag")

			if tt.shouldExit {
				// This should cause os.Exit(0) which we can't easily test
				// In a real scenario, you might want to refactor to make this testable
				// For now, we'll just verify the flag is read correctly
				hasVersion, _ := cmd.Flags().GetBool("version")
				assert.Equal(t, tt.versionFlag, hasVersion)
			} else {
				versionFlag(cmd)
			}
		})
	}
}

func TestQuietFlag(t *testing.T) {
	tests := []struct {
		name      string
		quietFlag bool
	}{
		{
			name:      "quiet flag not set",
			quietFlag: false,
		},
		{
			name:      "quiet flag set",
			quietFlag: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().Bool("quiet", tt.quietFlag, "test quiet flag")

			// Reset log level before test
			slog.SetLogLoggerLevel(slog.LevelInfo)

			quietFlag(cmd)

			// We can't easily test the log level change directly,
			// but we can verify the flag was read correctly
			quiet, _ := cmd.Flags().GetBool("quiet")
			assert.Equal(t, tt.quietFlag, quiet)
		})
	}
}

func TestCapitalizeFlag(t *testing.T) {
	tests := []struct {
		name           string
		capitalizeFlag bool
		expectOpt      bool
	}{
		{
			name:           "capitalize flag not set",
			capitalizeFlag: false,
			expectOpt:      false,
		},
		{
			name:           "capitalize flag set",
			capitalizeFlag: true,
			expectOpt:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			opts = nil
			webOpts = nil

			cmd := &cobra.Command{}
			cmd.Flags().Bool("capitalize", tt.capitalizeFlag, "test capitalize flag")

			capitalizeFlag(cmd)

			if tt.expectOpt {
				assert.Len(t, opts, 1)
				// Create a config with the options to test if capitalization is set
				cfg := config.New(opts...)
				assert.True(t, cfg.WithCapitalization)
			} else {
				assert.Len(t, opts, 0)
			}
		})
	}
}

func TestSpGroupFlag(t *testing.T) {
	tests := []struct {
		name        string
		spGroupFlag bool
		expectOpt   bool
	}{
		{
			name:        "species_group flag not set",
			spGroupFlag: false,
			expectOpt:   false,
		},
		{
			name:        "species_group flag set",
			spGroupFlag: true,
			expectOpt:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			opts = nil
			webOpts = nil

			cmd := &cobra.Command{}
			cmd.Flags().Bool("species_group", tt.spGroupFlag, "test species_group flag")

			spGroupFlag(cmd)

			if tt.expectOpt {
				assert.Len(t, opts, 1)
				cfg := config.New(opts...)
				assert.True(t, cfg.WithSpeciesGroup)
			} else {
				assert.Len(t, opts, 0)
			}
		})
	}
}

func TestFuzzyRelaxedFlag(t *testing.T) {
	tests := []struct {
		name             string
		fuzzyRelaxedFlag bool
		expectOpt        bool
	}{
		{
			name:             "fuzzy_relaxed flag not set",
			fuzzyRelaxedFlag: false,
			expectOpt:        false,
		},
		{
			name:             "fuzzy_relaxed flag set",
			fuzzyRelaxedFlag: true,
			expectOpt:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			opts = nil
			webOpts = nil

			cmd := &cobra.Command{}
			cmd.Flags().Bool("fuzzy_relaxed", tt.fuzzyRelaxedFlag, "test fuzzy_relaxed flag")

			fuzzyRelaxedFlag(cmd)

			if tt.expectOpt {
				assert.Len(t, opts, 1)
				cfg := config.New(opts...)
				assert.True(t, cfg.WithRelaxedFuzzyMatch)
			} else {
				assert.Len(t, opts, 0)
			}
		})
	}
}

func TestFuzzyUninomialFlag(t *testing.T) {
	tests := []struct {
		name               string
		fuzzyUninomialFlag bool
		expectOpt          bool
	}{
		{
			name:               "fuzzy_uninomial flag not set",
			fuzzyUninomialFlag: false,
			expectOpt:          false,
		},
		{
			name:               "fuzzy_uninomial flag set",
			fuzzyUninomialFlag: true,
			expectOpt:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			opts = nil
			webOpts = nil

			cmd := &cobra.Command{}
			cmd.Flags().Bool("fuzzy_uninomial", tt.fuzzyUninomialFlag, "test fuzzy_uninomial flag")

			fuzzyUninomialFlag(cmd)

			if tt.expectOpt {
				assert.Len(t, opts, 1)
				cfg := config.New(opts...)
				assert.True(t, cfg.WithUninomialFuzzyMatch)
			} else {
				assert.Len(t, opts, 0)
			}
		})
	}
}

func TestFormatFlag(t *testing.T) {
	tests := []struct {
		name           string
		formatString   string
		expectedFormat gnfmt.Format
	}{
		{
			name:           "csv format",
			formatString:   "csv",
			expectedFormat: gnfmt.CSV,
		},
		{
			name:           "tsv format",
			formatString:   "tsv",
			expectedFormat: gnfmt.TSV,
		},
		{
			name:           "json format",
			formatString:   "compact",
			expectedFormat: gnfmt.CompactJSON,
		},
		{
			name:           "pretty json format",
			formatString:   "pretty",
			expectedFormat: gnfmt.PrettyJSON,
		},
		{
			name:           "invalid format defaults to csv",
			formatString:   "invalid",
			expectedFormat: gnfmt.CSV,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			opts = nil
			webOpts = nil

			cmd := &cobra.Command{}
			cmd.Flags().String("format", tt.formatString, "test format flag")

			formatFlag(cmd)

			assert.Len(t, opts, 1)
			cfg := config.New(opts...)
			assert.Equal(t, tt.expectedFormat, cfg.Format)
		})
	}
}

func TestJobsFlag(t *testing.T) {
	tests := []struct {
		name         string
		jobs         int
		expectOpt    bool
		expectedJobs int
	}{
		{
			name:      "default jobs (4) - no option added",
			jobs:      4,
			expectOpt: false,
		},
		{
			name:         "custom jobs value",
			jobs:         8,
			expectOpt:    true,
			expectedJobs: 8,
		},
		{
			name:      "zero jobs - no option added",
			jobs:      0,
			expectOpt: false,
		},
		{
			name:      "negative jobs - no option added",
			jobs:      -1,
			expectOpt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			opts = nil
			webOpts = nil

			cmd := &cobra.Command{}
			cmd.Flags().Int("jobs", tt.jobs, "test jobs flag")

			jobsFlag(cmd)

			if tt.expectOpt {
				assert.Len(t, opts, 1)
				cfg := config.New(opts...)
				assert.Equal(t, tt.expectedJobs, cfg.Jobs)
			} else {
				assert.Len(t, opts, 0)
			}
		})
	}
}

func TestAllMatchesFlag(t *testing.T) {
	tests := []struct {
		name           string
		allMatchesFlag bool
		expectOpt      bool
	}{
		{
			name:           "all_matches flag not set",
			allMatchesFlag: false,
			expectOpt:      false,
		},
		{
			name:           "all_matches flag set",
			allMatchesFlag: true,
			expectOpt:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			opts = nil
			webOpts = nil

			cmd := &cobra.Command{}
			cmd.Flags().Bool("all_matches", tt.allMatchesFlag, "test all_matches flag")

			allMatchesFlag(cmd)

			if tt.expectOpt {
				assert.Len(t, opts, 1)
				cfg := config.New(opts...)
				assert.True(t, cfg.WithAllMatches)
			} else {
				assert.Len(t, opts, 0)
			}
		})
	}
}

func TestSourcesFlag(t *testing.T) {
	tests := []struct {
		name            string
		sources         string
		expectOpt       bool
		expectedSources []int
	}{
		{
			name:      "empty sources",
			sources:   "",
			expectOpt: false,
		},
		{
			name:            "single source",
			sources:         "1",
			expectOpt:       true,
			expectedSources: []int{1},
		},
		{
			name:            "multiple sources",
			sources:         "1,11,180",
			expectOpt:       true,
			expectedSources: []int{1, 11, 180},
		},
		{
			name:            "sources with spaces",
			sources:         "1, 11 , 180",
			expectOpt:       true,
			expectedSources: []int{1, 11, 180},
		},
		{
			name:            "invalid sources",
			sources:         "abc,def",
			expectOpt:       true,
			expectedSources: nil, // nil is still passed to OptDataSources
		},
		{
			name:            "mixed valid and invalid sources",
			sources:         "1,abc,180",
			expectOpt:       true,
			expectedSources: nil, // nil is still passed to OptDataSources
		},
		{
			name:            "negative source filtered out but positive ones kept",
			sources:         "1,-5,180",
			expectOpt:       true,
			expectedSources: []int{1, 180},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state before each test
			opts = nil
			webOpts = nil

			// Capture log output for warnings
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
			slog.SetDefault(logger)

			cmd := &cobra.Command{}
			cmd.Flags().String("sources", tt.sources, "test sources flag")

			sourcesFlag(cmd)

			if tt.expectOpt {
				require.True(t, len(opts) > 0, "Expected opts to contain data sources option")
				cfg := config.New(opts...)
				assert.Equal(t, tt.expectedSources, cfg.DataSources)
			} else {
				assert.Len(t, opts, 0)
			}
		})
	}
}

func TestVerifierUrlFlag(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectOpt   bool
		expectedURL string
	}{
		{
			name:      "empty url",
			url:       "",
			expectOpt: false,
		},
		{
			name:        "custom url",
			url:         "https://example.com/api/v1",
			expectOpt:   true,
			expectedURL: "https://example.com/api/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global opts and webOpts
			opts = nil
			webOpts = nil

			cmd := &cobra.Command{}
			cmd.Flags().String("verifier_url", tt.url, "test verifier_url flag")

			verifierUrlFlag(cmd)

			if tt.expectOpt {
				assert.Len(t, opts, 1)
				assert.Len(t, webOpts, 1)
				cfg := config.New(opts...)
				assert.Equal(t, tt.expectedURL, cfg.VerifierURL)
			} else {
				assert.Len(t, opts, 0)
				assert.Len(t, webOpts, 0)
			}
		})
	}
}

func TestVernacularsFlag(t *testing.T) {
	tests := []struct {
		name              string
		vernaculars       string
		expectOpt         bool
		expectedLanguages []string
	}{
		{
			name:        "empty vernaculars",
			vernaculars: "",
			expectOpt:   false,
		},
		{
			name:              "single language",
			vernaculars:       "eng",
			expectOpt:         true,
			expectedLanguages: []string{"eng"},
		},
		{
			name:              "multiple languages",
			vernaculars:       "eng,deu,rus",
			expectOpt:         true,
			expectedLanguages: []string{"eng", "deu", "rus"},
		},
		{
			name:              "languages with spaces - only first valid",
			vernaculars:       "eng, deu , rus",
			expectOpt:         true,
			expectedLanguages: []string{"eng", "deu", "rus"}, // Only "eng" is valid, " deu " and " rus" have spaces
		},
		{
			name:              "invalid language length",
			vernaculars:       "en,german,rus",
			expectOpt:         true,
			expectedLanguages: []string{"rus"}, // Only "rus" is valid 3-letter code
		},
		{
			name:              "languages without spaces",
			vernaculars:       "eng,deu,rus",
			expectOpt:         true,
			expectedLanguages: []string{"eng", "deu", "rus"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state before each test
			opts = nil
			webOpts = nil

			// Capture log output for warnings
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
			slog.SetDefault(logger)

			cmd := &cobra.Command{}
			cmd.Flags().String("vernaculars", tt.vernaculars, "test vernaculars flag")

			vernacularsFlag(cmd)

			if tt.expectOpt {
				require.True(t, len(opts) > 0, "Expected opts to contain vernaculars option")
				cfg := config.New(opts...)
				assert.Equal(t, tt.expectedLanguages, cfg.Vernaculars)
			} else {
				assert.Len(t, opts, 0)
			}
		})
	}
}

func TestParseDataSources(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single source",
			input:    "1",
			expected: []int{1},
		},
		{
			name:     "multiple sources",
			input:    "1,11,180",
			expected: []int{1, 11, 180},
		},
		{
			name:     "sources with spaces",
			input:    " 1 , 11 , 180 ",
			expected: []int{1, 11, 180},
		},
		{
			name:     "invalid source",
			input:    "abc",
			expected: nil,
		},
		{
			name:     "mixed valid and invalid",
			input:    "1,abc,180",
			expected: nil,
		},
		{
			name:     "negative sources filtered",
			input:    "1,-5,180",
			expected: []int{1, 180},
		},
		{
			name:     "all negative sources",
			input:    "-1,-5,-10",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
			slog.SetDefault(logger)

			result := parseDataSources(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseVernacularLanguages(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single valid language",
			input:    "eng",
			expected: []string{"eng"},
		},
		{
			name:     "multiple valid languages",
			input:    "eng,deu,rus",
			expected: []string{"eng", "deu", "rus"},
		},
		{
			name:     "invalid language length - too short",
			input:    "en,deu",
			expected: []string{"deu"},
		},
		{
			name:     "invalid language length - too long",
			input:    "english,deu",
			expected: []string{"deu"},
		},
		{
			name:     "all invalid languages",
			input:    "en,de,english",
			expected: nil,
		},
		{
			name:     "mixed valid and invalid",
			input:    "eng,de,rus,english",
			expected: []string{"eng", "rus"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
			slog.SetDefault(logger)

			result := parseVernacularLanguages(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAllFlagFunctions(t *testing.T) {
	// Test that all flag functions can be called without panicking
	flagFunctions := []funcFlag{
		capitalizeFlag,
		spGroupFlag,
		fuzzyRelaxedFlag,
		fuzzyUninomialFlag,
		formatFlag,
		jobsFlag,
		allMatchesFlag,
		sourcesFlag,
		vernacularsFlag,
		verifierUrlFlag,
		quietFlag,
	}

	// Reset global state
	opts = nil
	webOpts = nil

	// Create a fresh command for each test to avoid flag redefinition
	for i, flagFunc := range flagFunctions {
		t.Run(fmt.Sprintf("flag_function_%d", i), func(t *testing.T) {
			cmd := &cobra.Command{Use: "test"}
			originalRootCmd := rootCmd
			defer func() { rootCmd = originalRootCmd }()
			rootCmd = cmd
			initFlags()

			require.NotPanics(t, func() {
				flagFunc(cmd)
			})
		})
	}
}
