package prettier

var Registry = map[string]Prettier{}

func Register(name string, p Prettier) Prettier {
	Registry[name] = p
	return p
}

func Get(name string) (Prettier, bool) {
	p, ok := Registry[name]
	return p, ok
}
