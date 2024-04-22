package gomapper

import (
	"github.com/insei/gomapper/fields"
	"reflect"
	"strings"
)

var manualFieldRoutes = map[reflect.Type]map[string]string{}

func AutoRoute[TSource, TDest any | []any](options ...AutoMapperOption) error {
	s := new(TSource)
	d := new(TDest)
	sourceFields := fields.GetFrom(s)
	destFields := fields.GetFrom(d)
	sourceType := reflect.TypeOf(s)

	parseOptions(options, sourceType)

	mapFunc := func(source TSource, dest *TDest) error {
		for key, sourceFld := range sourceFields {
			destFld, ok := destFields[getDestFieldName(sourceType, key)]
			if !ok || strings.Contains(key, ".") {
				continue
			}
			if err := setFieldRecursive(sourceFld, destFld, source, dest); err != nil {
				return err
			}
		}
		return nil
	}
	return AddRoute[TSource, TDest](mapFunc)
}

func parseOptions(options []AutoMapperOption, sourceType reflect.Type) {
	for _, option := range options {
		switch autoMapperOption := option.(type) {
		case fieldPathOption:
			if manualFieldRoutes[sourceType] == nil {
				manualFieldRoutes[sourceType] = map[string]string{}
			}
			manualFieldRoutes[sourceType][autoMapperOption.source] = autoMapperOption.dest
		}
	}
}

func setFieldRecursive(sourceFld, destFld fields.Field, source, dest any) error {
	if r, ok := getRouteIfExists(sourceFld, destFld); ok {
		return r(sourceFld.Get(source), destFld.Get(dest))
	}

	if sourceFld.Type.Kind() != reflect.Struct {
		sourceVal := sourceFld.Get(source)
		if sourceVal != nil {
			destFld.Set(dest, sourceVal)
		}
		return nil
	}

	sourceStructField := sourceFld.Get(source)
	sourceFields := fields.GetFrom(sourceStructField)
	destStructField := destFld.Get(dest)
	destFields := fields.GetFrom(destStructField)

	for fieldName, sField := range sourceFields {
		dField, ok := destFields[getDestFieldName(sField.Type, fieldName)]
		if !ok || strings.Contains(fieldName, ".") {
			continue
		}
		err := setFieldRecursive(sField, dField, sourceStructField, destStructField)
		if err != nil {
			return err
		}
	}
	return nil
}

func getRouteIfExists(sourceFld, destFld fields.Field) (func(source interface{}, dest interface{}) error, bool) {
	destType := destFld.Type
	sourceType := sourceFld.Type
	for sourceType.Kind() == reflect.Ptr {
		sourceType = sourceType.Elem()
	}
	destType = reflect.PointerTo(destType)
	r, ok := routes[sourceType][destType]
	return r, ok
}

func getDestFieldName(sourceFieldType reflect.Type, sourceFieldName string) string {
	if destFieldName, ok := manualFieldRoutes[sourceFieldType][sourceFieldName]; ok {
		return destFieldName
	}
	return sourceFieldName
}
