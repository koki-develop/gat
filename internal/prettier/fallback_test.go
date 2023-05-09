package prettier

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFallbackPrettier_Pretty(t *testing.T) {
	tests := []struct {
		src     string
		want    string
		wantErr bool
	}{
		{
			src:     "CONTENT",
			want:    "CONTENT",
			wantErr: false,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			f := NewFallbackPrettier()
			got, err := f.Pretty(tt.src)

			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
