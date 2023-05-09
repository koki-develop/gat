package cmd

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/koki-develop/gat/internal/printer"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	version string

	flagLang   string
	flagFormat string
	flagTheme  string

	flagPretty bool

	flagListLangs   bool
	flagListFormats bool
	flagListThemes  bool

	flagForceColor bool
)

var rootCmd = &cobra.Command{
	Use:   "gat [file]...",
	Short: "cat alternative written in Go",
	Long:  "cat alternative written in Go.",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := printer.New(&printer.PrinterConfig{
			Lang:   flagLang,
			Format: flagFormat,
			Theme:  flagTheme,
			Pretty: flagPretty,
		})

		if strings.HasPrefix(flagFormat, "terminal") {
			ist := term.IsTerminal(int(os.Stdout.Fd()))
			if !ist && !flagForceColor {
				p.SetTheme("noop")
			}
		}

		switch {
		case flagListLangs:
			printer.PrintLangs()
			return nil
		case flagListFormats:
			printer.PrintFormats()
			return nil
		case flagListThemes:
			printer.PrintThemes()
			return nil
		}

		if len(args) == 0 {
			if err := p.Print(os.Stdin, os.Stdout); err != nil {
				return err
			}
		} else {
			for _, filename := range args {
				f, err := os.Open(filename)
				if err != nil {
					return err
				}
				defer f.Close()
				if err := p.Print(f, os.Stdout, printer.WithFilename(filename)); err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// version
	if version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			version = info.Main.Version
		}
	}

	rootCmd.Version = version

	// flags
	rootCmd.Flags().StringVarP(&flagLang, "lang", "l", "", "language for syntax highlighting")
	rootCmd.Flags().StringVarP(&flagFormat, "format", "f", printer.DefaultFormat, "output format")
	rootCmd.Flags().StringVarP(&flagTheme, "theme", "t", printer.DefaultTheme, "highlight theme")
	rootCmd.Flags().BoolVarP(&flagForceColor, "force-color", "c", false, "force colored output")

	rootCmd.Flags().BoolVarP(&flagPretty, "pretty", "p", false, "whether to format a content pretty")

	rootCmd.Flags().BoolVar(&flagListLangs, "list-langs", false, "print a list of supported languages for syntax highlighting")
	rootCmd.Flags().BoolVar(&flagListFormats, "list-formats", false, "print a list of supported output formats")
	rootCmd.Flags().BoolVar(&flagListThemes, "list-themes", false, "print a list of supported themes with preview")
	rootCmd.MarkFlagsMutuallyExclusive("list-langs", "list-formats", "list-themes")
}
