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

	"github.com/alecthomas/chroma/v2"
	"github.com/charmbracelet/glamour"
	"github.com/koki-develop/gat/internal/display"
	"github.com/koki-develop/gat/internal/formatters"
	"github.com/koki-develop/gat/internal/lexers"
	"github.com/koki-develop/gat/internal/prettier"
	"github.com/koki-develop/gat/internal/styles"
	"github.com/koki-develop/mask-go"
	"github.com/mattn/go-sixel"
	"golang.org/x/image/draw"
)

// masker is what --mask-secrets masks with: every pattern mask-go carries, so
// that a pattern added there is one gat masks with on the next upgrade.
var masker = mask.New(mask.WithPatterns(mask.AllBuiltinPatterns()...))

// markdownMasker is what the rendered markdown branch masks with. It redacts to
// a bare word because that branch masks the source glamour then reads, where a
// redaction is markup: asterisks alone on a line are a thematic break, so the
// value renders away to a horizontal rule and inside a list takes the rest of
// the list with it, and brackets are link syntax, so a redaction followed by a
// colon is a link reference definition and renders as nothing at all.
var markdownMasker = mask.New(
	mask.WithPatterns(mask.AllBuiltinPatterns()...),
	mask.WithRedactor(mask.Fixed("REDACTED")),
)

// maskMaxRetained is how much text the streaming path holds back while a
// pattern is still reading a value, before it gives up and masks what it holds.
//
// Giving up masks everything held and everything after it to the end of the
// stream, so on a file carrying one unbroken run — an unterminated PEM block, a
// megabytes-long base64 blob — the default mebibyte is low enough to turn the
// rest of the output into asterisks. The highlighting path buffers whole files
// already, so the limit is raised to where a real file does not reach it.
//
// A stream that never ends pays for that: a value it leaves open holds the
// output until the limit, so raising the limit lengthens the silence. A file
// pays nothing, its end settling every pattern still reading.
const maskMaxRetained = 16 << 20

type Config struct {
	Language       string
	Format         string
	Theme          string
	RenderMarkdown bool
	ForceBinary    bool
	NoResize       bool
}

type Gat struct {
	explicitLexer  chroma.Lexer
	formatter      chroma.Formatter
	style          *chroma.Style
	renderMarkdown bool
	forceBinary    bool
	noResize       bool
	noColor        bool
	terminalFormat bool
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
		g.explicitLexer = l
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

	g.noColor = cfg.Theme == "noop"
	g.terminalFormat = strings.HasPrefix(cfg.Format, "terminal")

	return g, nil
}

