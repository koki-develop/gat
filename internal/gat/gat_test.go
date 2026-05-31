package gat

import (
	"bytes"
	"strings"
	"testing"

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
