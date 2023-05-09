package prettier

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLPrettier_Pretty(t *testing.T) {
	tests := []struct {
		src     string
		want    string
		wantErr bool
	}{
		{
			src: `foo:            bar
baz:  qux    
`,
			want: `foo: bar
baz: qux
`,
			wantErr: false,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			p := NewYAMLPrettier()
			got, err := p.Pretty(tt.src)

			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
