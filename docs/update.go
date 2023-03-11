package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/styles"
	"github.com/koki-develop/gat/pkg/printer"
)

var (
	src = `package main

import "fmt"

func main() {
	fmt.Println("hello world")
}`
)

func String(s string) *string {
	return &s
}

func main() {
	updateThemes()
}

func updateThemes() {
	f, err := os.Create("docs/themes.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("# Highlight Themes\n\n")

	for _, s := range styles.Names() {
		f.WriteString(fmt.Sprintf("- [`%s`](#%s)\n", s, s))
	}
	f.WriteString("\n")

	for _, s := range styles.Names() {
		f.WriteString(fmt.Sprintf("## `%s`\n\n", s))

		p := printer.New(&printer.PrinterConfig{
			Format: "svg",
			Theme:  s,
		})

		b := new(bytes.Buffer)
		if err := p.Print(&printer.PrintInput{
			In:       strings.NewReader(src),
			Out:      b,
			Filename: String("main.go"),
		}); err != nil {
			panic(err)
		}

		img, err := os.Create(fmt.Sprintf("./docs/themes/%s.svg", s))
		if err != nil {
			panic(err)
		}
		defer img.Close()
		img.Write(b.Bytes())

		f.WriteString(fmt.Sprintf("![%s](./themes/%s.svg)\n\n", s, s))
	}
}
