package formatter

import (
	"bytes"
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"
)

var (
	JSONMinified = formatters.Register("json-min", NewJSONMinifiedFormatter())
)

type JSONMinifiedFormatter struct {
	jsonFormatter chroma.Formatter
	m             *minify.M
}

func NewJSONMinifiedFormatter() *JSONMinifiedFormatter {
	m := minify.New()
	m.AddFunc(mimeTypeJSON, json.Minify)

	return &JSONMinifiedFormatter{
		jsonFormatter: formatters.JSON,
		m:             m,
	}
}

func (f *JSONMinifiedFormatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	b := new(bytes.Buffer)
	if err := f.jsonFormatter.Format(b, style, iterator); err != nil {
		return err
	}

	if err := f.m.Minify(mimeTypeJSON, w, b); err != nil {
		return err
	}

	return nil
}
