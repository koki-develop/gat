package prettier

type FallbackPrettier struct{}

func NewFallbackPrettier() *FallbackPrettier {
	return &FallbackPrettier{}
}

func (*FallbackPrettier) Pretty(content string) (string, error) {
	return content, nil
}
