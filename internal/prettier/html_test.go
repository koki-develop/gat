package prettier

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLPrettier_Pretty(t *testing.T) {
	tests := []struct {
		src     string
		want    string
		wantErr bool
	}{
		{
			src: `<!DOCTYPE html><html><head><title>title</title></head><body><p>paragraph</p></body></html>`,
			want: `<!DOCTYPE html>
<html>
  <head>
    <title>
      title
    </title>
  </head>
  <body>
    <p>
      paragraph
    </p>
  </body>
</html>`,
			wantErr: false,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			p := NewHTMLPrettier()
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
