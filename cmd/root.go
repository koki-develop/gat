package cmd

import (
	"os"
	"runtime/debug"

	"github.com/koki-develop/gat/pkg/printer"
	"github.com/spf13/cobra"
)

var (
	version string

	lang   string
	format string
	theme  string
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
		})

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
}
