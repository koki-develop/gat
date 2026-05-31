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
			name:  "AWS Secret Access Key",
			input: "aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "aws_secret_access_key = " + strings.Repeat("*", 40),
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
			name:  "AWS Secret Access Key (real value with mixed case) is masked",
			input: "aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			want:  "aws_secret_access_key = " + strings.Repeat("*", 40),
		},
		{
			name:  "Git full SHA (lower-case hex) is not mistaken for an AWS secret key",
			input: "commit da39a3ee5e6b4b0d3255bfef95601890afd80709",
			want:  "commit da39a3ee5e6b4b0d3255bfef95601890afd80709",
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
