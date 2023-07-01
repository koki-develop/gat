package printer

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	_ "github.com/koki-develop/gat/internal/formatter"
	"github.com/koki-develop/gat/internal/prettier"
	_ "github.com/koki-develop/gat/internal/style"
	"github.com/mattn/go-sixel"
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

func (p *Printer) Print(in io.Reader, out io.Writer, opts ...Option) error {
	opt := &option{}
	for _, o := range opts {
		o(opt)
	}

	// read source
	b := new(bytes.Buffer)
	if _, err := io.Copy(b, in); err != nil {
		return err
	}

	// print image
	contentType := http.DetectContentType(b.Bytes())
	if strings.HasPrefix(contentType, "image/") {
		img, _, err := image.Decode(b)
		if err == nil {
			if err := sixel.NewEncoder(out).Encode(img); err != nil {
				return err
			}
			return nil
		}
	}

	src := b.String()

	// get lexer
	var l chroma.Lexer
	if p.lang == "" {
		if opt.filename != nil {
			l = lexers.Match(*opt.filename)
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
	if err := f.Format(out, s, it); err != nil {
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
		if err := p.Print(strings.NewReader(example), b); err != nil {
			panic(err)
		}

		sc := bufio.NewScanner(strings.NewReader(b.String()))
		for sc.Scan() {
			fmt.Printf("\t%s\n", sc.Text())
		}

		fmt.Print("\n\n")
	}
}
