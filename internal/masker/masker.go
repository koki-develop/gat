package masker

import (
	"regexp"
	"strings"
)

// pattern pairs a secret-matching regular expression with an optional
// validator. When validate is non-nil, a regex match is masked only if
// validate reports it as a genuine secret. This lets generic patterns reject
// look-alikes that happen to share the same shape.
type pattern struct {
	re       *regexp.Regexp
	validate func(string) bool
}

var patterns = []pattern{
	// AWS Access Key ID (permanent)
	{re: regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`)},
	// AWS Access Key ID (temporary, STS/SSO)
	{re: regexp.MustCompile(`\bASIA[0-9A-Z]{16}\b`)},
	// GitHub Tokens (ghp_, gho_, ghs_, ghr_)
	{re: regexp.MustCompile(`\bgh[pousr]_[a-zA-Z0-9]{36,}\b`)},
	// GitHub Fine-grained Personal Access Token (github_pat_ + 82 word chars)
	{re: regexp.MustCompile(`\bgithub_pat_\w{82}\b`)},
	// GitLab Personal Access Token
	{re: regexp.MustCompile(`\bglpat-[a-zA-Z0-9\-_]{20,}\b`)},
	// Slack Tokens
	{re: regexp.MustCompile(`\bxox[baprs]-[0-9a-zA-Z\-]+\b`)},
	// Slack App-level Token (xapp-)
	{re: regexp.MustCompile(`\bxapp-\d-[A-Z0-9]+-\d+-[a-z0-9]+\b`)},
	// Anthropic API Key (must be before OpenAI to avoid false matches)
	{re: regexp.MustCompile(`\bsk-ant-[a-zA-Z0-9\-_]+\b`)},
	// OpenAI API Key (both legacy sk- and new sk-proj- formats)
	{re: regexp.MustCompile(`\bsk-(?:proj-)?[a-zA-Z0-9_\-]{20,}\b`)},
	// Supabase Secret Key
	{re: regexp.MustCompile(`\bsb_secret_[a-zA-Z0-9\-_]+\b`)},
	// npm Access Token (npm_ + 36 base62 chars)
	{re: regexp.MustCompile(`\bnpm_[a-zA-Z0-9]{36}\b`)},
	// PyPI API Token (pypi- + base64-serialized macaroon, always prefixed
	// with the fixed AgEIcHlwaS5vcmc that encodes the "pypi.org" location)
	{re: regexp.MustCompile(`\bpypi-AgEIcHlwaS5vcmc[a-zA-Z0-9_\-]{50,}`)},
	// RubyGems API Key (rubygems_ + hex; the value is SecureRandom.hex(16) =
	// 32 hex chars, but match 32+ to also cover longer scanner-reported forms)
	{re: regexp.MustCompile(`\brubygems_[a-f0-9]{32,}\b`)},
	// JWT Tokens
	{re: regexp.MustCompile(`\beyJ[a-zA-Z0-9_-]*\.eyJ[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*\b`)},
	// Private Key Headers
	{re: regexp.MustCompile(`-----BEGIN\s+(RSA|DSA|EC|OPENSSH|PGP)\s+PRIVATE\s+KEY-----`)},
	// AWS Secret Access Key (must be last due to generic pattern that could match other secrets).
	// A real key is base64 of 30 random bytes, so it almost always mixes upper- and
	// lower-case letters. Requiring mixed case rejects common 40-char look-alikes such as
	// Git SHA-1 hashes, which are lower-case hex only.
	{re: regexp.MustCompile(`\b[a-zA-Z0-9+/]{40}\b`), validate: hasMixedCase},
}

// hasMixedCase reports whether s contains at least one ASCII upper-case letter
// and at least one ASCII lower-case letter.
func hasMixedCase(s string) bool {
	var hasUpper, hasLower bool
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		}
		if hasUpper && hasLower {
			return true
		}
	}
	return false
}

// Mask replaces sensitive patterns in content with asterisks of the same length
func Mask(content string) string {
	result := content
	for _, p := range patterns {
		result = p.re.ReplaceAllStringFunc(result, func(match string) string {
			if p.validate != nil && !p.validate(match) {
				return match
			}
			return strings.Repeat("*", len(match))
		})
	}
	return result
}
