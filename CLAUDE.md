# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Testing
```bash
go test ./...                    # Run all tests
go test ./internal/gat          # Test core functionality
go test ./internal/prettier     # Test code formatting
```

### Building
```bash
go build                        # Build for current platform
goreleaser check                 # Validate goreleaser config
goreleaser release --snapshot --clean  # Build cross-platform binaries
```

### Linting
```bash
golangci-lint run --verbose ./...  # Run linter (configured via mise.toml)
```

### Running
```bash
go run . [file]...               # Run from source
./gat [file]...                  # Run built binary
```

## Architecture

### Core Components

**CLI Layer (`cmd/`)**
- `root.go`: Main cobra command setup with terminal detection logic
- `flags.go`: All CLI flags and environment variable handling (`GAT_FORMAT`, `GAT_THEME`)
- `version.go`: Version command implementation

**Core Engine (`internal/gat/`)**
- `gat.go`: Main `Gat` struct and `Print()` method - handles content detection, lexer selection, and output formatting
- `formats.go` & `languages.go`: List available output formats and supported languages

**Formatters (`internal/formatters/`)**
- Wrapper around Chroma formatters for terminal, HTML, JSON, SVG output
- Includes minified variants for web formats

**Lexers (`internal/lexers/`)**
- Language detection logic using filename and content analysis
- Integrates with Chroma lexer registry

**Code Prettification (`internal/prettier/`)**
- Language-specific formatters (Go, JSON, HTML, CSS, XML, YAML)
- Registry pattern for extensible formatting support
- `fallback.go`: Default pass-through for unsupported languages

**Themes (`internal/styles/`)**
- Custom theme registry including `noop.xml` for no-color output
- Wraps Chroma's style system

### Key Data Flow

- **Input Processing**: `gat.Print()` reads content and detects MIME type
- **Content Detection**: Binary vs text, with special handling for images and gzip
- **Lexer Selection**: Auto-detect language from filename/content or use explicit `--lang`
- **Content Transformation**: 
  - Markdown rendering (glamour) if `--render-markdown`
  - Code prettification if `--pretty`
- **Output Formatting**: Apply syntax highlighting and format (terminal/HTML/JSON/SVG)

### Special Features

**Image Handling**
- Sixel encoding for terminal image display
- Automatic resizing (max 1800px edge, disabled with `--no-resize`)
- Supports JPEG, PNG, GIF

**Terminal Behavior**
- Auto-detects piped output and disables colors (unless `--force-color`)
- Uses `noop` theme for non-terminal output
- Forces binary output when piped

**Environment Integration**
- `GAT_FORMAT` and `GAT_THEME` environment variables
- Terminal detection for appropriate output formatting

## Tool Configuration

The project uses:
- **mise.toml**: Tool version management (Go 1.24.3, golangci-lint, vhs, goreleaser)
- **goreleaser**: Cross-platform building and release automation
- **GitHub Actions**: CI with test/build/lint jobs

## Testing Strategy

Tests focus on:
- Format and language listing functionality
- Code prettification for each supported language
- Error handling and edge cases in prettier modules