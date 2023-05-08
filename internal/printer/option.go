package printer

type option struct {
	filename *string
}

type Option func(*option)

func WithFilename(filename string) Option {
	return func(o *option) {
		if filename != "" {
			o.filename = &filename
		}
	}
}
