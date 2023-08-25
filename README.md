<p align="center">
<img src="./assets/logo_light.svg#gh-light-mode-only" >
<img src="./assets/logo_dark.svg#gh-dark-mode-only" >
</p>

<p align="center">
cat alternative written in Go.
</p>

<p align="center">
<a href="https://github.com/koki-develop/gat/releases/latest"><img src="https://img.shields.io/github/v/release/koki-develop/gat" alt="GitHub release (latest by date)"></a>
<a href="https://github.com/koki-develop/gat/releases/latest"><img alt="GitHub all releases" src="https://img.shields.io/github/downloads/koki-develop/gat/total?style=flat"></a>
<a href="https://github.com/koki-develop/gat/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/koki-develop/gat/ci.yml?logo=github" alt="GitHub Workflow Status"></a>
<a href="https://codeclimate.com/github/koki-develop/gat/maintainability"><img src="https://img.shields.io/codeclimate/maintainability/koki-develop/gat?style=flat&amp;logo=codeclimate" alt="Maintainability"></a>
<a href="https://goreportcard.com/report/github.com/koki-develop/gat"><img src="https://goreportcard.com/badge/github.com/koki-develop/gat" alt="Go Report Card"></a>
<a href="./LICENSE"><img src="https://img.shields.io/github/license/koki-develop/gat" alt="LICENSE"></a>
</p>

<p align="center">
<img src="./docs/demo.gif" >
</p>

## Contents

- [Installation](#installation)
- [Usage](#usage)
- [LICENSE](#license)

## Installation

### Homebrew

```console
$ brew install gat
```

### Homebrew Tap

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
  -c, --force-color     force colored output
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
Alternatively, it can be set using the `GAT_FORMAT` environment variable.
See [formats.md](./docs/formats.md) for valid formats.

### `-t`, `--theme`

Set the highlight theme ( default: `monokai` ).  
Alternatively, it can be set using the `GAT_THEME` environment variable.
See [themes.md](./docs/themes.md) for valid themes.

### `-p`, `--pretty`

Format a content pretty.  
For unsupported languages, this flag is ignored.

### `-c`, `--force-color`

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

### Print Image

If your terminal supports Sixel, you can print images.

![](./docs/image.png)

Supported image formats include:

- JPEG
- PNG
- GIF (animation not supported)

## LICENSE

[MIT](./LICENSE)
