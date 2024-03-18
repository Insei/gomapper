package gomapper

type AutoMapperOption interface {
}

type fieldPathOption struct {
	AutoMapperOption
	source string
	dest   string
}

func WithFieldRoute(source, dest string) AutoMapperOption {
	return fieldPathOption{
		source: source,
		dest:   dest,
	}
}
