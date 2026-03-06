package display

import (
	"bytes"
	"testing"
)

func TestNewWriter_NoOptions(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, &Options{})
	if w != &buf {
		t.Error("expected original writer when no options are set")
	}
}

func TestShowTabs(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, &Options{ShowTabs: true})
	_, _ = w.Write([]byte("hello\tworld\n"))
	want := "hello^Iworld\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}

func TestShowEnds(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, &Options{ShowEnds: true})
	_, _ = w.Write([]byte("hello\nworld\n"))
	want := "hello$\nworld$\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}

func TestShowNonPrinting_ControlChars(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, &Options{ShowNonPrinting: true})
	_, _ = w.Write([]byte("hello\x01\x02world\n"))
	want := "hello^A^Bworld\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}

func TestShowNonPrinting_DEL(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, &Options{ShowNonPrinting: true})
	_, _ = w.Write([]byte("hello\x7fworld\n"))
	want := "hello^?world\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}

func TestShowNonPrinting_HighBit(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, &Options{ShowNonPrinting: true})
	_, _ = w.Write([]byte{0x80, 0xC1, 0xFF, '\n'})
	// 0x80 -> M-^@ (low=0x00, control)
	// 0xC1 -> M-A (low=0x41, printable)
	// 0xFF -> M-^? (low=0x7F, DEL)
	want := "M-^@M-AM-^?\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}

func TestShowAll(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, &Options{ShowNonPrinting: true, ShowEnds: true, ShowTabs: true})
	_, _ = w.Write([]byte("hello\t\x01world\n"))
	want := "hello^I^Aworld$\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}

func TestANSIEscapePreserved(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, &Options{ShowNonPrinting: true, ShowEnds: true, ShowTabs: true})
	// ANSI escape: ESC[31m (red) ... ESC[0m (reset)
	input := "\x1b[31mhello\x1b[0m\n"
	_, _ = w.Write([]byte(input))
	want := "\x1b[31mhello\x1b[0m$\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}

func TestShowNonPrinting_TabAndNewlineNotConverted(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, &Options{ShowNonPrinting: true})
	_, _ = w.Write([]byte("a\tb\n"))
	want := "a\tb\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}
