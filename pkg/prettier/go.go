package prettier

import (
	"go/format"
)

var (
	Go = Register("Go", NewGoPrettier())
)

type GoPrettier struct{}

func NewGoPrettier() *GoPrettier {
	return &GoPrettier{}
}

func (*GoPrettier) Pretty(input string) (string, error) {
	b, err := format.Source([]byte(input))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
