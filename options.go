package gomapper

type withFuncOption[TSource, TDest any] struct {
	fn func(TSource, *TDest)
}

type options struct {
	Fns []any
}

type Option interface {
	apply(*options)
}

func (a withFuncOption[TSource, TDest]) apply(opts *options) {
	opts.Fns = append(opts.Fns, a.fn)
}

func WithFunc[TSource, TDest any](fn func(TSource, *TDest)) Option {
	return &withFuncOption[TSource, TDest]{fn: fn}
}
