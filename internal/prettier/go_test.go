package prettier

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoPrettier_Pretty(t *testing.T) {
	tests := []struct {
		src     string
		want    string
		wantErr bool
	}{
		{
			src: `package main
import "fmt"
func main() {
fmt.Println("Hello, world!")
}`,
			want: `package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}
`,
			wantErr: false,
		},
		{
			src:     "package main\nfunc main() {",
			want:    "",
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			g := NewGoPrettier()
			got, err := g.Pretty(tt.src)

			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
