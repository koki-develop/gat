package gat

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/koki-develop/gat/internal/display"
	"github.com/koki-develop/gat/internal/masker"
)

func TestGat_isPassthrough(t *testing.T) {
	tests := []struct {
		name           string
		theme          string
		format         string
		renderMarkdown bool
		pretty         bool
		want           bool
	}{
		{"noop + terminal256, no pretty/md", "noop", "terminal256", false, false, true},
		{"noop + terminal (plain)", "noop", "terminal", false, false, true},
		{"non-noop theme", "monokai", "terminal256", false, false, false},
		{"non-terminal format (html)", "noop", "html", false, false, false},
		{"pretty enabled", "noop", "terminal256", false, true, false},
		{"render-markdown enabled", "noop", "terminal256", true, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := New(&Config{
				Theme:          tt.theme,
				Format:         tt.format,
				RenderMarkdown: tt.renderMarkdown,
			})
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			got := g.isPassthrough(&printOption{Pretty: tt.pretty})
			if got != tt.want {
				t.Errorf("isPassthrough() = %v, want %v", got, tt.want)
			}
		})
	}
}

// newPassthroughGat builds a Gat configured for the passthrough path.
func newPassthroughGat(t *testing.T) *Gat {
	t.Helper()
	g, err := New(&Config{Theme: "noop", Format: "terminal256"})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return g
}

