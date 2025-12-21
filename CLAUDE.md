# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

gat is a cat command alternative written in Go that provides syntax highlighting, code formatting, and enhanced display capabilities for terminal output.

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
```

### Building
```bash
# Build the project
go build -o gat

# Build with version information
go build -ldflags "-X main.version=X.Y.Z" -o gat
```

### Linting
```bash
# Run linter (golangci-lint must be installed)
golangci-lint run --verbose ./...
```

### Release Management
```bash
# Check GoReleaser configuration
goreleaser check

# Build release snapshot (for testing)
goreleaser release --snapshot --clean
```

## Architecture

### Core Components

- **cmd/**: CLI command definitions using Cobra framework
  - `root.go`: Main command logic and flag handling
  - `flags.go`: Command-line flag definitions
  - `version.go`: Version command implementation

- **internal/gat/**: Core gat functionality
  - Main logic for file processing, syntax highlighting, and output formatting

- **internal/formatters/**: Output format processors
  - HTML minification, JSON formatting, SVG optimization

- **internal/lexers/**: Custom syntax highlighters
  - Terraform lexer implementation

- **internal/prettier/**: Code prettifiers
  - Language-specific formatting (CSS, Go, HTML, JSON, XML, YAML)

- **internal/masker/**: Sensitive information masking
  - Masks API keys, tokens, and other secrets in output

- **internal/styles/**: Theme definitions
  - Custom syntax highlighting themes

- **scripts/**: Build scripts
  - Shell completion generation

- **docs/**: Documentation assets
  - Demo GIFs, images, theme previews

- **tapes/**: VHS tape files for generating demo GIFs

- **assets/**: Logo files

### Key Dependencies

- **Chroma**: Syntax highlighting engine (200+ language support)
- **Cobra**: CLI framework for command parsing
- **Glamour**: Markdown rendering with terminal styling
- **go-sixel**: Image display in terminal via Sixel protocol

### Design Principles

1. **Modular Architecture**: Each formatter, lexer, and prettifier is isolated in its own package
2. **Internal Packages**: Core functionality is kept in `internal/` to prevent external imports
3. **Resource Management**: Proper cleanup of file handles and resources
4. **Smart Output Detection**: Automatic color handling based on terminal/pipe detection

## Release Process

The project uses Release Please for automated releases:
1. PRs are automatically created with changelog updates
2. Merging a release PR triggers GoReleaser
3. Binaries are built for multiple platforms and published to GitHub Releases
4. Homebrew formula is automatically updated

## Testing Approach

- Unit tests focus on formatters, prettifiers, and core functionality
- Test files follow Go convention: `*_test.go` alongside implementation
- Use table-driven tests where appropriate
- Mock external dependencies when needed

## Masker Package Patterns

When adding new API key patterns to `internal/masker/`:

### Pattern Ordering
- Place more specific patterns before general ones to avoid false matches
- Example: `sk-ant-` must be before `sk-` to prevent Anthropic keys from matching OpenAI pattern
- Example: AWS Secret Access Key (`[a-zA-Z0-9+/]{40}`) must be last due to its generic pattern

### Supported Patterns (in order of application)
- AWS Access Key ID (permanent): `AKIA[0-9A-Z]{16}`
- AWS Access Key ID (temporary/SSO): `ASIA[0-9A-Z]{16}`
- GitHub Tokens: `gh[pousr]_[a-zA-Z0-9]{36,}`
- GitLab PAT: `glpat-[a-zA-Z0-9\-_]{20,}`
- Slack Tokens: `xox[baprs]-[0-9a-zA-Z\-]+`
- Anthropic API Key: `sk-ant-[a-zA-Z0-9\-_]+`
- OpenAI API Key: `sk-(?:proj-)?[a-zA-Z0-9_\-]{20,}` (supports both legacy and project formats)
- Supabase Secret Key: `sb_secret_[a-zA-Z0-9\-_]+`
- JWT Tokens: `eyJ[a-zA-Z0-9_-]*\.eyJ[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*`
- Private Key Headers: `-----BEGIN\s+(RSA|DSA|EC|OPENSSH|PGP)\s+PRIVATE\s+KEY-----`
- AWS Secret Access Key: `[a-zA-Z0-9+/]{40}` (must be last due to generic pattern)

### Pattern Update Workflow
1. Add regex pattern to `internal/masker/masker.go`
2. Add test cases to `internal/masker/masker_test.go` with realistic examples
3. Run tests: `go test ./internal/masker/...`
4. Update README.md supported patterns list
5. Test cases should include special characters (`_`, `-`) where applicable