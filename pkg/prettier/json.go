package prettier

import (
	"bytes"
	"encoding/json"
)

var (
	JSON = Register("JSON", NewJSONPrettier())
)

type JSONPrettier struct{}

func NewJSONPrettier() *JSONPrettier {
	return &JSONPrettier{}
}

func (*JSONPrettier) Pretty(j string) (string, error) {
	b := new(bytes.Buffer)
	if err := json.Indent(b, []byte(j), "", "  "); err != nil {
		return "", err
	}
	return b.String(), nil
}
