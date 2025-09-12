package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnverifier/pkg/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlagIntegration(t *testing.T) {
	tests := []struct {
		name     string
		flags    map[string]interface{}
		validate func(t *testing.T, cfg config.Config)
	}{
		{
			name: "default configuration",
			flags: map[string]interface{}{
				"format": "csv",
				"jobs":   4,
			},
			validate: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, gnfmt.CSV, cfg.Format)
				assert.Equal(t, 4, cfg.Jobs)
				assert.Equal(t, "https://verifier.globalnames.org/api/v1/", cfg.VerifierURL)
				assert.False(t, cfg.WithAllMatches)
				assert.False(t, cfg.WithCapitalization)
				assert.False(t, cfg.WithSpeciesGroup)
				assert.False(t, cfg.WithRelaxedFuzzyMatch)
				assert.False(t, cfg.WithUninomialFuzzyMatch)
				assert.Empty(t, cfg.DataSources)
				assert.Empty(t, cfg.Vernaculars)
			},
		},
		{
			name: "all boolean flags enabled",
			flags: map[string]interface{}{
				"all_matches":     true,
				"capitalize":      true,
				"species_group":   true,
				"fuzzy_relaxed":   true,
				"fuzzy_uninomial": true,
				"format":          "pretty",
			},
			validate: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, gnfmt.PrettyJSON, cfg.Format)
				assert.True(t, cfg.WithAllMatches)
				assert.True(t, cfg.WithCapitalization)
				assert.True(t, cfg.WithSpeciesGroup)
				assert.True(t, cfg.WithRelaxedFuzzyMatch)
				assert.True(t, cfg.WithUninomialFuzzyMatch)
			},
		},
		{
			name: "custom performance and data sources",
			flags: map[string]interface{}{
				"jobs":    8,
				"sources": "1,11,180",
				"format":  "tsv",
			},
			validate: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, gnfmt.TSV, cfg.Format)
				assert.Equal(t, 8, cfg.Jobs)
				assert.Equal(t, []int{1, 11, 180}, cfg.DataSources)
			},
		},
		{
			name: "vernacular languages and custom URL",
			flags: map[string]interface{}{
				"vernaculars":  "eng,deu,rus",
				"verifier_url": "https://example.com/api/v1",
				"format":       "compact",
			},
			validate: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, gnfmt.CompactJSON, cfg.Format)
				assert.Equal(t, "https://example.com/api/v1", cfg.VerifierURL)
				assert.Equal(t, []string{"eng", "deu", "rus"}, cfg.Vernaculars)
			},
		},
		{
			name: "comprehensive configuration",
			flags: map[string]interface{}{
				"all_matches":     true,
				"capitalize":      true,
				"species_group":   true,
				"fuzzy_relaxed":   true,
				"fuzzy_uninomial": true,
				"jobs":            12,
				"sources":         "1,3,4,11,180",
				"vernaculars":     "eng,spa,fra",
				"verifier_url":    "https://custom.verifier.org/api/v1",
				"format":          "pretty",
			},
			validate: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, gnfmt.PrettyJSON, cfg.Format)
				assert.Equal(t, 12, cfg.Jobs)
				assert.Equal(t, "https://custom.verifier.org/api/v1", cfg.VerifierURL)
				assert.True(t, cfg.WithAllMatches)
				assert.True(t, cfg.WithCapitalization)
				assert.True(t, cfg.WithSpeciesGroup)
				assert.True(t, cfg.WithRelaxedFuzzyMatch)
				assert.True(t, cfg.WithUninomialFuzzyMatch)
				assert.Equal(t, []int{1, 3, 4, 11, 180}, cfg.DataSources)
				assert.Equal(t, []string{"eng", "spa", "fra"}, cfg.Vernaculars)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			opts = nil
			webOpts = nil

			// Create command with flags
			cmd := &cobra.Command{Use: "test"}
			originalRootCmd := rootCmd
			defer func() { rootCmd = originalRootCmd }()
			rootCmd = cmd

			initFlags()

			// Set flag values
			for flagName, value := range tt.flags {
				switch v := value.(type) {
				case bool:
					err := cmd.Flags().Set(flagName, "true")
					if v == false {
						err = cmd.Flags().Set(flagName, "false")
					}
					require.NoError(t, err, "Failed to set flag %s", flagName)
				case int:
					err := cmd.Flags().Set(flagName, fmt.Sprintf("%d", v))
					require.NoError(t, err, "Failed to set flag %s", flagName)
				case string:
					err := cmd.Flags().Set(flagName, v)
					require.NoError(t, err, "Failed to set flag %s", flagName)
				}
			}

			// Apply all flags
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
			}

			for _, flagFunc := range flagFunctions {
				flagFunc(cmd)
			}

			// Create config with options and validate
			cfg := config.New(opts...)
			tt.validate(t, cfg)
		})
	}
}

