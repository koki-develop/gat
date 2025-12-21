package masker

import (
	"regexp"
	"strings"
)

var patterns = []*regexp.Regexp{
	// AWS Access Key ID
	regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
	// GitHub Tokens (ghp_, gho_, ghs_, ghr_)
	regexp.MustCompile(`gh[pousr]_[a-zA-Z0-9]{36,}`),
	// GitLab Personal Access Token
	regexp.MustCompile(`glpat-[a-zA-Z0-9\-_]{20,}`),
	// Slack Tokens
	regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z\-]+`),
	// JWT Tokens
	regexp.MustCompile(`eyJ[a-zA-Z0-9_-]*\.eyJ[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*`),
	// Private Key Headers
	regexp.MustCompile(`-----BEGIN\s+(RSA|DSA|EC|OPENSSH|PGP)\s+PRIVATE\s+KEY-----`),
}

// Mask replaces sensitive patterns in content with asterisks of the same length
func Mask(content string) string {
	result := content
	for _, p := range patterns {
		result = p.ReplaceAllStringFunc(result, func(match string) string {
			return strings.Repeat("*", len(match))
		})
	}
	return result
}
