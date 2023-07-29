package formatters

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
)

func Get(name string) (chroma.Formatter, bool) {
	f, ok := formatters.Registry[name]
	return f, ok
}

func List() []string {
	return formatters.Names()
}
