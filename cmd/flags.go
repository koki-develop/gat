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

	// -M, --render-markdown
	flagRenderMarkdown bool

	// --force-color
	flagForceColor bool

	// --force-binary
	flagForceBinary bool

	// --no-resize
	flagNoResize bool

	// --pretty
	flagPretty bool

	// --mask-secrets
	flagMaskSecrets bool

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
	rootCmd.Flags().BoolVarP(&flagRenderMarkdown, "render-markdown", "M", false, "render markdown")
	rootCmd.Flags().BoolVarP(&flagForceColor, "force-color", "c", false, "force colored output")
	rootCmd.Flags().BoolVarP(&flagForceBinary, "force-binary", "b", false, "force binary output")
	rootCmd.Flags().BoolVar(&flagNoResize, "no-resize", false, "do not resize images")

	rootCmd.Flags().BoolVarP(&flagPretty, "pretty", "p", false, "whether to format a content pretty")
	rootCmd.Flags().BoolVar(&flagMaskSecrets, "mask-secrets", false, "mask sensitive information (API keys, tokens)")

	rootCmd.Flags().BoolVar(&flagListLangs, "list-langs", false, "print a list of supported languages for syntax highlighting")
	rootCmd.Flags().BoolVar(&flagListFormats, "list-formats", false, "print a list of supported output formats")
	rootCmd.Flags().BoolVar(&flagListThemes, "list-themes", false, "print a list of supported themes with preview")
	rootCmd.MarkFlagsMutuallyExclusive("list-langs", "list-formats", "list-themes")
}
