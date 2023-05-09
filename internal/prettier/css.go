package prettier

import (
	"bytes"
	"strings"

	"github.com/client9/csstool"
)

var (
	CSS = Register("CSS", NewCSSPrettier())
)

type CSSPrettier struct{}

func NewCSSPrettier() *CSSPrettier {
	return &CSSPrettier{}
}

func (p *CSSPrettier) Pretty(c string) (string, error) {
	f := csstool.NewCSSFormat(2, false, nil)
	f.AlwaysSemicolon = true

	b := new(bytes.Buffer)
	if err := f.Format(strings.NewReader(c), b); err != nil {
		return "", err
	}

	return b.String(), nil
}
