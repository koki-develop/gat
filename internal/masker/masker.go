package masker

import (
	"regexp"
	"strings"
)

// pattern pairs a secret-matching regular expression with optional refinements.
// When validate is non-nil, a match is masked only if validate reports the
// masked portion as a genuine secret. When maskGroup is > 0, only that capture
// group is replaced with asterisks instead of the whole match — used for
// contextual patterns where the regex also has to match a surrounding keyword
// to gain confidence, but only the secret value should be masked.
type pattern struct {
	re        *regexp.Regexp
	validate  func(string) bool
	maskGroup int
}

var patterns = []pattern{
	// AWS Access Key ID (permanent)
	{re: regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`)},
	// AWS Access Key ID (temporary, STS/SSO)
	{re: regexp.MustCompile(`\bASIA[0-9A-Z]{16}\b`)},
	// GitHub App installation token (stateless, JWT format). GitHub is rolling
	// out a new ghs_-prefixed installation token that is a ~520-char JWT and so
	// contains dots, hyphens, and underscores. The generic GitHub pattern below
	// only matches [a-zA-Z0-9], so it would stop at the first dot (or miss the
	// token entirely when the leading segment is under 36 chars); this pattern,
	// placed first so it masks the whole token, follows GitHub's own guidance of
	// ghs_[A-Za-z0-9.\-_]{36,}. No trailing \b because the value may end in a
	// non-word char (- or _), like the SendGrid and PyPI patterns below.
	// https://github.blog/changelog/2026-05-15-github-app-installation-tokens-per-request-override-header/
	{re: regexp.MustCompile(`\bghs_[a-zA-Z0-9._\-]{36,}`)},
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
	// Google (GCP) / Firebase API Key (AIza + 35 chars). Firebase web API keys
	// share this same AIza format, so this single pattern covers both.
	{re: regexp.MustCompile(`\bAIza[0-9A-Za-z_\-]{35}\b`)},
	// Stripe Secret / Restricted API Key. Per Stripe docs the modes are test,
	// live, and org (organization keys, sk_org_); prod is kept for gitleaks parity.
	// Distinct from OpenAI's sk- thanks to the underscore separator, so there is
	// no overlap with the OpenAI pattern above.
	{re: regexp.MustCompile(`\b(?:sk|rk)_(?:test|live|prod|org)_[a-zA-Z0-9]{10,99}\b`)},
	// SendGrid API Key (SG. + 66 chars). No trailing \b because the value may
	// end in a non-word character (=, ., -), like the PyPI pattern above.
	{re: regexp.MustCompile(`\bSG\.[a-zA-Z0-9=_.\-]{66}`)},
	// JWT Tokens
	{re: regexp.MustCompile(`\beyJ[a-zA-Z0-9_-]*\.eyJ[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*\b`)},
	// Private Key Headers
	{re: regexp.MustCompile(`-----BEGIN\s+(RSA|DSA|EC|OPENSSH|PGP)\s+PRIVATE\s+KEY-----`)},
	// AWS Secret Access Key (contextual). The value alone is 40 base64 chars
	// with no fixed prefix, which collides with file paths, hashes, and other
	// 40-char strings — even gitleaks ships no standalone rule for it, and
	// trufflehog / detect-secrets only confirm matches by calling AWS STS. We
	// can't validate live, so we mask only when the value is preceded by an
	// obvious AWS_SECRET_*_KEY-style identifier. Covers the env-var, AWS
	// credentials file, and JSON/YAML forms; uses [ \t]* (not \s*) to keep the
	// match single-line, per TestPatternsAreSingleLine.
	{
		re:        regexp.MustCompile(`(?i)\baws[_.\-]?secret[_.\-]?(?:access[_.\-]?)?key\b["']?[ \t]*[:=][ \t]*["']?([a-zA-Z0-9+/]{40})`),
		maskGroup: 1,
	},
}

// Mask replaces sensitive patterns in content with asterisks of the same length.
func Mask(content string) string {
	result := content
	for _, p := range patterns {
		result = p.apply(result)
	}
	return result
}

func (p pattern) apply(s string) string {
	if p.maskGroup == 0 {
		return p.re.ReplaceAllStringFunc(s, func(match string) string {
			if p.validate != nil && !p.validate(match) {
				return match
			}
			return strings.Repeat("*", len(match))
		})
	}

	matches := p.re.FindAllStringSubmatchIndex(s, -1)
	if len(matches) == 0 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	last := 0
	gs, ge := p.maskGroup*2, p.maskGroup*2+1
	for _, m := range matches {
		if ge >= len(m) || m[gs] < 0 {
			continue
		}
		if p.validate != nil && !p.validate(s[m[gs]:m[ge]]) {
			continue
		}
		b.WriteString(s[last:m[gs]])
		b.WriteString(strings.Repeat("*", m[ge]-m[gs]))
		last = m[ge]
	}
	b.WriteString(s[last:])
	return b.String()
}
