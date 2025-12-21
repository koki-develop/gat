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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mask(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
