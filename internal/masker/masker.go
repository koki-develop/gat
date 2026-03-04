package masker

import (
	"regexp"
	"strings"
)

var patterns = []*regexp.Regexp{
	// AWS Access Key ID (permanent)
	regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`),
	// AWS Access Key ID (temporary, STS/SSO)
	regexp.MustCompile(`\bASIA[0-9A-Z]{16}\b`),
	// GitHub Tokens (ghp_, gho_, ghs_, ghr_)
	regexp.MustCompile(`\bgh[pousr]_[a-zA-Z0-9]{36,}\b`),
	// GitLab Personal Access Token
	regexp.MustCompile(`\bglpat-[a-zA-Z0-9\-_]{20,}\b`),
	// Slack Tokens
	regexp.MustCompile(`\bxox[baprs]-[0-9a-zA-Z\-]+\b`),
	// Anthropic API Key (must be before OpenAI to avoid false matches)
	regexp.MustCompile(`\bsk-ant-[a-zA-Z0-9\-_]+\b`),
	// OpenAI API Key (both legacy sk- and new sk-proj- formats)
	regexp.MustCompile(`\bsk-(?:proj-)?[a-zA-Z0-9_\-]{20,}\b`),
	// Supabase Secret Key
	regexp.MustCompile(`\bsb_secret_[a-zA-Z0-9\-_]+\b`),
	// JWT Tokens
	regexp.MustCompile(`\beyJ[a-zA-Z0-9_-]*\.eyJ[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*\b`),
	// Private Key Headers
	regexp.MustCompile(`-----BEGIN\s+(RSA|DSA|EC|OPENSSH|PGP)\s+PRIVATE\s+KEY-----`),
	// AWS Secret Access Key (must be last due to generic pattern that could match other secrets)
	regexp.MustCompile(`\b[a-zA-Z0-9+/]{40}\b`),
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
