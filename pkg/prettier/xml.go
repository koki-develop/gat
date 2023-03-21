package prettier

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

var (
	XML = Register("XML", NewXMLPrettier())
)

type XMLPrettier struct{}

func NewXMLPrettier() *XMLPrettier {
	return &XMLPrettier{}
}

func (*XMLPrettier) Pretty(input string) (string, error) {
	var b bytes.Buffer
	dec := xml.NewDecoder(strings.NewReader(input))
	enc := xml.NewEncoder(&b)
	enc.Indent("", "  ")

	for {
		token, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		if err := enc.EncodeToken(token); err != nil {
			return "", err
		}

		if declaration, ok := token.(xml.ProcInst); ok && declaration.Target == "xml" {
			if err := enc.Flush(); err != nil {
				return "", err
			}
			if _, err := b.WriteString("\n"); err != nil {
				return "", err
			}
		}
	}

	if err := enc.Flush(); err != nil {
		return "", err
	}

	return b.String(), nil
}
