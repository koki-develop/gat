package cmd

import (
	"os"

	"github.com/koki-develop/gat/pkg/printer"
	"github.com/spf13/cobra"
)

var (
	format string
	theme  string
)

var rootCmd = &cobra.Command{
	Use:   "gat [file]...",
	Short: "cat alternative written in Go",
	Long:  "cat alternative written in Go.",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := printer.New(&printer.PrinterConfig{
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
	rootCmd.Flags().StringVarP(&format, "format", "f", printer.DefaultFormat, "output format")
	rootCmd.Flags().StringVarP(&theme, "theme", "t", printer.DefaultTheme, "highlight theme")
}
