package masker

import (
	"regexp"
	"strings"
)

var patterns = []*regexp.Regexp{
	// AWS Access Key ID (permanent)
	regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
	// AWS Access Key ID (temporary, STS/SSO)
	regexp.MustCompile(`ASIA[0-9A-Z]{16}`),
	// GitHub Tokens (ghp_, gho_, ghs_, ghr_)
	regexp.MustCompile(`gh[pousr]_[a-zA-Z0-9]{36,}`),
	// GitLab Personal Access Token
	regexp.MustCompile(`glpat-[a-zA-Z0-9\-_]{20,}`),
	// Slack Tokens
	regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z\-]+`),
	// Anthropic API Key (must be before OpenAI to avoid false matches)
	regexp.MustCompile(`sk-ant-[a-zA-Z0-9\-_]+`),
	// OpenAI API Key (both legacy sk- and new sk-proj- formats)
	regexp.MustCompile(`sk-(?:proj-)?[a-zA-Z0-9_\-]{20,}`),
	// Supabase Secret Key
	regexp.MustCompile(`sb_secret_[a-zA-Z0-9\-_]+`),
	// JWT Tokens
	regexp.MustCompile(`eyJ[a-zA-Z0-9_-]*\.eyJ[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*`),
	// Private Key Headers
	regexp.MustCompile(`-----BEGIN\s+(RSA|DSA|EC|OPENSSH|PGP)\s+PRIVATE\s+KEY-----`),
	// AWS Secret Access Key (must be last due to generic pattern that could match other secrets)
	regexp.MustCompile(`[a-zA-Z0-9+/]{40}`),
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
