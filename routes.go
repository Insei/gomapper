package gomapper

import (
	"fmt"
	"reflect"
)

var routes = map[reflect.Type]map[reflect.Type]func(source interface{}, dest interface{}) error{}

func addSliceRoute[TSliceSource any, TSliceDest any](sliceMapFunc func(sourceSlice TSliceSource, destSlice TSliceDest) error) {
	funcConverted := func(source any, dest any) error {
		return sliceMapFunc(source.(TSliceSource), dest.(TSliceDest))
	}
	sourceSlice := *new(TSliceSource)
	destSlice := *new(TSliceDest)
	var route map[reflect.Type]func(source interface{}, dest interface{}) error
	var ok bool
	if route, ok = routes[reflect.TypeOf(sourceSlice)]; !ok {
		route = map[reflect.Type]func(source interface{}, dest interface{}) error{}
		routes[reflect.TypeOf(sourceSlice)] = route
	}
	route[reflect.TypeOf(destSlice)] = funcConverted
}

func addSliceRoutes[TSource, TDest any]() {
	//source slice is a value, dest slice is a pointer
	addSliceRoute(func(sourceSlice []TSource, pointerDestSlice *[]TDest) error {
		if len(sourceSlice) == 0 {
			*pointerDestSlice = make([]TDest, 0)
		}
		for _, source := range sourceSlice {
			dest, err := MapTo[TDest](source)
			if err != nil {
				return err
			}
			*pointerDestSlice = append(*pointerDestSlice, dest)
		}
		return nil
	})
	//source slice is a value, dest slice is a pointer with pointer elements
	addSliceRoute(func(sourceSlice []TSource, pointerDestSlice *[]*TDest) error {
		if len(sourceSlice) == 0 {
			*pointerDestSlice = make([]*TDest, 0)
		}
		for _, source := range sourceSlice {
			dest, err := MapTo[TDest](source)
			if err != nil {
				return err
			}
			*pointerDestSlice = append(*pointerDestSlice, &dest)
		}
		return nil
	})
	//source slice is a value, dest slice is a pointer
	addSliceRoute(func(sourceSlice []*TSource, pointerDestSlice *[]TDest) error {
		if len(sourceSlice) == 0 {
			*pointerDestSlice = make([]TDest, 0)
		}
		for _, source := range sourceSlice {
			dest, err := MapTo[TDest](source)
			if err != nil {
				return err
			}
			*pointerDestSlice = append(*pointerDestSlice, dest)
		}
		return nil
	})
	addSliceRoute(func(sourceSlice []*TSource, pointerDestSlice *[]*TDest) error {
		if len(sourceSlice) == 0 {
			*pointerDestSlice = make([]*TDest, 0)
		}
		for _, source := range sourceSlice {
			dest, err := MapTo[TDest](source)
			if err != nil {
				return err
			}
			*pointerDestSlice = append(*pointerDestSlice, &dest)
		}
		return nil
	})
}

func addRoute[TSource, TDest any | []any](mapFunc func(source TSource, dest *TDest) error) error {
	source := *new(TSource)
	dest := *new(TDest)

	sourceValueOf := reflect.ValueOf(source)
	if sourceValueOf.Kind() == reflect.Ptr {
		return fmt.Errorf("source type can't be reference type, route: %s -> %s", getTypeName(source), getTypeName(dest))
	}
	var route map[reflect.Type]func(source interface{}, dest interface{}) error
	route, ok := routes[reflect.TypeOf(source)]
	if !ok {
		route = map[reflect.Type]func(source interface{}, dest interface{}) error{}
		routes[reflect.TypeOf(source)] = route
	}
	funcConverted := func(source any, dest any) error {
		sourceValueOf := reflect.ValueOf(source)
		for sourceValueOf.Kind() == reflect.Ptr {
			if sourceValueOf.IsNil() {
				return nil
			}
			sourceValueOf = sourceValueOf.Elem()
		}
		return mapFunc(sourceValueOf.Interface().(TSource), dest.(*TDest))
	}
	route[reflect.TypeOf(&dest)] = funcConverted
	// source is value, dest is ptr - its important
	addSliceRoutes[TSource, TDest]()
	return nil
}

func AddRoute[TSource, TDest any | []any](mapFunc func(source TSource, dest *TDest) error) error {
	return addRoute[TSource, TDest](mapFunc)
}