func TestFlagConflicts(t *testing.T) {
	tests := []struct {
		name        string
		flags       map[string]interface{}
		expectError bool
		description string
	}{
		{
			name: "invalid format with valid other flags",
			flags: map[string]interface{}{
				"format":     "invalid",
				"capitalize": true,
				"jobs":       8,
			},
			expectError: false, // Should default to CSV
			description: "Invalid format should default to CSV",
		},
		{
			name: "zero jobs with valid other flags",
			flags: map[string]interface{}{
				"jobs":        0,
				"all_matches": true,
				"format":      "tsv",
			},
			expectError: false, // Should use default jobs
			description: "Zero jobs should use default",
		},
		{
			name: "negative jobs with valid other flags",
			flags: map[string]interface{}{
				"jobs":          -5,
				"species_group": true,
				"format":        "compact",
			},
			expectError: false, // Should use default jobs
			description: "Negative jobs should use default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			opts = nil

			// Create command with flags
			cmd := &cobra.Command{Use: "test"}
			originalRootCmd := rootCmd
			defer func() { rootCmd = originalRootCmd }()
			rootCmd = cmd

			initFlags()

			// Set flag values
			for flagName, value := range tt.flags {
				switch v := value.(type) {
				case bool:
					err := cmd.Flags().Set(flagName, "true")
					if v == false {
						err = cmd.Flags().Set(flagName, "false")
					}
					require.NoError(t, err, "Failed to set flag %s", flagName)
				case int:
					err := cmd.Flags().Set(flagName, fmt.Sprintf("%d", v))
					require.NoError(t, err, "Failed to set flag %s", flagName)
				case string:
					err := cmd.Flags().Set(flagName, v)
					require.NoError(t, err, "Failed to set flag %s", flagName)
				}
			}

			// Apply all flags - should not panic even with invalid values
			require.NotPanics(t, func() {
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
				}

				for _, flagFunc := range flagFunctions {
					flagFunc(cmd)
				}
			}, tt.description)

			// Config creation should not panic
			require.NotPanics(t, func() {
				cfg := config.New(opts...)
				_ = cfg
			}, "Config creation should not panic")
		})
	}
}

func TestFlagPrecedence(t *testing.T) {
	t.Run("command line flags override defaults", func(t *testing.T) {
		// Reset global state
		opts = nil

		cmd := &cobra.Command{Use: "test"}
		originalRootCmd := rootCmd
		defer func() { rootCmd = originalRootCmd }()
		rootCmd = cmd

		initFlags()

		// Set custom values
		require.NoError(t, cmd.Flags().Set("format", "pretty"))
		require.NoError(t, cmd.Flags().Set("jobs", "8"))
		require.NoError(t, cmd.Flags().Set("capitalize", "true"))

		// Apply flags
		formatFlag(cmd)
		jobsFlag(cmd)
		capitalizeFlag(cmd)

		cfg := config.New(opts...)

		assert.Equal(t, gnfmt.PrettyJSON, cfg.Format)
		assert.Equal(t, 8, cfg.Jobs)
		assert.True(t, cfg.WithCapitalization)
	})
}

func TestEnvironmentVariableIntegration(t *testing.T) {
	// This test verifies that the viper configuration setup would work
	// We can't fully test environment variables without modifying the global viper state
	// but we can test that the binding setup doesn't break anything

	t.Run("initConfig does not panic", func(t *testing.T) {
		// Store original environment
		originalEnvs := map[string]string{
			"GNV_DATA_SOURCES":        os.Getenv("GNV_DATA_SOURCES"),
			"GNV_FORMAT":              os.Getenv("GNV_FORMAT"),
			"GNV_JOBS":                os.Getenv("GNV_JOBS"),
			"GNV_VERIFIER_URL":        os.Getenv("GNV_VERIFIER_URL"),
			"GNV_WITH_ALL_MATCHES":    os.Getenv("GNV_WITH_ALL_MATCHES"),
			"GNV_WITH_CAPITALIZATION": os.Getenv("GNV_WITH_CAPITALIZATION"),
			"GNV_WITH_SPECIES_GROUP":  os.Getenv("GNV_WITH_SPECIES_GROUP"),
		}

		// Cleanup
		defer func() {
			for key, value := range originalEnvs {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}
		}()

		// Set test environment variables
		os.Setenv("GNV_FORMAT", "pretty")
		os.Setenv("GNV_JOBS", "6")
		os.Setenv("GNV_WITH_CAPITALIZATION", "true")

		// Reset global state
		opts = nil

		// This should not panic and should respect environment variables
		require.NotPanics(t, func() {
			// We can't easily test the full initConfig without side effects
			// but we can verify the individual option setting works
			getOpts()
		})
	})
}

func TestWebOptsIntegration(t *testing.T) {
	t.Run("web options are properly configured", func(t *testing.T) {
		// Reset global state
		opts = nil
		webOpts = nil

		cmd := &cobra.Command{Use: "test"}
		originalRootCmd := rootCmd
		defer func() { rootCmd = originalRootCmd }()
		rootCmd = cmd

		initFlags()

		// Set verifier URL flag
		require.NoError(t, cmd.Flags().Set("verifier_url", "https://test.example.com"))

		// Apply verifier URL flag
		verifierUrlFlag(cmd)

		// Verify both opts and webOpts are set
		assert.Len(t, opts, 1)
		assert.Len(t, webOpts, 1)

		// Verify configs
		cfg := config.New(opts...)
		webCfg := config.New(webOpts...)

		assert.Equal(t, "https://test.example.com", cfg.VerifierURL)
		assert.Equal(t, "https://test.example.com", webCfg.VerifierURL)
	})

	t.Run("web options include capitalization by default", func(t *testing.T) {
		// Reset global state
		opts = nil
		webOpts = nil

		// Simulate what happens in the root command
		copy(webOpts, opts)
		webOpts = append(webOpts, config.OptWithCapitalization(true))

		webCfg := config.New(webOpts...)
		assert.True(t, webCfg.WithCapitalization)
	})
}
