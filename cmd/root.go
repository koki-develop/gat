package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gat",
	Short: "cat alternative written in Go",
	Long:  "cat alternative written in Go.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hello World")
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
