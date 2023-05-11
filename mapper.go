package gomapper

import (
	"fmt"
	"reflect"
)

var routes = map[reflect.Type]map[reflect.Type]func(source interface{}, dest interface{}) error{}

func AddRoute[TSource, TDest interface{}](mapFunc func(source TSource, dest TDest) error) error {
	source := *new(TSource)
	dest := *new(TDest)
	destValueOf := reflect.ValueOf(dest)
	if destValueOf.Kind() != reflect.Ptr {
		return fmt.Errorf("destination object must be of reference type")
	}
	sourceValueOf := reflect.ValueOf(source)
	if sourceValueOf.Kind() == reflect.Ptr {
		return fmt.Errorf("source object must not be of reference type")
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
		return mapFunc(sourceForMap.(TSource), dest.(TDest))
	}
	route[reflect.TypeOf(dest)] = funcConverted
	return nil
}

//Map source to dest
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
		return *dest, nil
	}
	mapFunc, ok := route[reflect.TypeOf(dest)]
	if !ok {
		return *dest, nil
	}
	err := mapFunc(source, dest)
	return *dest, err
}