type printOption struct {
	Pretty   bool
	Mask     bool
	Filename string
	Display  *display.Options
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

func WithMask(m bool) PrintOption {
	return func(o *printOption) {
		o.Mask = m
	}
}

func WithDisplay(d *display.Options) PrintOption {
	return func(o *printOption) {
		o.Display = d
	}
}

// maskedReader wraps r so that a path streaming its input straight out masks
// on the way. The reader holds back only the tail its patterns are still
// reading, so a value written across two reads is masked whole and ordinary
// text goes straight through.
//
// A newline does not settle text; the next write does. So a feed being followed
// stalls a line behind whenever its newest line ends on something a pattern has
// opened — a bare aws_secret_access_key=, a PEM begin boundary, a trailing sk-.
// That is the price of masking across chunks. Masking a line at a time would
// not stall, and would release the body of a private key one line at a time.
func (g *Gat) maskedReader(r io.Reader, opt *printOption) io.Reader {
	if !opt.Mask {
		return r
	}
	return mask.NewReader(r, masker, mask.WithMaxRetained(maskMaxRetained))
}

// isPassthrough reports whether the configuration requires no tokenization, so
// the input can be streamed out as it is read instead of being buffered in
// full. Display transforms still apply per chunk. Masking applies across
// chunks: a value may begin in one read and end in another, so the text a
// pattern is still reading is held back rather than written out in pieces.
func (g *Gat) isPassthrough(opt *printOption) bool {
	return g.noColor && g.terminalFormat && !g.renderMarkdown && !opt.Pretty
}

// detectionHead returns the leading bytes used for content-type and binary
// detection. Normally it reads up to 1024 bytes. In passthrough mode it blocks
// only for the first read and returns whatever is already buffered (capped at
// 1024), so a slowly trickling text stream is not held back waiting for a full
// 1024-byte read. For bulk input the first read already fills the buffer, so
// detection is identical to the non-passthrough path, and image/gzip/binary
// content is routed to its normal handler by the caller.
//
// Trade-off: when non-text content (gzip, binary) arrives in a first read
// smaller than its detection window, it goes unrecognized and is streamed
// through as text. This never happens for files or bulk streams (the first
// read is large); it only affects content trickling in tiny fragments, which
// in practice is text. The common non-tty passthrough path also sets
// forceBinary, which already streams binary raw regardless.
func (g *Gat) detectionHead(br *bufio.Reader, opt *printOption) ([]byte, error) {
	if !g.isPassthrough(opt) {
		return br.Peek(1024)
	}
	if _, err := br.Peek(1); err != nil {
		return nil, err // includes io.EOF for empty input
	}
	n := br.Buffered()
	if n > 1024 {
		n = 1024
	}
	return br.Peek(n)
}

func (g *Gat) Print(w io.Writer, r io.Reader, opts ...PrintOption) error {
	// parse options
	opt := &printOption{}
	for _, o := range opts {
		o(opt)
	}

	br := bufio.NewReader(r)
	head, err := g.detectionHead(br, opt)
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
				if _, err := io.Copy(w, g.maskedReader(br, opt)); err != nil {
					return err
				}
			} else {
				if _, err := w.Write([]byte("+----------------------------------------------------------------------------+\n| NOTE: This is a binary file. To force output, use the --force-binary flag. |\n+----------------------------------------------------------------------------+\n")); err != nil {
					return err
				}
			}
			return nil
		}

		// Stream the input out as it is read instead of buffering the whole
		// source (only valid when isPassthrough; see its doc).
		if g.isPassthrough(opt) {
			out := w
			if opt.Display != nil {
				out = display.NewWriter(w, opt.Display)
			}

			_, err := io.Copy(out, g.maskedReader(br, opt))
			return err
		}

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, br); err != nil {
			return err
		}
		src = buf.String()
	}

	// analyse lexer
	lexer := g.explicitLexer
	if lexer == nil {
		l, err := lexers.Get(lexers.WithFilename(opt.Filename), lexers.WithSource(src))
		if err != nil {
			return err
		}
		lexer = l
	}

	if g.renderMarkdown && lexer.Config().Name == "markdown" {
		// Masking happens here rather than at the shared call below, which this
		// branch returns before reaching, and it masks the source rather than
		// the rendered output: rendering reflows text, and a value it wrapped
		// across two lines is one no pattern would find.
		if opt.Mask {
			src = markdownMasker.Mask(src)
		}

		r, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(-1),
		)
		if err != nil {
			return err
		}
		defer func() { _ = r.Close() }()

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
		p, ok := prettier.Get(lexer.Config().Name)
		if ok {
			s, err := p.Pretty(src)
			if err == nil {
				src = s
			}
		}
	}

	// mask sensitive information
	if opt.Mask {
		src = masker.Mask(src)
	}

	// display transformation
	out := w
	if opt.Display != nil {
		out = display.NewWriter(w, opt.Display)
	}

	// print
	it, err := lexer.Tokenise(nil, src)
	if err != nil {
		return err
	}
	if err := g.formatter.Format(out, g.style, it); err != nil {
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
	defer func() { _ = gz.Close() }()

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

func PrintThemes(w io.Writer, withColor bool) error {
	if withColor {
		src := `package main

import "fmt"

func main() {
	fmt.Println("hello world")
}`

		for _, t := range styles.List() {
			if _, err := fmt.Fprintf(w, "\x1b[1m%s\x1b[0m\n\n", t); err != nil {
				return err
			}

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
			if _, err := fmt.Fprintln(w, t); err != nil {
				return err
			}
		}
	}

	return nil
}
