package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
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
	updateLanguages()
	updateThemes()
	updateFormats()
}

func updateLanguages() {
	f, err := os.Create("docs/languages.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("# Languages\n\n")

	for _, l := range lexers.Registry.Lexers {
		cfg := l.Config()
		f.WriteString(fmt.Sprintf("- `%s`", cfg.Name))
		if len(cfg.Aliases) > 0 {
			f.WriteString("( alias: ")
			aliases := []string{}
			for _, a := range cfg.Aliases {
				aliases = append(aliases, fmt.Sprintf("`%s`", a))
			}
			f.WriteString(strings.Join(aliases, ", "))
			f.WriteString(" )")
		}
		f.WriteString("\n")
	}
}

func updateFormats() {
	f, err := os.Create("docs/formats.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("# Output Formats\n\n")

	for _, format := range formatters.Names() {
		f.WriteString(fmt.Sprintf("- [`%s`](#%s)\n", format, format))
	}
	f.WriteString("\n")

	for _, format := range formatters.Names() {
		f.WriteString(fmt.Sprintf("## `%s`\n\n", format))

		p := printer.New(&printer.PrinterConfig{
			Format: format,
			Theme:  printer.DefaultTheme,
		})

		b := new(bytes.Buffer)
		if err := p.Print(&printer.PrintInput{
			In:       strings.NewReader(src),
			Out:      b,
			Filename: String("main.go"),
		}); err != nil {
			panic(err)
		}

		f.WriteString(fmt.Sprintf("```%s\n", format))
		f.WriteString(strings.TrimSpace(b.String()))
		f.WriteString("\n")
		f.WriteString("```\n")

		f.WriteString("\n")
	}
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
