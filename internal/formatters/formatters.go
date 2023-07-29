package formatters

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
)

func Get(name string) (chroma.Formatter, bool) {
	l, ok := formatters.Registry[name]
	return l, ok
}
