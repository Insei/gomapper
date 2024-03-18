package fields

import (
	"reflect"
	"unsafe"
)

type Field reflect.StructField

var fields = map[reflect.Type]map[string]Field{}

// Get returns the value of the fields in the provided object.
// It takes a parameter `obj` of type `interface{}`, representing the object.
// It returns the value of the fields as an `interface{}`.
// If the fields type is `string`, it returns the value as a `string`.
// If the fields type is `int`, it returns the value as an `int`.
// If the fields type is `bool`, it returns the value as a `bool`.
// If the fields type is not handled, it panics with an error message.
func (f Field) Get(obj interface{}) interface{} {
	ptrToField := f.getPtr(obj)
	kind := f.Type.Kind()
	isPtr := false
	if kind == reflect.Ptr {
		isPtr = true
		kind = f.Type.Elem().Kind()
	}
	if isPtr {
		switch kind {
		case reflect.String:
			return getPtrValue[*string](ptrToField)
		case reflect.Int:
			return getPtrValue[*int](ptrToField)
		case reflect.Int8:
			return getPtrValue[*int8](ptrToField)
		case reflect.Int16:
			return getPtrValue[*int16](ptrToField)
		case reflect.Int32:
			return getPtrValue[*int32](ptrToField)
		case reflect.Int64:
			return getPtrValue[*int64](ptrToField)
		case reflect.Float32:
			return getPtrValue[*float32](ptrToField)
		case reflect.Float64:
			return getPtrValue[*float64](ptrToField)
		case reflect.Bool:
			return getPtrValue[*bool](ptrToField)
		case reflect.Struct:
			return reflect.NewAt(f.Type, ptrToField).Interface()
		default:
			panic("unhandled default case")
		}
	} else {
		switch kind {
		case reflect.String:
			return getPtrValue[string](ptrToField)
		case reflect.Int:
			return getPtrValue[int](ptrToField)
		case reflect.Int8:
			return getPtrValue[int8](ptrToField)
		case reflect.Int16:
			return getPtrValue[int16](ptrToField)
		case reflect.Int32:
			return getPtrValue[int32](ptrToField)
		case reflect.Int64:
			return getPtrValue[int64](ptrToField)
		case reflect.Float32:
			return getPtrValue[float32](ptrToField)
		case reflect.Float64:
			return getPtrValue[float64](ptrToField)
		case reflect.Bool:
			return getPtrValue[bool](ptrToField)
		case reflect.Struct:
			return reflect.NewAt(f.Type, ptrToField).Interface()
		default:
			panic("unhandled default case")
		}
	}
}

// getPtr returns a pointer to the field's value in the provided configuration object.
// It takes a parameter `conf` of type `any`, representing the configuration object.
// It returns an `unsafe.Pointer` to the `field's` value in the configuration object.
func (f Field) getPtr(conf interface{}) unsafe.Pointer {
	confPointer := ((*[2]unsafe.Pointer)(unsafe.Pointer(&conf)))[1]
	ptToField := unsafe.Add(confPointer, f.Offset)
	return ptToField
}

func setPtrValue[T any](ptr unsafe.Pointer, val any) {
	valSet := (*T)(ptr)
	*valSet = val.(T)
}

func getPtrValue[T any](ptr unsafe.Pointer) T {
	return *(*T)(ptr)
}

// Set updates the value of the fields in the provided object with the provided value.
// It takes two parameters:
//   - obj: interface{}, representing the object containing the fields.
//   - val: interface{}, representing the new value for the fields.
//
// The Set method uses the getPtr method to get a pointer to the fields in the object.
// It then performs a type switch on the kind of the fields to determine its type, and sets the value accordingly.
// The supported fields types are string, int, and bool.
// If the fields type is not one of the supported types, it panics with the message "unhandled default case".
func (f Field) Set(obj interface{}, val interface{}) {
	ptrToField := f.getPtr(obj)
	kind := f.Type.Kind()
	isPtr := false
	if kind == reflect.Ptr {
		isPtr = true
		kind = f.Type.Elem().Kind()
	}
	if isPtr {
		switch kind {
		case reflect.String:
			setPtrValue[*string](ptrToField, val)
		case reflect.Int:
			setPtrValue[*int](ptrToField, val)
		case reflect.Int8:
			setPtrValue[*int8](ptrToField, val)
		case reflect.Int16:
			setPtrValue[*int16](ptrToField, val)
		case reflect.Int32:
			setPtrValue[*int32](ptrToField, val)
		case reflect.Int64:
			setPtrValue[*int64](ptrToField, val)
		case reflect.Float32:
			setPtrValue[*float32](ptrToField, val)
		case reflect.Float64:
			setPtrValue[*float64](ptrToField, val)
		case reflect.Bool:
			setPtrValue[*bool](ptrToField, val)
		default:
			panic("unhandled default case")
		}
	} else {
		switch kind {
		case reflect.String:
			setPtrValue[string](ptrToField, val)
		case reflect.Int:
			setPtrValue[int](ptrToField, val)
		case reflect.Int8:
			setPtrValue[int8](ptrToField, val)
		case reflect.Int16:
			setPtrValue[int16](ptrToField, val)
		case reflect.Int32:
			setPtrValue[int32](ptrToField, val)
		case reflect.Int64:
			setPtrValue[int64](ptrToField, val)
		case reflect.Float32:
			setPtrValue[float32](ptrToField, val)
		case reflect.Float64:
			setPtrValue[float64](ptrToField, val)
		case reflect.Bool:
			setPtrValue[bool](ptrToField, val)
		default:
			panic("unhandled default case")
		}
	}
}

// Get returns a map of FieldGetSet objects representing the fields of the provided object.
// It takes a parameter `conf` of type `any`, representing the object.
// It returns a map with string keys (fields path) and FieldGetSet values.
// FieldGetSet is an interface that defines Get and Set methods for fields.
func Get[T interface{}]() map[string]Field {
	obj := new(T)
	return GetFrom(obj)
}

func GetFrom(obj interface{}) map[string]Field {
	typeOf := reflect.TypeOf(obj)
	if tFields, ok := fields[typeOf]; ok {
		return tFields
	}
	tFields := map[string]Field{}
	getFieldsMapRecursive(obj, "", &tFields)
	fields[typeOf] = tFields
	return tFields
}

func getFieldsMapRecursive(conf any, path string, f *map[string]Field) {
	typeOf := reflect.TypeOf(conf)
	valueOf := reflect.ValueOf(conf)
	if reflect.ValueOf(conf).Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	if path != "" {
		path += "."
	}
	for i := 0; i < typeOf.NumField(); i++ {
		fieldTypeOf := typeOf.Field(i)
		fieldValueOf := valueOf.Field(i)
		switch fieldTypeOf.Type.Kind() {
		case reflect.Slice:
			break
		case reflect.Struct:
			(*f)[path+fieldTypeOf.Name] = Field(fieldTypeOf)
			getFieldsMapRecursive(fieldValueOf.Addr().Interface(), path+fieldTypeOf.Name, f)
		default:
			(*f)[path+fieldTypeOf.Name] = Field(fieldTypeOf)
		}
	}
}
