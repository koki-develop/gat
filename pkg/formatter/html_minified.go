package formatter

import (
	"bytes"
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	htmlformatter "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

var (
	HTMLMinified = formatters.Register("html-min", NewHTMLMinifiedFormatter())
)

type HTMLMinifiedFormatter struct {
	htmlFormatter chroma.Formatter
	m             *minify.M
}

func NewHTMLMinifiedFormatter() *HTMLMinifiedFormatter {
	m := minify.New()
	m.Add(mimeTypeHTML, &html.Minifier{
		KeepDocumentTags: true,
		KeepQuotes:       true,
	})
	m.AddFunc(mimeTypeCSS, css.Minify)

	return &HTMLMinifiedFormatter{
		htmlFormatter: htmlformatter.New(htmlformatter.Standalone(true), htmlformatter.WithClasses(true)),
		m:             m,
	}
}

func (f *HTMLMinifiedFormatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	b := new(bytes.Buffer)
	if err := f.htmlFormatter.Format(b, style, iterator); err != nil {
		return err
	}

	if err := f.m.Minify(mimeTypeHTML, w, b); err != nil {
		return err
	}

	return nil
}