func TestGat_Print_passthrough(t *testing.T) {
	const ghToken = "ghp_0123456789abcdefghijklmnopqrstuvwxyz123456"

	t.Run("mask disabled: output equals input", func(t *testing.T) {
		in := "line1\nline2 " + ghToken + "\nline3\n"
		var buf bytes.Buffer
		if err := newPassthroughGat(t).Print(&buf, strings.NewReader(in), WithMask(false)); err != nil {
			t.Fatalf("Print() error = %v", err)
		}
		if buf.String() != in {
			t.Errorf("output = %q, want %q", buf.String(), in)
		}
	})

	t.Run("mask enabled: equals full-text Mask", func(t *testing.T) {
		in := "line1\nline2 " + ghToken + "\nline3\n"
		var buf bytes.Buffer
		if err := newPassthroughGat(t).Print(&buf, strings.NewReader(in), WithMask(true)); err != nil {
			t.Fatalf("Print() error = %v", err)
		}
		want := masker.Mask(in)
		if buf.String() != want {
			t.Errorf("output = %q, want %q", buf.String(), want)
		}
		if strings.Contains(buf.String(), ghToken) {
			t.Errorf("token not masked: %q", buf.String())
		}
	})

	t.Run("no trailing newline", func(t *testing.T) {
		in := "no newline at end"
		var buf bytes.Buffer
		if err := newPassthroughGat(t).Print(&buf, strings.NewReader(in), WithMask(true)); err != nil {
			t.Fatalf("Print() error = %v", err)
		}
		if buf.String() != in {
			t.Errorf("output = %q, want %q", buf.String(), in)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		var buf bytes.Buffer
		if err := newPassthroughGat(t).Print(&buf, strings.NewReader(""), WithMask(true)); err != nil {
			t.Fatalf("Print() error = %v", err)
		}
		if buf.String() != "" {
			t.Errorf("output = %q, want empty", buf.String())
		}
	})

	t.Run("display option (show-ends), mask disabled", func(t *testing.T) {
		in := "a\nb\n"
		var buf bytes.Buffer
		err := newPassthroughGat(t).Print(&buf, strings.NewReader(in),
			WithMask(false), WithDisplay(&display.Options{ShowEnds: true}))
		if err != nil {
			t.Fatalf("Print() error = %v", err)
		}
		want := "a$\nb$\n"
		if buf.String() != want {
			t.Errorf("output = %q, want %q", buf.String(), want)
		}
	})
}

// A newline-split (malformed) PEM header is missed by per-line masking
// (intended divergence); a single-line PEM header is masked as before.
func TestGat_Print_passthrough_pemDivergence(t *testing.T) {
	t.Run("split PEM header is NOT masked (intended divergence)", func(t *testing.T) {
		in := "-----BEGIN\nRSA\nPRIVATE\nKEY-----\n"
		var buf bytes.Buffer
		if err := newPassthroughGat(t).Print(&buf, strings.NewReader(in), WithMask(true)); err != nil {
			t.Fatalf("Print() error = %v", err)
		}
		if buf.String() != in {
			t.Errorf("split PEM unexpectedly altered: output = %q, want %q", buf.String(), in)
		}
	})

	t.Run("single-line PEM header IS masked", func(t *testing.T) {
		in := "-----BEGIN RSA PRIVATE KEY-----\n"
		var buf bytes.Buffer
		if err := newPassthroughGat(t).Print(&buf, strings.NewReader(in), WithMask(true)); err != nil {
			t.Fatalf("Print() error = %v", err)
		}
		want := masker.Mask(in)
		if buf.String() != want {
			t.Errorf("output = %q, want %q", buf.String(), want)
		}
		if strings.Contains(buf.String(), "BEGIN RSA PRIVATE KEY") {
			t.Errorf("single-line PEM not masked: %q", buf.String())
		}
	})
}

func TestGat_Print_passthrough_maskWithDisplay(t *testing.T) {
	const ghToken = "ghp_0123456789abcdefghijklmnopqrstuvwxyz123456"
	in := "first\tcol " + ghToken + "\nsecond line\n"
	opts := &display.Options{ShowEnds: true, ShowTabs: true}

	// Oracle: full-text mask, then the same display transform.
	var want bytes.Buffer
	if _, err := display.NewWriter(&want, opts).Write([]byte(masker.Mask(in))); err != nil {
		t.Fatalf("building oracle: %v", err)
	}

	var got bytes.Buffer
	if err := newPassthroughGat(t).Print(&got, strings.NewReader(in),
		WithMask(true), WithDisplay(opts)); err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	if got.String() != want.String() {
		t.Errorf("output = %q, want %q", got.String(), want.String())
	}
}

func TestGat_Print_passthrough_multipleSecrets(t *testing.T) {
	const ghToken = "ghp_0123456789abcdefghijklmnopqrstuvwxyz123456"
	const oaToken = "sk-0123456789abcdefghijklmnopqrstuvwxyz"
	in := "header\ntoken1 " + ghToken + "\nmiddle\ntoken2 " + oaToken + "\nfooter\n"

	var buf bytes.Buffer
	if err := newPassthroughGat(t).Print(&buf, strings.NewReader(in), WithMask(true)); err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	if want := masker.Mask(in); buf.String() != want {
		t.Errorf("output = %q, want %q", buf.String(), want)
	}
	if strings.Contains(buf.String(), ghToken) || strings.Contains(buf.String(), oaToken) {
		t.Errorf("a token survived masking: %q", buf.String())
	}
}

// Non-passthrough configurations must still route through the highlighting
// path, so the streaming guard does not accidentally capture them.
func TestGat_Print_nonPassthroughIsHighlighted(t *testing.T) {
	const src = "package main\n\nfunc main() {}\n"

	t.Run("colored terminal theme emits ANSI escapes", func(t *testing.T) {
		g, err := New(&Config{Theme: "monokai", Format: "terminal256", Language: "go"})
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		var buf bytes.Buffer
		if err := g.Print(&buf, strings.NewReader(src)); err != nil {
			t.Fatalf("Print() error = %v", err)
		}
		if !strings.Contains(buf.String(), "\x1b[") {
			t.Errorf("expected ANSI escapes in highlighted output, got %q", buf.String())
		}
	})

	t.Run("html format emits markup", func(t *testing.T) {
		g, err := New(&Config{Theme: "monokai", Format: "html", Language: "go"})
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		var buf bytes.Buffer
		if err := g.Print(&buf, strings.NewReader(src)); err != nil {
			t.Fatalf("Print() error = %v", err)
		}
		if !strings.Contains(buf.String(), "<") {
			t.Errorf("expected HTML markup in output, got %q", buf.String())
		}
	})
}

// signalWriter records output and closes firstWrite on the first non-empty
// write, so a test can detect that output was emitted before EOF.
type signalWriter struct {
	mu         sync.Mutex
	buf        bytes.Buffer
	once       sync.Once
	firstWrite chan struct{}
}

func (w *signalWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(p) > 0 {
		w.once.Do(func() { close(w.firstWrite) })
	}
	return w.buf.Write(p)
}

func (w *signalWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}

// gatedReader serves its chunks, then blocks the next Read until release is
// closed before reporting EOF. It models a slow stream that has not yet
// terminated.
type gatedReader struct {
	chunks  [][]byte
	idx     int
	release chan struct{}
}

func (r *gatedReader) Read(p []byte) (int, error) {
	if r.idx < len(r.chunks) {
		n := copy(p, r.chunks[r.idx])
		r.idx++
		return n, nil
	}
	<-r.release
	return 0, io.EOF
}

// In passthrough mode, a sub-1024-byte chunk must be emitted without waiting
// for the input to terminate (i.e. content detection must not block for a full
// 1024-byte read).
func TestGat_Print_passthrough_streamsBeforeEOF(t *testing.T) {
	release := make(chan struct{})
	r := &gatedReader{chunks: [][]byte{[]byte("first\n")}, release: release}
	w := &signalWriter{firstWrite: make(chan struct{})}

	done := make(chan error, 1)
	go func() { done <- newPassthroughGat(t).Print(w, r) }()

	select {
	case <-w.firstWrite:
		// emitted before EOF — good
	case <-time.After(2 * time.Second):
		close(release)
		<-done
		t.Fatal("output was buffered: first chunk not emitted before input terminated")
	}

	close(release)
	if err := <-done; err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	if got := w.String(); got != "first\n" {
		t.Errorf("output = %q, want %q", got, "first\n")
	}
}

// Bulk gzip input is still detected and decompressed in passthrough mode (the
// lazy detection must fall back to the normal content-type routing).
func TestGat_Print_passthrough_gzipStillDecompresses(t *testing.T) {
	const content = "plain text content\n"
	var gz bytes.Buffer
	zw := gzip.NewWriter(&gz)
	if _, err := zw.Write([]byte(content)); err != nil {
		t.Fatalf("gzip write: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}

	var out bytes.Buffer
	if err := newPassthroughGat(t).Print(&out, bytes.NewReader(gz.Bytes())); err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	if out.String() != content {
		t.Errorf("gzip not decompressed: got %q, want %q", out.String(), content)
	}
}

// Bulk binary input is still guarded in passthrough mode rather than streamed
// raw to the terminal.
func TestGat_Print_passthrough_binaryStillGuarded(t *testing.T) {
	in := append([]byte("some text"), 0x00)
	in = append(in, []byte("more bytes")...)

	var out bytes.Buffer
	if err := newPassthroughGat(t).Print(&out, bytes.NewReader(in)); err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	if !strings.Contains(out.String(), "binary file") {
		t.Errorf("binary not guarded: got %q", out.String())
	}
}

// The masked streaming path (the motivating use case for --mask-secrets) also
// emits each line as it arrives, without waiting for the input to terminate.
func TestGat_Print_passthrough_streamsBeforeEOF_masked(t *testing.T) {
	const ghToken = "ghp_0123456789abcdefghijklmnopqrstuvwxyz123456"
	// The first chunk must contain a newline; the masked path reads line by line.
	firstLine := "secret " + ghToken + "\n"

	release := make(chan struct{})
	r := &gatedReader{chunks: [][]byte{[]byte(firstLine)}, release: release}
	w := &signalWriter{firstWrite: make(chan struct{})}

	done := make(chan error, 1)
	go func() { done <- newPassthroughGat(t).Print(w, r, WithMask(true)) }()

	select {
	case <-w.firstWrite:
		// emitted before EOF — good
	case <-time.After(2 * time.Second):
		close(release)
		<-done
		t.Fatal("masked output was buffered: first line not emitted before input terminated")
	}

	close(release)
	if err := <-done; err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	got := w.String()
	if strings.Contains(got, ghToken) {
		t.Errorf("token not masked: %q", got)
	}
	if want := masker.Mask(firstLine); got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

// Accepted limitation: gzip arriving in a first read too small to contain its
// magic is not recognized and streams through raw. Files and bulk streams are
// unaffected (their first read is large). This pins the trade-off documented
// on detectionHead.
func TestGat_Print_passthrough_gzipTrickleStreamsRaw(t *testing.T) {
	const content = "plain text content\n"
	var gz bytes.Buffer
	zw := gzip.NewWriter(&gz)
	if _, err := zw.Write([]byte(content)); err != nil {
		t.Fatalf("gzip write: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}
	raw := gz.Bytes()

	// Deliver one byte on the first read so detection cannot see the gzip magic.
	release := make(chan struct{})
	close(release)
	r := &gatedReader{chunks: [][]byte{raw[:1], raw[1:]}, release: release}

	var out bytes.Buffer
	if err := newPassthroughGat(t).Print(&out, r); err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	if out.String() != string(raw) {
		t.Errorf("expected raw gzip passthrough, got %d bytes (want %d)", out.Len(), len(raw))
	}
}

// Accepted limitation: a NUL byte arriving after the first read is not caught,
// so the content streams through. This mirrors the forceBinary CLI behavior on
// the common passthrough path.
func TestGat_Print_passthrough_binaryTrickleStreamsRaw(t *testing.T) {
	release := make(chan struct{})
	close(release)
	r := &gatedReader{chunks: [][]byte{[]byte("some text"), {0x00}, []byte("more")}, release: release}

	var out bytes.Buffer
	if err := newPassthroughGat(t).Print(&out, r); err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	if want := "some text\x00more"; out.String() != want {
		t.Errorf("expected raw passthrough, got %q", out.String())
	}
}
