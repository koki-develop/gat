package prettier

type Prettier interface {
	Pretty(content string) (string, error)
}
