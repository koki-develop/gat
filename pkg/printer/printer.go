package printer

import (
	"io"
	"os"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var (
	DefaultFormatter = "terminal256"
	DefaultStyle     = "monokai"
)

type Printer struct{}

func New() *Printer {
	return &Printer{}
}

type PrintFileInput struct {
	Out      io.Writer
	Filename string
}

func (p *Printer) PrintFile(ipt *PrintFileInput) error {
	f, err := os.Open(ipt.Filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := p.Print(&PrintInput{
		In:       f,
		Out:      ipt.Out,
		Filename: &ipt.Filename,
	}); err != nil {
		return err
	}

	return nil
}

type PrintInput struct {
	In       io.Reader
	Out      io.Writer
	Filename *string
}

func (p *Printer) Print(ipt *PrintInput) error {
	// read source
	b, err := io.ReadAll(ipt.In)
	if err != nil {
		return err
	}
	src := string(b)

	// get lexer
	var l chroma.Lexer
	if ipt.Filename != nil {
		l = lexers.Match(*ipt.Filename)
	}
	if l == nil {
		l = lexers.Analyse(src)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// get formatter
	f := formatters.Get(DefaultFormatter)
	if f == nil {
		f = formatters.Fallback
	}

	// get style
	s := styles.Get(DefaultStyle)
	if s == nil {
		s = styles.Fallback
	}

	// format
	it, err := l.Tokenise(nil, src)
	if err != nil {
		return err
	}
	if err := f.Format(ipt.Out, s, it); err != nil {
		return err
	}

	return nil
}
