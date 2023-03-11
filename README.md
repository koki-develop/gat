# gat

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/koki-develop/gat)](https://github.com/koki-develop/gat/releases/latest)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/koki-develop/gat/build.yml?logo=github)](https://github.com/koki-develop/gat/actions/workflows/build.yml)
[![Maintainability](https://img.shields.io/codeclimate/maintainability/koki-develop/gat?style=flat&logo=codeclimate)](https://codeclimate.com/github/koki-develop/gat/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/koki-develop/gat)](https://goreportcard.com/report/github.com/koki-develop/gat)
[![LICENSE](https://img.shields.io/github/license/koki-develop/gat)](./LICENSE)

cat alternative written in Go.

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
  -t, --theme string    highlight theme (default "monokai")
  -v, --version         version for gat
```

### Language

`--lang` flag to set the language for syntax highlighting.  
See [languages.md](./docs/languages.md) for valid languages.

### Format

`--format` flag to explicitly set the output format ( default: `terminal256` ).  
See [formats.md](./docs/formats.md) for valid formats.

### Theme

`--theme` flag to set the highlight theme ( default: `monokai` ).  
See [themes.md](./docs/themes.md) for valid thtmes.

## LICENSE

[MIT](./LICENSE)
