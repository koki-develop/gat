package gat

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_printLanguages(t *testing.T) {
	langs := []Language{
		{Name: "AAA", Aliases: []string{"a", "aa"}},
		{Name: "BB", Aliases: []string{"b"}},
		{Name: "CC", Aliases: []string{}},
		{Name: "DDD", Aliases: []string{"d", "ddd"}},
	}
	want := `AAA a, aa
BB  b
CC  
DDD d, ddd
`

	buf := new(bytes.Buffer)
	if err := printLanguages(buf, langs); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, buf.String())
}
