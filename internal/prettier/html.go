package prettier

import (
	"github.com/yosssi/gohtml"
)

var (
	HTML = Register("HTML", NewHTMLPrettier())
)

type HTMLPrettier struct{}

func NewHTMLPrettier() *HTMLPrettier {
	return &HTMLPrettier{}
}

func (p *HTMLPrettier) Pretty(h string) (string, error) {
	return gohtml.Format(h), nil
}
