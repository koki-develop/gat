package printer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
        "strings"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	_ "github.com/koki-develop/gat/pkg/formatter"
	"github.com/koki-develop/gat/pkg/prettier"
	_ "github.com/koki-develop/gat/pkg/style"
)

var (
	DefaultFormat = "terminal256"
	DefaultTheme  = "monokai"
)

type Printer struct {
	lang   string
	format string
	theme  string

	pretty bool
}

type PrinterConfig struct {
	Lang   string
	Format string
	Theme  string

	Pretty bool
}

func New(cfg *PrinterConfig) *Printer {
	return &Printer{
		lang:   cfg.Lang,
		format: cfg.Format,
		theme:  cfg.Theme,
		pretty: cfg.Pretty,
	}
}

func (p *Printer) SetTheme(t string) {
	p.theme = t
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
	// b, err := io.ReadAll(ipt.In)
	strB := new(strings.Builder)
        _,err := io.Copy(strB, ipt.In)
        if err != nil {
                return err
        }
        src := strB.String()
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

	// pretty
	if p.pretty {
		pt := prettier.Get(l.Config().Name)
		if prettied, err := pt.Pretty(src); err == nil {
			src = prettied
		}
	}

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

func PrintLangs() {
	for _, l := range lexers.GlobalLexerRegistry.Lexers {
		cfg := l.Config()
		fmt.Print(cfg.Name)
		if len(cfg.Aliases) > 0 {
			fmt.Printf(" (%s)", strings.Join(cfg.Aliases, ", "))
		}
		fmt.Print("\n")
	}
}

func PrintFormats() {
	for _, f := range formatters.Names() {
		fmt.Println(f)
	}
}

var (
	example = `package main

import "fmt"

func main() {
	fmt.Println("hello world")
}`
)

func PrintThemes() {
	for _, t := range styles.Names() {
		fmt.Printf("\x1b[1m%s\x1b[0m\n\n", t)

		b := new(bytes.Buffer)
		p := New(&PrinterConfig{
			Lang:   "go",
			Format: DefaultFormat,
			Theme:  t,
		})
		if err := p.Print(&PrintInput{
			In:  strings.NewReader(example),
			Out: b,
		}); err != nil {
			panic(err)
		}

		sc := bufio.NewScanner(strings.NewReader(b.String()))
		for sc.Scan() {
			fmt.Printf("\t%s\n", sc.Text())
		}

		fmt.Print("\n\n")
	}
}
