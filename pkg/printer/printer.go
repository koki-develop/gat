package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var (
	DefaultFormat = "terminal256"
	DefaultTheme  = "monokai"
)

type Printer struct {
	lang   string
	format string
	theme  string
}

type PrinterConfig struct {
	Lang   string
	Format string
	Theme  string
}

func New(cfg *PrinterConfig) *Printer {
	return &Printer{
		lang:   cfg.Lang,
		format: cfg.Format,
		theme:  cfg.Theme,
	}
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
	if p.lang == "" {
		if ipt.Filename != nil {
			l = lexers.Match(*ipt.Filename)
		}
		if l == nil {
			l = lexers.Analyse(src)
		}
		if l == nil {
			l = lexers.Fallback
		}
	} else {
		l = lexers.Get(p.lang)
		if l == nil {
			return fmt.Errorf("unknown lang: %s", p.lang)
		}
	}
	l = chroma.Coalesce(l)

	// get formatter
	f, ok := formatters.Registry[p.format]
	if !ok {
		return fmt.Errorf("unknown formatter: %s", p.format)
	}

	// get style
	s, ok := styles.Registry[p.theme]
	if !ok {
		return fmt.Errorf("unknown theme: %s", p.theme)
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
