package cmd

var (
	flagLang   string
	flagFormat string
	flagTheme  string

	flagPretty bool

	flagListLangs   bool
	flagListFormats bool
	flagListThemes  bool

	flagForceColor bool
)

func init() {
	rootCmd.Flags().StringVarP(&flagLang, "lang", "l", "", "language for syntax highlighting")
	rootCmd.Flags().StringVarP(&flagFormat, "format", "f", "terminal256", "output format")
	rootCmd.Flags().StringVarP(&flagTheme, "theme", "t", "monokai", "highlight theme")
	rootCmd.Flags().BoolVarP(&flagForceColor, "force-color", "c", false, "force colored output")

	rootCmd.Flags().BoolVarP(&flagPretty, "pretty", "p", false, "whether to format a content pretty")

	rootCmd.Flags().BoolVar(&flagListLangs, "list-langs", false, "print a list of supported languages for syntax highlighting")
	rootCmd.Flags().BoolVar(&flagListFormats, "list-formats", false, "print a list of supported output formats")
	rootCmd.Flags().BoolVar(&flagListThemes, "list-themes", false, "print a list of supported themes with preview")
	rootCmd.MarkFlagsMutuallyExclusive("list-langs", "list-formats", "list-themes")
}
