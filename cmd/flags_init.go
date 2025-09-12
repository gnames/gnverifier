package cmd

func initFlags() {
	baseFlags()
	webFlags()
	verificationFlags()
	formatFlags()
	performanceFlags()
	dataSourcesFlags()
}

func baseFlags() {
	rootCmd.Flags().BoolP("version", "V", false, "Prints version information")
}

func webFlags() {
	rootCmd.Flags().IntP("port", "p", 0, "Port to run web GUI.")
}

func verificationFlags() {
	rootCmd.Flags().StringP("verifier_url", "v", "",
		`URL for verification service.
  Default: https://verifier.globalnames.org/api/v1`)
	rootCmd.Flags().IntP("name_field", "n", 1, "Set position of ScientificName field, the first field is 1.")
	rootCmd.Flags().BoolP("all_matches", "M", false, "return all matched results per source, not just the best one.")
	rootCmd.Flags().BoolP("species_group", "G", false, "searching for species names also searches their species groups.")
	rootCmd.Flags().BoolP("fuzzy_relaxed", "R", false,
		"relaxes fuzzy matching rules, decreses max names to 50.")
	rootCmd.Flags().BoolP("fuzzy_uninomial", "U", false,
		"allows fuzzy matching for uninomial names.")
	rootCmd.Flags().StringP(
		"vernaculars", "r", "",
		`sets languages for vernacular names search (e.g., "eng,deu,rus")
limited to 50 scientific names, try it with iNaturalist (id 180)`,
	)
}

func formatFlags() {
	rootCmd.Flags().BoolP("quiet", "q", false, "do not show progress")
	rootCmd.Flags().BoolP("capitalize", "c", false, "capitalizes first character")
	rootCmd.Flags().StringP("format", "f", "csv", `Format of the output: "compact", "pretty", "csv", "tsv".
  compact: compact JSON,
  pretty: pretty JSON,
  csv: CSV (DEFAULT)`)
}

func performanceFlags() {
	rootCmd.Flags().IntP("jobs", "j", 4, "Number of jobs running in parallel.")
}

func dataSourcesFlags() {
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
  180 - iNaturalist
  181 - IRMNG
  194 - PLAZI
  195 - AlgaeBase`)
}
