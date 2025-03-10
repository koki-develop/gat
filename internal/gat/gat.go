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
	"github.com/charmbracelet/glamour"
	"github.com/koki-develop/gat/internal/formatters"
	"github.com/koki-develop/gat/internal/lexers"
	"github.com/koki-develop/gat/internal/prettier"
	"github.com/koki-develop/gat/internal/styles"
	"github.com/mattn/go-sixel"
	"golang.org/x/image/draw"
)

type Config struct {
	Language       string
	Format         string
	Theme          string
	RenderMarkdown bool
	ForceBinary    bool
	NoResize       bool
}

type Gat struct {
	lexer          chroma.Lexer
	formatter      chroma.Formatter
	style          *chroma.Style
	renderMarkdown bool
	forceBinary    bool
	noResize       bool
}

func New(cfg *Config) (*Gat, error) {
	g := &Gat{
		renderMarkdown: cfg.RenderMarkdown,
		forceBinary:    cfg.ForceBinary,
		noResize:       cfg.NoResize,
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

	br := bufio.NewReader(r)
	head, err := br.Peek(1024)
	if err != nil && err != io.EOF {
		return err
	}

	// detect content type
	contentType := http.DetectContentType(head)

	// print image
	if strings.HasPrefix(contentType, "image/") && !g.forceBinary {
		if err := g.printImage(w, br); err == nil {
			return nil
		}
	}

	// read source
	var src string
	switch contentType {
	case "application/x-gzip":
		s, err := g.readGzip(br)
		if err != nil {
			return err
		}
		src = s
	default:
		if isBinary(head) {
			if g.forceBinary {
				if _, err := br.WriteTo(w); err != nil {
					return err
				}
			} else {
				if _, err := w.Write([]byte("+----------------------------------------------------------------------------+\n| NOTE: This is a binary file. To force output, use the --force-binary flag. |\n+----------------------------------------------------------------------------+\n")); err != nil {
					return err
				}
			}
			return nil
		}

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, br); err != nil {
			return err
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

	if g.renderMarkdown && g.lexer.Config().Name == "markdown" {
		r, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(-1),
		)
		if err != nil {
			return err
		}
		defer r.Close()

		s, err := r.Render(src)
		if err != nil {
			return err
		}
		if _, err := w.Write([]byte(s)); err != nil {
			return err
		}
		return nil
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

func (g *Gat) printImage(w io.Writer, r io.Reader) error {
	maxEdge := 1800

	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()

	if g.noResize || (imgWidth <= maxEdge && imgHeight <= maxEdge) {
		if err := sixel.NewEncoder(w).Encode(img); err != nil {
			return err
		}
	} else {
		var dstWidth, dstHeight int
		aspectRatio := float64(imgHeight) / float64(imgWidth)
		if imgWidth > imgHeight {
			dstWidth, dstHeight = maxEdge, int(float64(maxEdge)*aspectRatio)
		} else {
			dstWidth, dstHeight = int(float64(maxEdge)/aspectRatio), maxEdge
		}

		dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))
		draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Src, nil)
		if err := sixel.NewEncoder(w).Encode(dst); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte{'\n'}); err != nil {
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

func isBinary(data []byte) bool {
	if len(data) < 1024 {
		return bytes.IndexByte(data, 0) != -1
	}
	return bytes.IndexByte(data[:1024], 0) != -1
}

func PrintLanguages(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	for _, l := range lexers.List() {
		cfg := l.Config()
		if _, err := tw.Write([]byte(fmt.Sprintf("%s\t%s\n", cfg.Name, strings.Join(cfg.Aliases, ", ")))); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func PrintFormats(w io.Writer) {
	for _, f := range formatters.List() {
		fmt.Fprintln(w, f)
	}
}

func PrintThemes(w io.Writer, withColor bool) error {
	if withColor {
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
	} else {
		for _, t := range styles.List() {
			fmt.Fprintln(w, t)
		}
	}

	return nil
}
