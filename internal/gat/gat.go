package gat

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strings"
	"text/tabwriter"

	"github.com/alecthomas/chroma/v2"
	"github.com/koki-develop/gat/internal/formatters"
	"github.com/koki-develop/gat/internal/lexers"
	"github.com/koki-develop/gat/internal/prettier"
	"github.com/koki-develop/gat/internal/styles"
	"github.com/mattn/go-sixel"
)

type Config struct {
	Language    string
	Format      string
	Theme       string
	ForceBinary bool
}

type Gat struct {
	lexer       chroma.Lexer
	formatter   chroma.Formatter
	style       *chroma.Style
	forceBinary bool
}

func New(cfg *Config) (*Gat, error) {
	g := &Gat{
		forceBinary: cfg.ForceBinary,
	}

	// lexer
	if cfg.Language != "" {
		l, err := lexers.Get(lexers.WithLanguage(cfg.Language))
		if err != nil {
			return nil, err
		}
		g.lexer = l
	}

	// formatter
	f, ok := formatters.Get(cfg.Format)
	if !ok {
		return nil, fmt.Errorf("unknown format: %s", cfg.Format)
	}
	g.formatter = f

	// style
	s, ok := styles.Get(cfg.Theme)
	if !ok {
		return nil, fmt.Errorf("unknown theme: %s", cfg.Theme)
	}
	g.style = s

	return g, nil
}

type printOption struct {
	Pretty   bool
	Filename string
}

type PrintOption func(*printOption)

func WithPretty(p bool) PrintOption {
	return func(o *printOption) {
		o.Pretty = p
	}
}

func WithFilename(name string) PrintOption {
	return func(o *printOption) {
		o.Filename = name
	}
}

func (g *Gat) Print(w io.Writer, r io.Reader, opts ...PrintOption) error {
	// parse options
	opt := &printOption{}
	for _, o := range opts {
		o(opt)
	}

	// read w
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return err
	}

	// detect content type
	contentType := http.DetectContentType(buf.Bytes())

	// print image
	if strings.HasPrefix(contentType, "image/") && !g.forceBinary {
		if err := g.printImage(w, buf); err == nil {
			return nil
		}
	}

	// read source
	var src string
	switch contentType {
	case "application/x-gzip":
		s, err := g.readGzip(buf)
		if err != nil {
			return err
		}
		src = s
	default:
		isBin, err := isBinary(buf.Bytes())
		if err != nil {
			return err
		}
		if isBin {
			if g.forceBinary {
				if _, err := buf.WriteTo(w); err != nil {
					return err
				}
			} else {
				if _, err := w.Write([]byte("+----------------------------------------------------------------------------+\n| NOTE: This is a binary file. To force output, use the --force-binary flag. |\n+----------------------------------------------------------------------------+\n")); err != nil {
					return err
				}
			}
			return nil
		}

		src = buf.String()
	}

	// analyse lexer
	if g.lexer == nil {
		l, err := lexers.Get(lexers.WithFilename(opt.Filename), lexers.WithSource(src))
		if err != nil {
			return err
		}
		g.lexer = l
	}

	// pretty code
	if opt.Pretty {
		p, ok := prettier.Get(g.lexer.Config().Name)
		if ok {
			s, err := p.Pretty(src)
			if err != nil {
				return err
			}
			src = s
		}
	}

	// print
	it, err := g.lexer.Tokenise(nil, src)
	if err != nil {
		return err
	}
	if err := g.formatter.Format(w, g.style, it); err != nil {
		return err
	}

	return nil
}

func (*Gat) printImage(w io.Writer, r io.Reader) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	if err := sixel.NewEncoder(w).Encode(img); err != nil {
		return err
	}

	return nil
}

func (*Gat) readGzip(r io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	gz, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	if _, err := io.Copy(buf, gz); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func isBinary(data []byte) (bool, error) {
	if len(data) < 1024 {
		return bytes.IndexByte(data, 0) != -1, nil
	}
	return bytes.IndexByte(data[:1024], 0) != -1, nil
}

func PrintLanguages(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	if _, err := tw.Write([]byte("NAME\tALIASES\n")); err != nil {
		return err
	}

	for _, l := range lexers.List() {
		cfg := l.Config()
		if _, err := tw.Write([]byte(fmt.Sprintf("%s\t%s\n", cfg.Name, strings.Join(cfg.Aliases, ", ")))); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func PrintFormats(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	if _, err := tw.Write([]byte("NAME\n")); err != nil {
		return err
	}

	for _, f := range formatters.List() {
		if _, err := tw.Write([]byte(fmt.Sprintf("%s\n", f))); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func PrintThemes(w io.Writer) error {
	src := `package main

import "fmt"

func main() {
	fmt.Println("hello world")
}`

	for _, t := range styles.List() {
		fmt.Fprintf(w, "\x1b[1m%s\x1b[0m\n\n", t)

		g, err := New(&Config{
			Language: "go",
			Theme:    t,
			Format:   "terminal256",
		})
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		if err := g.Print(buf, strings.NewReader(src)); err != nil {
			return err
		}

		// indent source
		sc := bufio.NewScanner(buf)
		for sc.Scan() {
			if _, err := fmt.Fprintf(w, "\t%s\n", sc.Text()); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}

	return nil
}
