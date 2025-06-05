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

### Coverage
```bash
go test -cover ./...             # Run tests with coverage
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out  # Generate HTML coverage report
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
- **mise**: Tool version management (Go 1.24.3, golangci-lint, vhs, goreleaser)
- **goreleaser**: Cross-platform building and release automation
- **GitHub Actions**: CI with test/build/lint jobs

## Technical Challenges & Improvements

### 🚨 Critical Issues (Fix Immediately)

#### 1. Resource Leak in File Handling
**File**: `cmd/root.go:54-63`
**Issue**: Files are deferred to close at function end, not loop iteration end
```go
for _, filename := range args {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer func() { _ = f.Close() }() // BUG: All files close at function end
    // ...
}
```
**Fix**: Close files immediately after processing each one:
```go
for _, filename := range args {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    err = g.Print(os.Stdout, f, gat.WithPretty(flagPretty), gat.WithFilename(filename))
    f.Close() // Close immediately
    if err != nil {
        return err
    }
}
```

#### 2. Extremely Low Test Coverage (4.7%)
**Missing Coverage**:
- Core functionality (`internal/gat/gat.go`)
- CLI layer (`cmd/`)
- Essential modules (`internal/lexers/`, `internal/formatters/`, `internal/styles/`)

**Needed Tests**:
- Integration tests for end-to-end workflows
- Table-driven tests for binary detection edge cases
- Error injection tests for file I/O operations
- Tests for different image formats and sizes

#### 3. Memory Safety Issues
**File**: `internal/gat/gat.go:142-146, 202-236`
**Issues**:
- No size limits when loading files/images into memory
- Potential DoS via large file attacks
- Image processing lacks memory bounds

**Recommendations**:
- Add configurable size limits
- Implement streaming processing for large files
- Add timeout controls for I/O operations

### ⚠️ High Priority Issues

#### 4. Silent Error Suppression
**Files**: `cmd/root.go:59`, `internal/gat/gat.go:166, 244`
```go
defer func() { _ = f.Close() }() // Ignoring close errors
```
**Fix**: Log errors or return them appropriately:
```go
defer func() {
    if err := f.Close(); err != nil {
        // Log error or handle appropriately
    }
}()
```

#### 5. Performance Bottlenecks
**String Concatenation**: `internal/gat/gat.go:142-146`
```go
buf := new(bytes.Buffer)
if _, err := io.Copy(buf, br); err != nil {
    return err
}
src = buf.String() // Unnecessary copy
```
**Fix**: Use `buf.Bytes()` and work with byte slices where possible.

**Image Processing**: `internal/gat/gat.go:202-236`
- Creates new RGBA image in memory for resizing without size limits
- No memory pool for frequent image operations
- Loads entire image into memory before size checking

#### 6. Hard-coded Magic Values
**File**: `internal/gat/gat.go:103-135`
- Binary message string is hard-coded and very long (line 135)
- Magic number 1024 for peek size and binary detection (lines 104, 254-257)
- Magic number 1800 for image resize (line 203)

**Fix**: Extract as constants:
```go
const (
    DefaultPeekSize = 1024
    MaxImageEdge = 1800
    BinaryFileMessage = "+----------------------------------------------------------------------------+\n| NOTE: This is a binary file. To force output, use the --force-binary flag. |\n+----------------------------------------------------------------------------+\n"
)
```

### 🔧 Medium Priority Issues

#### 7. Global State in Prettier Registry
**File**: `internal/prettier/registry.go:3`
```go
var Registry = map[string]Prettier{} // Global mutable state
```
**Fix**: Use dependency injection or constructor pattern:
```go
type Registry struct {
    prettiers map[string]Prettier
}

func NewRegistry() *Registry {
    return &Registry{prettiers: make(map[string]Prettier)}
}
```

#### 8. Tight Coupling
**File**: `internal/gat/gat.go:44-75`
The `New()` function directly couples to multiple internal packages, making it hard to test in isolation.

#### 9. Security Concerns
**Dependency Vulnerabilities**:
- Several dependencies at v0.x versions
- `github.com/yosssi/gohtml` (last updated 2020)
- Various packages with commit hashes instead of proper versions

**Potential DoS**:
- No size limits when reading files into memory
- Could cause OOM with very large files
- No timeout for file operations

#### 10. Documentation Issues
**Missing Godoc Comments**:
- `type Config struct` - no documentation for fields
- `func New()` - no documentation for error conditions
- `func (g *Gat) Print()` - no documentation for complex behavior

**Missing Architecture Documentation**:
- How lexer selection works
- Image processing pipeline
- Error handling strategy

#### 11. Dependency Management
**Outdated Dependencies**:
- `github.com/alecthomas/chroma/v2` v2.17.2 → v2.18.0
- `github.com/bmatcuk/doublestar/v4` v4.7.1 → v4.8.1

**Inconsistent Versioning**:
- Some dependencies use commit hashes instead of semantic versions

### 📈 Action Plan

#### Phase 1 (Week 1) - Critical Fixes
- [ ] Fix file handle leak in `cmd/root.go:54-63`
- [ ] Add basic integration tests to achieve >50% coverage
- [ ] Add size limits for file and image processing
- [ ] Fix error handling - don't silently ignore errors

#### Phase 2 (Week 2) - High Priority
- [ ] Extract magic numbers to constants
- [ ] Add comprehensive error context
- [ ] Improve string concatenation performance
- [ ] Update critical dependencies

#### Phase 3 (Month 2) - Medium Priority
- [ ] Refactor global prettier registry
- [ ] Add comprehensive test suite (target 70% coverage)
- [ ] Optimize image processing with streaming
- [ ] Add performance benchmarks

#### Phase 4 (Future) - Enhancements
- [ ] Improve documentation coverage
- [ ] Consider adding configuration file support
- [ ] Add memory profiling and optimization
- [ ] Implement progressive image loading

## Testing Strategy

**Current Coverage**: 4.7% (critically low)

Tests focus on:
- Format and language listing functionality
- Code prettification for each supported language
- Error handling and edge cases in prettier modules

### Current Coverage Issues
- **Core functionality**: 0% coverage for binary detection, image processing, gzip handling
- **CLI layer**: 0% coverage for flag parsing, terminal detection
- **Essential modules**: 0% coverage for lexers, formatters, styles

### Recommended Test Structure
```
tests/
├── integration/
│   ├── cli_test.go          # End-to-end CLI testing
│   ├── formats_test.go      # Output format testing
│   └── images_test.go       # Image processing workflows
├── unit/
│   ├── binary_test.go       # Binary detection edge cases
│   ├── lexers_test.go       # Language detection
│   └── prettier_test.go     # Code formatting
└── fixtures/
    ├── samples/             # Test files for each format
    └── images/              # Test images of various sizes
```

### Performance Testing
- [ ] Add benchmarks for large file processing
- [ ] Add memory usage tests
- [ ] Add concurrent access tests
- [ ] Add stress tests for resource limits