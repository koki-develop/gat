package cmd

import "os"

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

var (
	// --lang
	flagLang string

	// --format
	flagFormat        string
	flagFormatDefault = envOrDefault("GAT_FORMAT", "terminal256")

	// --theme
	flagTheme        string
	flagThemeDefault = envOrDefault("GAT_THEME", "monokai")

	// --force-color
	flagForceColor bool

	// --pretty
	flagPretty bool

	// --list-langs
	flagListLangs bool

	// --list-formats
	flagListFormats bool

	// --list-themes
	flagListThemes bool
)

func init() {
	rootCmd.Flags().StringVarP(&flagLang, "lang", "l", "", "language for syntax highlighting")
	rootCmd.Flags().StringVarP(&flagFormat, "format", "f", flagFormatDefault, "output format")
	rootCmd.Flags().StringVarP(&flagTheme, "theme", "t", flagThemeDefault, "highlight theme")
	rootCmd.Flags().BoolVarP(&flagForceColor, "force-color", "c", false, "force colored output")

	rootCmd.Flags().BoolVarP(&flagPretty, "pretty", "p", false, "whether to format a content pretty")

	rootCmd.Flags().BoolVar(&flagListLangs, "list-langs", false, "print a list of supported languages for syntax highlighting")
	rootCmd.Flags().BoolVar(&flagListFormats, "list-formats", false, "print a list of supported output formats")
	rootCmd.Flags().BoolVar(&flagListThemes, "list-themes", false, "print a list of supported themes with preview")
	rootCmd.MarkFlagsMutuallyExclusive("list-langs", "list-formats", "list-themes")
}
