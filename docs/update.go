package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/koki-develop/gat/internal/formatters"
	"github.com/koki-develop/gat/internal/gat"
	"github.com/koki-develop/gat/internal/lexers"
	"github.com/koki-develop/gat/internal/styles"
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

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func OrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func Map[T, I any](ts []T, f func(t T) I) []I {
	is := []I{}
	for _, t := range ts {
		is = append(is, f(t))
	}
	return is
}

func updateLanguages() {
	f := Must(os.Create("docs/languages.md"))
	defer func() { _ = f.Close() }()

	Must(f.WriteString("# Languages\n\n"))

	Must(f.WriteString("| Language | Aliases |\n"))
	Must(f.WriteString("| --- | --- |\n"))

	for _, l := range lexers.List() {
		cfg := l.Config()
		Must(fmt.Fprintf(f, "| `%s` ", cfg.Name))

		if len(cfg.Aliases) > 0 {
			Must(fmt.Fprintf(
				f,
				"| %s |",
				strings.Join(
					Map(
						cfg.Aliases,
						func(a string) string { return fmt.Sprintf("`%s`", a) },
					),
					", ",
				),
			))
		} else {
			Must(f.WriteString("| |"))
		}
		Must(f.WriteString("\n"))
	}
}

func updateFormats() {
	f := Must(os.Create("docs/formats.md"))
	defer func() { _ = f.Close() }()

	Must(f.WriteString("# Output Formats\n\n"))

	for _, format := range formatters.List() {
		Must(fmt.Fprintf(f, "- [`%s`](#%s)\n", format, format))
	}
	Must(f.WriteString("\n"))

	for _, format := range formatters.List() {
		Must(fmt.Fprintf(f, "## `%s`\n\n", format))

		g := Must(gat.New(&gat.Config{
			Format:   format,
			Theme:    "monokai",
			Language: "go",
		}))

		b := new(bytes.Buffer)
		OrPanic(g.Print(b, strings.NewReader(src)))

		Must(fmt.Fprintf(f, "```%s\n", strings.TrimSuffix(format, "-min")))
		if strings.HasPrefix(format, "terminal") {
			Must(f.WriteString(strings.Trim(strings.ReplaceAll(fmt.Sprintf("%#v", strings.TrimSpace(b.String())), "\\n", "\n"), "\"")))
		} else {
			Must(f.WriteString(strings.TrimSpace(b.String())))
		}
		Must(f.WriteString("\n"))
		Must(f.WriteString("```\n"))

		Must(f.WriteString("\n"))
	}
}

func updateThemes() {
	f := Must(os.Create("docs/themes.md"))
	defer func() { _ = f.Close() }()

	Must(f.WriteString("# Highlight Themes\n\n"))

	for _, s := range styles.List() {
		Must(fmt.Fprintf(f, "- [`%s`](#%s)\n", s, s))
	}
	Must(f.WriteString("\n"))

	for _, s := range styles.List() {
		Must(fmt.Fprintf(f, "## `%s`\n\n", s))

		g := Must(gat.New(&gat.Config{
			Format:   "svg",
			Theme:    s,
			Language: "go",
		}))

		b := new(bytes.Buffer)
		OrPanic(g.Print(b, strings.NewReader(src)))

		img := Must(os.Create(fmt.Sprintf("./docs/themes/%s.svg", s)))
		defer func() { _ = img.Close() }()
		Must(img.Write(b.Bytes()))

		Must(fmt.Fprintf(f, "![%s](./themes/%s.svg)\n\n", s, s))
	}
}
