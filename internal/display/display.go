package display

import "io"

// Options controls cat-compatible display transformations.
type Options struct {
	ShowNonPrinting bool // -v
	ShowEnds        bool // -E
	ShowTabs        bool // -T
}

// NewWriter returns an io.Writer that applies display transformations.
// If no options are enabled, it returns w unchanged.
func NewWriter(w io.Writer, opts *Options) io.Writer {
	if !opts.ShowNonPrinting && !opts.ShowEnds && !opts.ShowTabs {
		return w
	}
	return &writer{w: w, opts: opts}
}

type writer struct {
	w        io.Writer
	opts     *Options
	inEscape bool
}

func (dw *writer) Write(p []byte) (int, error) {
	var buf []byte
	for _, b := range p {
		if dw.inEscape {
			buf = append(buf, b)
			// ESC + one byte (0x40-0x7E) ends the sequence or CSI introducer
			if b >= 0x40 && b <= 0x7E {
				dw.inEscape = false
			}
			continue
		}

		if b == '\x1b' {
			dw.inEscape = true
			buf = append(buf, b)
			continue
		}

		switch {
		case b == '\n':
			if dw.opts.ShowEnds {
				buf = append(buf, '$')
			}
			buf = append(buf, '\n')
		case b == '\t':
			if dw.opts.ShowTabs {
				buf = append(buf, '^', 'I')
			} else {
				buf = append(buf, '\t')
			}
		case b < 0x20 && dw.opts.ShowNonPrinting:
			// Control characters (except \t, \n handled above)
			buf = append(buf, '^', b+'@')
		case b == 0x7F && dw.opts.ShowNonPrinting:
			buf = append(buf, '^', '?')
		case b >= 0x80 && dw.opts.ShowNonPrinting:
			buf = append(buf, 'M', '-')
			low := b & 0x7F
			if low < 0x20 {
				buf = append(buf, '^', low+'@')
			} else if low == 0x7F {
				buf = append(buf, '^', '?')
			} else {
				buf = append(buf, low)
			}
		default:
			buf = append(buf, b)
		}
	}

	if len(buf) > 0 {
		_, err := dw.w.Write(buf)
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
