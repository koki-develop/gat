package cmd

import (
	"os"
	"strings"

	"github.com/koki-develop/gat/internal/gat"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// processFile handles opening, processing, and closing a single file with proper defer scope
func processFile(g *gat.Gat, filename string, opts ...gat.PrintOption) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	return g.Print(os.Stdout, f, opts...)
}

var rootCmd = &cobra.Command{
	Use:   "gat [file]...",
	Short: "cat alternative written in Go",
	Long:  "cat alternative written in Go.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ist := term.IsTerminal(int(os.Stdout.Fd()))

		switch {
		case flagListLangs:
			return gat.PrintLanguages(os.Stdout)
		case flagListFormats:
			return gat.PrintFormats(os.Stdout)
		case flagListThemes:
			return gat.PrintThemes(os.Stdout, ist)
		}

		if strings.HasPrefix(flagFormat, "terminal") {
			if !ist {
				if !flagForceColor {
					flagTheme = "noop"
				}
				flagForceBinary = true
				flagNoResize = true
			}
		}

		g, err := gat.New(&gat.Config{
			Language:       flagLang,
			Format:         flagFormat,
			Theme:          flagTheme,
			RenderMarkdown: flagRenderMarkdown,
			ForceBinary:    flagForceBinary,
			NoResize:       flagNoResize,
		})
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return g.Print(os.Stdout, os.Stdin, gat.WithPretty(flagPretty))
		}

		for _, filename := range args {
			if err := processFile(g, filename, gat.WithPretty(flagPretty), gat.WithFilename(filename)); err != nil {
				return err
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
