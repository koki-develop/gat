package gat

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_printFormats(t *testing.T) {
	fmts := []string{
		"html",
		"html-min",
		"json",
		"json-min",
		"noop",
		"svg",
		"svg-min",
		"terminal",
		"terminal16",
		"terminal16m",
		"terminal256",
		"terminal8",
		"tokens",
	}

	want := `html
html-min
json
json-min
noop
svg
svg-min
terminal
terminal16
terminal16m
terminal256
terminal8
tokens
`

	buf := new(bytes.Buffer)
	err := printFormats(buf, fmts)

	assert.NoError(t, err)
	assert.Equal(t, want, buf.String())
}
