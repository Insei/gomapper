package gomapper

import (
	"fmt"
	"reflect"
)

var routes = map[reflect.Type]map[reflect.Type]func(source interface{}, dest interface{}) error{}

func addSliceRoute[TSliceSource any, TSliceDest any](sliceMapFunc func(sourceSlice TSliceSource, destSlice TSliceDest) error) {
	funcConverted1 := func(source any, dest any) error {
		s := *new(TSliceSource)
		d := *new(TSliceDest)
		fmt.Print(s, d)
		return sliceMapFunc(source.(TSliceSource), dest.(TSliceDest))
	}
	sourceSlice := *new(TSliceSource)
	destSlice := *new(TSliceDest)
	var route map[reflect.Type]func(source interface{}, dest interface{}) error
	var ok bool
	if route, ok = routes[reflect.TypeOf(sourceSlice)]; !ok {
		route = map[reflect.Type]func(source interface{}, dest interface{}) error{}
		route[reflect.TypeOf(destSlice)] = funcConverted1
		routes[reflect.TypeOf(sourceSlice)] = route
	}
}

func addSliceRoutes[TSource, TDest any]() {
	//source slice is a value, dest slice is a pointer
	addSliceRoute[[]TSource, *[]TDest](func(sourceSlice []TSource, pointerDestSlice *[]TDest) error {
		for _, source := range sourceSlice {
			dest, err := MapTo[TDest](source)
			if err != nil {
				return err
			}
			*pointerDestSlice = append(*pointerDestSlice, dest)
		}
		return nil
	})
	//source slice is a value, dest slice is a pointer
	addSliceRoute[*[]TSource, *[]TDest](func(sourceSlice *[]TSource, pointerDestSlice *[]TDest) error {
		for _, source := range *sourceSlice {
			dest, err := MapTo[TDest](source)
			if err != nil {
				return err
			}
			*pointerDestSlice = append(*pointerDestSlice, dest)
		}
		return nil
	})
}

func AddRoute[TSource, TDest any](mapFunc func(source TSource, dest *TDest) error) error {
	source := *new(TSource)
	dest := *new(TDest)
	destValueOf := reflect.ValueOf(dest)
	if destValueOf.Kind() == reflect.Ptr {
		return fmt.Errorf("destination type can't be reference type")
	}
	sourceValueOf := reflect.ValueOf(source)
	if sourceValueOf.Kind() == reflect.Ptr {
		return fmt.Errorf("source type can't be reference type")
	}
	var route map[reflect.Type]func(source interface{}, dest interface{}) error
	route, ok := routes[reflect.TypeOf(source)]
	if !ok {
		route = map[reflect.Type]func(source interface{}, dest interface{}) error{}
		routes[reflect.TypeOf(source)] = route
	}
	funcConverted := func(source any, dest any) error {
		sourceForMap := source
		sourceValueOf := reflect.ValueOf(source)
		if sourceValueOf.Kind() == reflect.Ptr {
			sourceForMap = sourceValueOf.Elem().Interface()
		}
		return mapFunc(sourceForMap.(TSource), dest.(*TDest))
	}
	route[reflect.TypeOf(&dest)] = funcConverted
	// source is value, dest is ptr - its important
	addSliceRoutes[TSource, TDest]()
	return nil
}

// Map source to dest
func Map(source interface{}, dest interface{}) error {
	if source == nil || dest == nil {
		return fmt.Errorf("")
	}
	destValueOf := reflect.ValueOf(dest)
	if destValueOf.Kind() != reflect.Ptr {
		return fmt.Errorf("")
	}

	var sourceToMap any
	sourceToMap = source
	sourceValueOf := reflect.ValueOf(source)
	if sourceValueOf.Kind() == reflect.Ptr {
		sourceToMap = sourceValueOf.Elem().Interface()
	}
	route, ok := routes[reflect.TypeOf(sourceToMap)]
	if !ok {
		return fmt.Errorf("route not found")
	}
	mapFunc, ok := route[reflect.TypeOf(dest)]
	if !ok {
		return fmt.Errorf("route not found")
	}
	return mapFunc(source, dest)
}

// MapTo Map source to the new dest object
func MapTo[TDest interface{}](source interface{}) (TDest, error) {
	dest := new(TDest)
	if source == nil {
		return *dest, nil
	}
	var sourceToMap any
	sourceToMap = source
	sourceValueOf := reflect.ValueOf(source)
	if sourceValueOf.Kind() == reflect.Ptr {
		sourceToMap = sourceValueOf.Elem().Interface()
	}
	route, ok := routes[reflect.TypeOf(sourceToMap)]
	if !ok {
		return *dest, fmt.Errorf("route does not exist in routes map")
	}
	mapFunc, ok := route[reflect.TypeOf(dest)]
	if !ok {
		return *dest, fmt.Errorf("route does not exist in routes map")
	}
	err := mapFunc(sourceToMap, dest)
	return *dest, err
}
