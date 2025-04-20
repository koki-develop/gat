package gat

import (
	"fmt"
	"io"

	"github.com/koki-develop/gat/internal/formatters"
)

func PrintFormats(w io.Writer) {
	for _, f := range formatters.List() {
		fmt.Fprintln(w, f)
	}
}
