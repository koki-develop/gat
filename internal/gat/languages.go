package gat

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/koki-develop/gat/internal/lexers"
)

type Language struct {
	Name    string
	Aliases []string
}

func PrintLanguages(w io.Writer) error {
	return printLanguages(w, listLanguages())
}

func printLanguages(w io.Writer, langs []Language) error {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	for _, l := range langs {
		if _, err := tw.Write([]byte(fmt.Sprintf("%s\t%s\n", l.Name, strings.Join(l.Aliases, ", ")))); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func listLanguages() []Language {
	ls := lexers.List()

	rtn := make([]Language, len(ls))
	for i, l := range ls {
		cfg := l.Config()
		rtn[i] = Language{Name: cfg.Name, Aliases: cfg.Aliases}
	}

	return rtn
}
