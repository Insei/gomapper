package gomapper

import (
	"github.com/insei/fmap/v3"
)

type withFuncOption[TSource, TDest any] struct {
	fn func(TSource, *TDest)
}

type withFieldSkip[TDest any] struct {
	field fmap.Field
}

type options struct {
	Fns      []any
	Excluded []fmap.Field
}

type Option interface {
	apply(*options)
}

func (a withFuncOption[TSource, TDest]) apply(opts *options) {
	opts.Fns = append(opts.Fns, a.fn)
}

func (a withFieldSkip[TDest]) apply(opts *options) {
	if opts.Excluded == nil {
		opts.Excluded = make([]fmap.Field, 0)
	}
	opts.Excluded = append(opts.Excluded, a.field)
}

func WithFunc[TSource, TDest any](fn func(TSource, *TDest)) Option {
	return &withFuncOption[TSource, TDest]{fn: fn}
}

func WithFieldSkip[TSource any](fn func(*TSource) any) Option {
	source := new(TSource)

	storage, err := fmap.GetFrom(source)
	if err != nil {
		panic(err)
	}

	fieldPtr := fn(source)

	field, err := storage.GetFieldByPtr(source, fieldPtr)
	if err != nil {
		panic(err)
	}

	return &withFieldSkip[TSource]{field: field}
}
