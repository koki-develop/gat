package masker

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMask(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "AWS Access Key",
			input: "aws_access_key_id = AKIAIOSFODNN7EXAMPLE",
			want:  "aws_access_key_id = " + strings.Repeat("*", 20),
		},
		{
			name:  "AWS Access Key (temporary/SSO)",
			input: "aws_access_key_id = ASIAISEXAMPLEKEY1234",
			want:  "aws_access_key_id = " + strings.Repeat("*", 20),
		},
		{
			name:  "AWS Secret Access Key (credentials file format)",
			input: "aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "aws_secret_access_key = " + strings.Repeat("*", 40),
		},
		{
			name:  "AWS Secret Access Key (env var format)",
			input: "AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "AWS_SECRET_ACCESS_KEY=" + strings.Repeat("*", 40),
		},
		{
			name:  "AWS Secret Access Key (JSON format)",
			input: `"aws_secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"`,
			want:  `"aws_secret_access_key": "` + strings.Repeat("*", 40) + `"`,
		},
		{
			name:  "AWS Secret Access Key (short AWS_SECRET_KEY alias)",
			input: "AWS_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "AWS_SECRET_KEY=" + strings.Repeat("*", 40),
		},
		{
			// Exercises apply()'s multi-match submatch path: two contextual
			// secrets in one input. Verifies the inter-match text and tail
			// are preserved verbatim when output is rebuilt index-by-index
			// via FindAllStringSubmatchIndex.
			name:  "AWS Secret Access Key (multiple secrets in one input)",
			input: "aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY\n# next entry\nAWS_SECRET_KEY=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			want:  "aws_secret_access_key = " + strings.Repeat("*", 40) + "\n# next entry\nAWS_SECRET_KEY=" + strings.Repeat("*", 40),
		},
		{
			// The character class [a-zA-Z0-9+/] includes `+` alongside `/`.
			// Locks in `+` so a future tightening (e.g. to \w{40}) is caught.
			name:  "AWS Secret Access Key (value containing +)",
			input: "AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI+K7MDENG+bPxRfiCYEXAMPLEKEY",
			want:  "AWS_SECRET_ACCESS_KEY=" + strings.Repeat("*", 40),
		},
		{
			// `[_.\-]?` allows `.` between identifier segments (Java
			// properties / system-property form).
			name:  "AWS Secret Access Key (dot-separated identifier)",
			input: "aws.secret.access.key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "aws.secret.access.key=" + strings.Repeat("*", 40),
		},
		{
			// `[_.\-]?` allows `-` between identifier segments (kebab-case).
			name:  "AWS Secret Access Key (kebab-case identifier)",
			input: "aws-secret-access-key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "aws-secret-access-key=" + strings.Repeat("*", 40),
		},
		{
			// `[ \t]*` (not `\s*`, to preserve the single-line invariant)
			// allows tabs around the separator.
			name:  "AWS Secret Access Key (tab around separator)",
			input: "AWS_SECRET_ACCESS_KEY\t=\twJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "AWS_SECRET_ACCESS_KEY\t=\t" + strings.Repeat("*", 40),
		},
		{
			// `(?i)` makes the identifier match case-insensitively.
			name:  "AWS Secret Access Key (mixed-case identifier)",
			input: "Aws_Secret_Access_Key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "Aws_Secret_Access_Key=" + strings.Repeat("*", 40),
		},
		{
			// Each `[_.\-]?` is individually optional, so a concatenated
			// identifier with no internal separators also matches. Locks in
			// the behavior — flip `[_.\-]?` to `[_.\-]` to reject it.
			name:  "AWS Secret Access Key (no internal separators)",
			input: "awssecretkey=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "awssecretkey=" + strings.Repeat("*", 40),
		},
		{
			name:  "GitHub Personal Access Token",
			input: "token: ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			want:  "token: " + strings.Repeat("*", 40),
		},
		{
			name:  "GitHub OAuth Token",
			input: "GITHUB_TOKEN=gho_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			want:  "GITHUB_TOKEN=" + strings.Repeat("*", 40),
		},
		{
			name:  "GitHub Fine-grained Personal Access Token",
			input: "token: github_pat_11ABCDE0Y0abcdefghijkl_mnopqrstuvwxyzABCDEFGHIJKLMNOPqrstuvwxyz0123456789ABCDEFGHI",
			want:  "token: " + strings.Repeat("*", 93),
		},
		{
			name:  "GitHub App installation token (stateless, JWT format)",
			input: "token: ghs_eyJhbGciOiJIUzI1NiJ9.eyJpbnN0YWxsYXRpb25faWQiOjEyMzQ1Njc4OTB9.EXAMPLE_SIGNATURE_DO_NOT_USE_THIS_012345678",
			want:  "token: " + strings.Repeat("*", 109),
		},
		{
			name:  "GitLab Personal Access Token",
			input: "GITLAB_TOKEN=glpat-xxxxxxxxxxxxxxxxxxxx",
			want:  "GITLAB_TOKEN=" + strings.Repeat("*", 26),
		},
		{
			name:  "Slack Bot Token",
			input: "SLACK_TOKEN=xoxb-123456789-abcdefgh",
			want:  "SLACK_TOKEN=" + strings.Repeat("*", 23),
		},
		{
			name:  "Slack App-level Token",
			input: "SLACK_APP_TOKEN=xapp-1-A0123ABCDEF-1234567890123-abcdef0123456789abcdef0123456789",
			want:  "SLACK_APP_TOKEN=" + strings.Repeat("*", 65),
		},
		{
			name:  "Anthropic API Key",
			input: "ANTHROPIC_API_KEY=sk-ant-api03-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			want:  "ANTHROPIC_API_KEY=" + strings.Repeat("*", 57),
		},
		{
			name:  "OpenAI API Key (legacy)",
			input: "OPENAI_API_KEY=sk-1234567890_abcdef-1234567890_abcdef-1234567890",
			want:  "OPENAI_API_KEY=" + strings.Repeat("*", 49),
		},
		{
			name:  "OpenAI API Key (project)",
			input: "OPENAI_API_KEY=sk-proj-abcd_1234-efgh_5678-ijkl_9012-mnop",
			want:  "OPENAI_API_KEY=" + strings.Repeat("*", 42),
		},
		{
			name:  "Supabase Secret Key",
			input: "SUPABASE_KEY=sb_secret_1234567890abcdef1234567890abcdef",
			want:  "SUPABASE_KEY=" + strings.Repeat("*", 42),
		},
		{
			name:  "npm Access Token",
			input: "NPM_TOKEN=npm_abcdefghijklmnopqrstuvwxyz0123456789",
			want:  "NPM_TOKEN=" + strings.Repeat("*", 40),
		},
		{
			name:  "PyPI API Token",
			input: "PYPI_TOKEN=pypi-AgEIcHlwaS5vcmcAbCdEfGhIjKlMnOpQrStUvWxYz0123456789AbCdEfGhIjKlMn",
			want:  "PYPI_TOKEN=" + strings.Repeat("*", 70),
		},
		{
			name:  "RubyGems API Key",
			input: "RUBYGEMS_API_KEY=rubygems_deadbeefdeadbeefdeadbeefdeadbeef",
			want:  "RUBYGEMS_API_KEY=" + strings.Repeat("*", 41),
		},
		{
			name:  "Google API Key",
			input: "GOOGLE_API_KEY=AIzaSyEXAMPLE_GOOGLE_API_KEY_0123456789",
			want:  "GOOGLE_API_KEY=" + strings.Repeat("*", 39),
		},
		{
			name:  "Firebase API Key",
			input: "FIREBASE_API_KEY=AIzaSyEXAMPLE_FIREBASE_API_KEY_01234567",
			want:  "FIREBASE_API_KEY=" + strings.Repeat("*", 39),
		},
		{
			name:  "Stripe Secret Key (live)",
			input: "STRIPE_SECRET_KEY=sk_live_EXAMPLESTRIPESECRETKEYDONOTUSE012345678901234567890123456789012345678901234567890123456789012345678",
			want:  "STRIPE_SECRET_KEY=" + strings.Repeat("*", 107),
		},
		{
			name:  "Stripe Restricted Key (test)",
			input: "STRIPE_RK=rk_test_EXAMPLEDUMMYKEY123456789",
			want:  "STRIPE_RK=" + strings.Repeat("*", 32),
		},
		{
			name:  "Stripe Organization Key",
			input: "STRIPE_ORG_KEY=sk_org_EXAMPLEDUMMYKEY123456789",
			want:  "STRIPE_ORG_KEY=" + strings.Repeat("*", 31),
		},
		{
			name:  "SendGrid API Key",
			input: "SENDGRID_API_KEY=SG.EXAMPLE_SENDGRID_TOKEN.EXAMPLE_SENDGRID_SIGNATURE_DO_NOT_USE_00000",
			want:  "SENDGRID_API_KEY=" + strings.Repeat("*", 69),
		},
		{
			name:  "JWT Token",
			input: "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			want:  "Authorization: Bearer " + strings.Repeat("*", 108),
		},
		{
			name:  "RSA Private Key Header",
			input: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...",
			want:  strings.Repeat("*", 31) + "\nMIIEpAIBAAKCAQEA...",
		},
		{
			name:  "No sensitive data",
			input: "const message = 'Hello World'",
			want:  "const message = 'Hello World'",
		},
		{
			name:  "Multiple secrets",
			input: "GITHUB=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\nAWS=AKIAIOSFODNN7EXAMPLE",
			want:  "GITHUB=" + strings.Repeat("*", 40) + "\nAWS=" + strings.Repeat("*", 20),
		},
		{
			name:  "Empty string",
			input: "",
			want:  "",
		},
		{
			name:  "Preserves length",
			input: "KEY=AKIAIOSFODNN7EXAMPLE",
			want:  "KEY=" + strings.Repeat("*", 20),
		},
		{
			name:  "Long alphanumeric string should not be partially masked",
			input: strings.Repeat("A", 76),
			want:  strings.Repeat("A", 76),
		},
		{
			// Regression: a 40-char filesystem path segment used to false-positive
			// against the old standalone AWS Secret Access Key pattern. With the
			// contextual pattern, no AWS keyword precedes the value, so it stays.
			name:  "40-char path segment without AWS keyword is not masked",
			input: "/Users/koki/.opensrc/repos/github.com/clauderic/dnd-kit/0.4.0/packages/react/src/sortable/useSortable.ts",
			want:  "/Users/koki/.opensrc/repos/github.com/clauderic/dnd-kit/0.4.0/packages/react/src/sortable/useSortable.ts",
		},
		{
			name:  "Git full SHA without AWS keyword is not masked",
			input: "commit da39a3ee5e6b4b0d3255bfef95601890afd80709",
			want:  "commit da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
		{
			// A bare 40-char base64 value with no preceding AWS_SECRET_*_KEY
			// identifier is intentionally not masked — see gitleaks' decision
			// to drop the standalone rule for the same reason.
			name:  "Bare 40-char base64 without AWS keyword is not masked",
			input: "value: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "value: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mask(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestPatternsAreSingleLine guards an invariant the streaming passthrough path
// in internal/gat depends on: every pattern is confined to a single line,
// except the private-key header. That path masks input line by line, so a
// newly added pattern that can span newlines would silently fail to mask there.
// If this test fails, reconsider the per-line masking in internal/gat before
// adding the pattern.
func TestPatternsAreSingleLine(t *testing.T) {
	// spansNewline reports whether a pattern's source contains a construct that
	// can match across a newline in Go's RE2: \s matches \n, (?s) makes . match
	// \n, and a literal \n obviously spans lines.
	spansNewline := func(src string) bool {
		return strings.Contains(src, `\s`) ||
			strings.Contains(src, `(?s)`) ||
			strings.Contains(src, `\n`)
	}

	for _, p := range patterns {
		src := p.re.String()
		if strings.Contains(src, "PRIVATE") {
			// The private-key header is the sole, documented exception.
			assert.True(t, spansNewline(src),
				"expected private-key pattern to span newlines (\\s+): %q", src)
			continue
		}
		assert.False(t, spansNewline(src),
			"pattern %q can match across newlines; per-line masking in internal/gat "+
				"(streaming passthrough) would silently miss it", src)
	}
}
