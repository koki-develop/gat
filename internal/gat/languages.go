package gat

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/koki-develop/gat/internal/lexers"
)

func PrintLanguages(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	for _, l := range lexers.List() {
		cfg := l.Config()
		if _, err := tw.Write([]byte(fmt.Sprintf("%s\t%s\n", cfg.Name, strings.Join(cfg.Aliases, ", ")))); err != nil {
			return err
		}
	}

	return tw.Flush()
}
