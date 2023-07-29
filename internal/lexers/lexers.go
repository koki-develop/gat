package lexers

import (
	"fmt"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

type option struct {
	Language string
	Filename string
	Source   string
}

type Option func(*option)

func WithLanguage(lang string) Option {
	return func(o *option) {
		o.Language = lang
	}
}

func WithFilename(name string) Option {
	return func(o *option) {
		o.Filename = name
	}
}

func WithSource(src string) Option {
	return func(o *option) {
		o.Source = src
	}
}

func Get(opts ...Option) (chroma.Lexer, error) {
	opt := &option{}
	for _, o := range opts {
		o(opt)
	}

	var l chroma.Lexer
	if opt.Language != "" {
		l = lexers.Get(opt.Language)
		if l == nil {
			return nil, fmt.Errorf("unknown language: %s", opt.Language)
		}
	}

	if l == nil && opt.Filename != "" {
		l = lexers.Match(opt.Filename)
	}

	if l == nil && opt.Source != "" {
		l = lexers.Analyse(opt.Source)
	}

	if l == nil {
		l = lexers.Fallback
	}
	return chroma.Coalesce(l), nil
}
