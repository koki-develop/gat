package cmd

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/koki-develop/gat/pkg/printer"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	version string

	lang   string
	format string
	theme  string

	pretty bool

	listLangs   bool
	listFormats bool
	listThemes  bool

	forceColor bool
)

var rootCmd = &cobra.Command{
	Use:   "gat [file]...",
	Short: "cat alternative written in Go",
	Long:  "cat alternative written in Go.",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := printer.New(&printer.PrinterConfig{
			Lang:   lang,
			Format: format,
			Theme:  theme,
			Pretty: pretty,
		})

		if strings.HasPrefix(format, "terminal") {
			ist := term.IsTerminal(int(os.Stdout.Fd()))
			if !ist && !forceColor {
				p.SetTheme("noop")
			}
		}

		switch {
		case listLangs:
			printer.PrintLangs()
			return nil
		case listFormats:
			printer.PrintFormats()
			return nil
		case listThemes:
			printer.PrintThemes()
			return nil
		}

		if len(args) == 0 {
			if err := p.Print(&printer.PrintInput{
				In:  os.Stdin,
				Out: os.Stdout,
			}); err != nil {
				return err
			}
		} else {
			for _, filename := range args {
				if err := p.PrintFile(&printer.PrintFileInput{
					Out:      os.Stdout,
					Filename: filename,
				}); err != nil {
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
	rootCmd.Flags().StringVarP(&lang, "lang", "l", "", "language for syntax highlighting")
	rootCmd.Flags().StringVarP(&format, "format", "f", printer.DefaultFormat, "output format")
	rootCmd.Flags().StringVarP(&theme, "theme", "t", printer.DefaultTheme, "highlight theme")
	rootCmd.Flags().BoolVarP(&forceColor, "force-color", "c", false, "force colored output")

	rootCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "whether to format a content pretty")

	rootCmd.Flags().BoolVar(&listLangs, "list-langs", false, "print a list of supported languages for syntax highlighting")
	rootCmd.Flags().BoolVar(&listFormats, "list-formats", false, "print a list of supported output formats")
	rootCmd.Flags().BoolVar(&listThemes, "list-themes", false, "print a list of supported themes with preview")
	rootCmd.MarkFlagsMutuallyExclusive("list-langs", "list-formats", "list-themes")
}
