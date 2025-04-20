package gat

import (
	"fmt"
	"io"

	"github.com/koki-develop/gat/internal/formatters"
)

func PrintFormats(w io.Writer) error {
	return printFormats(w, listFormats())
}

func printFormats(w io.Writer, formats []string) error {
	for _, f := range formats {
		if _, err := fmt.Fprintln(w, f); err != nil {
			return err
		}
	}
	return nil
}

func listFormats() []string {
	return formatters.List()
}
