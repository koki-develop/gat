package prettier

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMLPrettier_Pretty(t *testing.T) {
	tests := []struct {
		src     string
		want    string
		wantErr bool
	}{
		{
			src: `<?xml version="1.0" encoding="UTF-8"?><root><child>text</child></root>`,
			want: `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <child>text</child>
</root>`,
			wantErr: false,
		},
		{
			src:     "<INVALID_XML>",
			want:    "",
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			x := NewXMLPrettier()
			got, err := x.Pretty(tt.src)

			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
