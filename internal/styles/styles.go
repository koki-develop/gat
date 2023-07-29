package styles

import (
	"embed"
	"io/fs"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
)

//go:embed *.xml
var embedded embed.FS

func init() {
	files, err := fs.ReadDir(embedded, ".")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		r, err := embedded.Open(file.Name())
		if err != nil {
			panic(err)
		}
		style, err := chroma.NewXMLStyle(r)
		if err != nil {
			panic(err)
		}
		styles.Register(style)
		_ = r.Close()
	}
}

func Get(name string) (*chroma.Style, bool) {
	s, ok := styles.Registry[name]
	return s, ok
}

func List() []string {
	return styles.Names()
}
