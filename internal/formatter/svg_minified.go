package formatter

import (
	"bytes"
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/svg"
)

var (
	SVGMinified = formatters.Register("svg-min", NewSVGMinifiedFormatter())
)

type SVGMinifiedFormatter struct {
	svgFormatter chroma.Formatter
	m            *minify.M
}

func NewSVGMinifiedFormatter() *SVGMinifiedFormatter {
	m := minify.New()
	m.AddFunc(mimeTypeSVG, svg.Minify)

	return &SVGMinifiedFormatter{
		svgFormatter: formatters.SVG,
		m:            m,
	}
}

func (f *SVGMinifiedFormatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	b := new(bytes.Buffer)
	if err := f.svgFormatter.Format(b, style, iterator); err != nil {
		return err
	}

	if err := f.m.Minify(mimeTypeSVG, w, b); err != nil {
		return err
	}

	return nil
}
