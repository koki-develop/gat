package prettier

import (
	"github.com/google/yamlfmt"
	"github.com/google/yamlfmt/formatters/basic"
)

var (
	YAML = Register("YAML", NewYAMLPrettier())
)

type YAMLPrettier struct {
	factory yamlfmt.Factory
}

func NewYAMLPrettier() *YAMLPrettier {
	return &YAMLPrettier{
		factory: &basic.BasicFormatterFactory{},
	}
}

func (p *YAMLPrettier) Pretty(input string) (string, error) {
	formatter, err := p.factory.NewFormatter(nil)
	if err != nil {
		return "", err
	}

	b, err := formatter.Format([]byte(input))
	if err != nil {
		return "", nil
	}

	return string(b), nil
}
