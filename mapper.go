package gomapper

import (
	"fmt"
	"reflect"
)

const RouteNotFoundError = "route does not exist in route map"

func validateSource(source any) error {
	sValueOf := reflect.ValueOf(source)
	if source == nil || (sValueOf.Kind() == reflect.Ptr && sValueOf.IsNil()) {
		return fmt.Errorf("source value cant't be nil")
	}
	if sValueOf.Kind() == reflect.Ptr && sValueOf.Elem().Kind() == reflect.Ptr {
		return fmt.Errorf("source can have a pointer type, but not a pointer to pointer")
	}
	return nil
}

func prepareSource(source any) any {
	sourceToMap := source
	sourceValueOf := reflect.ValueOf(source)
	if sourceValueOf.Kind() == reflect.Ptr {
		sourceToMap = sourceValueOf.Elem().Interface()
	}
	return sourceToMap
}

func validateDest(dest any) error {
	dValueOf := reflect.ValueOf(dest)
	if dest == nil || dValueOf.IsNil() {
		return fmt.Errorf("destenation value can't be nil")
	}
	if dValueOf.Kind() != reflect.Ptr {
		return fmt.Errorf("destenation value should have a pointer type")
	}
	if dValueOf.Elem().Kind() == reflect.Ptr {
		return fmt.Errorf("destenation value should have a pointer type, not a pointer to pointer")
	}
	return nil
}

// Map source to dest
func Map(source interface{}, dest interface{}) error {
	err := validateSource(source)
	if err != nil {
		return err
	}
	sourceForMap := prepareSource(source)
	err = validateDest(dest)
	if err != nil {
		return err
	}
	route, ok := routes[reflect.TypeOf(sourceForMap)]
	if !ok {
		return fmt.Errorf(RouteNotFoundError)
	}
	mapFunc, ok := route[reflect.TypeOf(dest)]
	if !ok {
		return fmt.Errorf(RouteNotFoundError)
	}
	return mapFunc(sourceForMap, dest)
}

// MapTo Map source to the new dest object
func MapTo[TDest interface{}](source interface{}) (TDest, error) {
	dest := new(TDest)
	err := Map(source, dest)
	if err != nil {
		return *dest, err
	}
	return *dest, nil
}
