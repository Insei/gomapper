package gomapper

import (
	"fmt"
	"reflect"
)

func validateSource(source any) error {
	sourceTypeName := getTypeName(source)
	if source == nil {
		return fmt.Errorf("source value can't be nil, source type: %s", sourceTypeName)
	}
	typeOf := reflect.TypeOf(source)
	if typeOf.Kind() == reflect.Ptr && typeOf.Elem().Kind() == reflect.Ptr {
		return fmt.Errorf("source can have a pointer type, but not a pointer to pointer, source type: %s", sourceTypeName)
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

func getTypeNameRecursive(target reflect.Type, typeName string) string {
	if target.Kind() == reflect.Ptr || target.Kind() == reflect.Slice {
		newTarget := target.Elem()
		newTypeName := "*" + typeName
		if target.Kind() == reflect.Slice {
			newTypeName = "[]" + typeName
		}
		return getTypeNameRecursive(newTarget, newTypeName)
	}
	return fmt.Sprintf("%s%s.%s", typeName, target.PkgPath(), target.Name())
}

func getTypeName(target any) string {
	tt := reflect.TypeOf(target)
	if target != nil || tt != nil {
		return getTypeNameRecursive(tt, "")
	}
	return "undefined"
}

func validateDest(dest any) error {
	dValueOf := reflect.ValueOf(dest)
	dTypeName := getTypeName(dest)
	if dValueOf.Kind() != reflect.Ptr {
		return fmt.Errorf("destenation value should have a pointer type, but has %s type", dTypeName)
	}
	if dest == nil || dValueOf.IsNil() {
		return fmt.Errorf("destenation value can't be nil, destenation type: %s", dTypeName)
	}
	if dValueOf.Kind() == reflect.Ptr && dValueOf.Elem().Kind() == reflect.Ptr {
		return fmt.Errorf("destenation value should have a pointer type, not a pointer to pointer, but has %s type", dTypeName)
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
		return fmt.Errorf("route not found for type %s to type %s",
			getTypeName(sourceForMap), getTypeName(dest))
	}
	mapFunc, ok := route[reflect.TypeOf(dest)]
	if !ok {
		return fmt.Errorf("route not found for type %s to type %s",
			getTypeName(sourceForMap), getTypeName(dest))
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
