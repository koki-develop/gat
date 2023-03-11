# gat

cat alternative written in Go.

- [Installation](#installation)
- [Usage](#usage)
- [LICENSE](#license)

## Installation

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
  -t, --theme string    highlight theme (default "monokai")
```

### Format

`--format` flag to set the output format ( default: `terminal256` ).  
See [formats.md](./docs/formats.md) for valid formats.

### Theme

`--theme` flag to set the highlight theme ( default: `monokai` ).  
See [themes.md](./docs/themes.md) for valid thtmes.

## LICENSE

[MIT](./LICENSE)
