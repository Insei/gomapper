package gomapper

import (
	"reflect"

	"github.com/insei/fmap/v3"
)

var (
	manualFieldRoutes = map[reflect.Type]map[string]string{}
)

func AutoRoute[TSource, TDest any | []any](options ...AutoMapperOption) error {
	s := new(TSource)
	d := new(TDest)
	sourceStorage, _ := fmap.GetFrom(s)
	destStorage, _ := fmap.GetFrom(d)
	sourceType := reflect.TypeOf(s)

	parseOptions(options, sourceType)

	mapFunc := func(source TSource, dest *TDest) error {
		for _, sourcePath := range sourceStorage.GetAllPaths() {
			destFld, ok := destStorage.Find(getDestFieldName(sourceType, sourcePath))
			if !ok {
				continue
			}

			srcFld := sourceStorage.MustFind(sourcePath)
			if err := setFieldRecursive(srcFld, destFld, source, dest); err != nil {
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

func setFieldRecursive(sourceFld, destFld fmap.Field, source, dest any) error {
	if r, ok := getRouteIfExists(sourceFld, destFld); ok {
		return r(sourceFld.Get(source), destFld.GetPtr(dest))
	}

	if sourceFld.GetType().Kind() != reflect.Struct {
		sourceVal := sourceFld.Get(source)
		if sourceVal != nil {
			destFld.Set(dest, sourceVal)
		}
		return nil
	}

	sourceStructField := sourceFld.GetPtr(source)
	sourceStorage, _ := fmap.GetFrom(sourceStructField)

	destStructField := destFld.GetPtr(dest)
	destStorage, _ := fmap.GetFrom(destStructField)

	for _, sPath := range sourceStorage.GetAllPaths() {
		sField := sourceStorage.MustFind(sPath)
		dPath := getDestFieldName(sField.GetType(), sPath)
		dField, ok := destStorage.Find(dPath)
		if !ok {
			continue
		}
		err := setFieldRecursive(sField, dField, sourceStructField, destStructField)
		if err != nil {
			return err
		}
	}
	return nil
}

func getRouteIfExists(sourceFld, destFld fmap.Field) (func(source interface{}, dest interface{}) error, bool) {
	destType := destFld.GetType()
	sourceType := sourceFld.GetType()
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
