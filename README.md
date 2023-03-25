# gat

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/koki-develop/gat)](https://github.com/koki-develop/gat/releases/latest)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/koki-develop/gat/ci.yml?logo=github)](https://github.com/koki-develop/gat/actions/workflows/ci.yml)
[![Maintainability](https://img.shields.io/codeclimate/maintainability/koki-develop/gat?style=flat&logo=codeclimate)](https://codeclimate.com/github/koki-develop/gat/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/koki-develop/gat)](https://goreportcard.com/report/github.com/koki-develop/gat)
[![LICENSE](https://img.shields.io/github/license/koki-develop/gat)](./LICENSE)

cat alternative written in Go.

![demo](./docs/demo.gif)

- [Installation](#installation)
- [Usage](#usage)
- [LICENSE](#license)

## Installation

### Homebrew

```console
$ brew install koki-develop/tap/gat
```

### `go install`

```console
$ go install github.com/koki-develop/gat@latest
```

### Releases

Download the binary from the [releases page](https://github.com/koki-develop/gat/releases/latest).

## Usage

```console
$ gat --help
cat alternative written in Go.

Usage:
  gat [file]... [flags]

Flags:
  -f, --format string   output format (default "terminal256")
  -h, --help            help for gat
  -l, --lang string     language for syntax highlighting
      --list-formats    print a list of supported output formats
      --list-langs      print a list of supported languages for syntax highlighting
      --list-themes     print a list of supported themes with preview
  -p, --pretty          whether to format a content pretty
  -t, --theme string    highlight theme (default "monokai")
  -v, --version         version for gat
```

### `-l`, `--lang`

Explicitly set the language for syntax highlighting.  
See [languages.md](./docs/languages.md) for valid languages.

### `-f`, `--format`

Set the output format ( default: `terminal256` ).  
See [formats.md](./docs/formats.md) for valid formats.

### `-t`, `--theme`

Set the highlight theme ( default: `monokai` ).  
See [themes.md](./docs/themes.md) for valid thtmes.

## `-p`, `--pretty`

Format a content pretty.  
For unsupported languages, this flag is ignored.

## `-c`, `--force-color`

`gat` disables colored output when piped to another program.  
Settings the `--force-color` forces colored output to be enabled.  
This is useful, for example, when used in combination with the `less -R` command.

![](/docs/gess.gif)

It is also useful to declare the following function to allow `gat` to be used with a pager.

```sh
function gess() {
  gat --force-color "$@" | less -R
}
```

## LICENSE

[MIT](./LICENSE)
