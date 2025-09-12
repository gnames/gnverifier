# Flag Testing Documentation

This document describes the comprehensive test suite for all CLI flags in the gnverifier application.

## Test Files Overview

The flag testing suite consists of three main test files:

### 1. `flags_test.go` - Core Flag Function Tests
Tests individual flag processing functions and parsing utilities.

**Coverage:**
- `TestVersionFlag` - Tests version flag behavior
- `TestQuietFlag` - Tests quiet logging flag
- `TestCapitalizeFlag` - Tests capitalization flag
- `TestSpGroupFlag` - Tests species group flag
- `TestFuzzyRelaxedFlag` - Tests relaxed fuzzy matching flag
- `TestFuzzyUninomialFlag` - Tests uninomial fuzzy matching flag
- `TestFormatFlag` - Tests output format flag with all valid formats
- `TestJobsFlag` - Tests parallel jobs flag with boundary conditions
- `TestAllMatchesFlag` - Tests all matches flag
- `TestSourcesFlag` - Tests data sources flag with validation
- `TestVerifierUrlFlag` - Tests custom verifier URL flag
- `TestVernacularsFlag` - Tests vernacular languages flag
- `TestParseDataSources` - Tests data source parsing logic
- `TestParseVernacularLanguages` - Tests language parsing logic
- `TestAllFlagFunctions` - Integration test ensuring all flag functions work

### 2. `flags_init_test.go` - Flag Initialization Tests
Tests the flag initialization system and ensures all flags are properly registered.

**Coverage:**
- `TestInitFlags` - Verifies all 14 expected flags are created
- `TestBaseFlags` - Tests version flag registration
- `TestWebFlags` - Tests port flag registration
- `TestVerificationFlags` - Tests all verification-related flags
- `TestFormatFlags` - Tests formatting-related flags
- `TestPerformanceFlags` - Tests performance flags (jobs)
- `TestDataSourcesFlags` - Tests data source flags
- `TestFlagTypes` - Ensures correct flag types (bool, int, string)
- `TestFlagShorthands` - Verifies all shorthand mappings are unique and correct
- `TestFlagDefaults` - Tests default values for all flags

### 3. `flags_integration_test.go` - Integration and Edge Case Tests
Tests flag combinations, conflicts, and advanced scenarios.

**Coverage:**
- `TestFlagIntegration` - Tests complex flag combinations
- `TestFlagConflicts` - Tests invalid inputs and error handling
- `TestFlagPrecedence` - Tests command line flag precedence
- `TestEnvironmentVariableIntegration` - Tests environment variable support
- `TestWebOptsIntegration` - Tests web-specific option handling

## Flag Coverage

The test suite covers all 14 CLI flags:

### Base Flags
- `--version, -V` - Version information flag

### Web Flags  
- `--port, -p` - Web GUI port flag

### Verification Flags
- `--verifier_url, -v` - Custom verifier URL
- `--name_field, -n` - Scientific name field position
- `--all_matches, -M` - Return all matches flag
- `--species_group, -G` - Species group search flag
- `--fuzzy_relaxed, -R` - Relaxed fuzzy matching flag
- `--fuzzy_uninomial, -U` - Uninomial fuzzy matching flag
- `--vernaculars, -r` - Vernacular language search flag

### Format Flags
- `--quiet, -q` - Quiet progress flag
- `--capitalize, -c` - Capitalization flag
- `--format, -f` - Output format flag

### Performance Flags
- `--jobs, -j` - Parallel jobs flag

### Data Source Flags
- `--sources, -s` - Data source IDs flag

## Testing Strategy

### 1. Unit Testing
Each flag function is tested in isolation with:
- Valid inputs
- Invalid inputs
- Edge cases
- Default behavior

### 2. Parsing Testing
Helper functions `parseDataSources()` and `parseVernacularLanguages()` are tested with:
- Empty inputs
- Single values
- Multiple values
- Invalid formats
- Mixed valid/invalid inputs

### 3. Integration Testing
Tests combinations of flags to ensure:
- No conflicts between flags
- Proper configuration generation
- Correct option precedence
- Error handling with invalid combinations

### 4. Initialization Testing
Ensures the flag system is properly set up:
- All expected flags exist
- Correct types and defaults
- Proper shorthand mappings
- No duplicate flags

## Key Testing Features

### Global State Management
Tests properly reset global variables (`opts`, `webOpts`) between test runs to ensure isolation.

### Error Handling
Tests verify that invalid inputs are handled gracefully:
- Invalid formats default to CSV
- Invalid data sources are filtered out
- Invalid vernacular languages generate warnings

### Configuration Integration
Tests verify that flag options correctly translate to configuration settings by creating `config.Config` instances and checking their properties.

### Cobra Integration
Tests work with the Cobra CLI framework by:
- Creating test commands
- Setting flag values programmatically
- Verifying flag registration and lookup

## Running the Tests

```bash
# Run all flag tests
go test ./cmd/... -v

# Run specific test categories
go test ./cmd -run "TestFlag" -v           # Individual flag tests
go test ./cmd -run "TestInit" -v           # Initialization tests  
go test ./cmd -run "TestIntegration" -v    # Integration tests

# Run tests with coverage
go test ./cmd/... -cover
```

## Test Patterns

### Flag Function Test Pattern
```go
func TestSomeFlag(t *testing.T) {
    tests := []struct {
        name      string
        flagValue interface{}
        expectOpt bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Reset global state
            opts = nil
            webOpts = nil
            
            // Create command and set flag
            cmd := &cobra.Command{}
            cmd.Flags().Bool("flag_name", tt.flagValue, "test flag")
            
            // Execute flag function
            flagFunction(cmd)
            
            // Verify results
            if tt.expectOpt {
                assert.Len(t, opts, 1)
                cfg := config.New(opts...)
                // Assert configuration properties
            } else {
                assert.Len(t, opts, 0)
            }
        })
    }
}
```

This comprehensive test suite ensures that all CLI flags work correctly in isolation and in combination, providing confidence in the application's command-line interface.