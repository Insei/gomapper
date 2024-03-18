package fields

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type NestedStruct struct {
	String string
}

type TestStruct struct {
	String       string
	Int          int
	Int8         int8
	Int16        int16
	Int32        int32
	Int64        int64
	Float32      float32
	Float64      float64
	Bool         bool
	PtrString    *string
	PtrInt       *int
	PtrInt8      *int8
	PtrInt16     *int16
	PtrInt32     *int32
	PtrInt64     *int64
	PtrFloat32   *float32
	PtrFloat64   *float64
	PtrBool      *bool
	NestedStruct NestedStruct
	Slice        []string
}

func TestGetFrom(t *testing.T) {
	t.Run("Get fields map from struct", func(t *testing.T) {
		fieldMap := GetFrom(&TestStruct{})
		assert.Contains(t, fieldMap, "String")
		assert.Contains(t, fieldMap, "Int")
		assert.Contains(t, fieldMap, "Int8")
		assert.Contains(t, fieldMap, "Int16")
		assert.Contains(t, fieldMap, "Int32")
		assert.Contains(t, fieldMap, "Int64")
		assert.Contains(t, fieldMap, "Float32")
		assert.Contains(t, fieldMap, "Float64")
		assert.Contains(t, fieldMap, "Bool")
		assert.Contains(t, fieldMap, "PtrString")
		assert.Contains(t, fieldMap, "PtrInt")
		assert.Contains(t, fieldMap, "PtrInt8")
		assert.Contains(t, fieldMap, "PtrInt16")
		assert.Contains(t, fieldMap, "PtrInt32")
		assert.Contains(t, fieldMap, "PtrInt64")
		assert.Contains(t, fieldMap, "PtrFloat32")
		assert.Contains(t, fieldMap, "PtrFloat64")
		assert.Contains(t, fieldMap, "PtrBool")
		assert.Contains(t, fieldMap, "NestedStruct")
		assert.Contains(t, fieldMap, "NestedStruct.String")
		assert.NotContains(t, fieldMap, "Slice")
	})
}

func TestGet(t *testing.T) {
	t.Run("Get fields from struct", func(t *testing.T) {
		strVal := "Test2"
		intVal := 6
		int8Val := int8(7)
		int16Val := int16(8)
		int32Val := int32(9)
		int64Val := int64(10)
		f32Val := float32(11.11)
		f64Val := 12.12
		boolVal := true
		source := &TestStruct{
			String:       "Test1",
			Int:          1,
			Int8:         2,
			Int16:        3,
			Int32:        4,
			Int64:        5,
			Float32:      6.6,
			Float64:      7.7,
			Bool:         true,
			PtrString:    &strVal,
			PtrInt:       &intVal,
			PtrInt8:      &int8Val,
			PtrInt16:     &int16Val,
			PtrInt32:     &int32Val,
			PtrInt64:     &int64Val,
			PtrFloat32:   &f32Val,
			PtrFloat64:   &f64Val,
			PtrBool:      &boolVal,
			NestedStruct: NestedStruct{String: "Nested"},
		}
		mainStructMap := GetFrom(source)
		nestedStructMap := GetFrom(mainStructMap["NestedStruct"].Get(source))
		nestedField := nestedStructMap["String"].Get(mainStructMap["NestedStruct"].Get(source))
		assert.Equal(t, "Test1", mainStructMap["String"].Get(source))
		assert.Equal(t, 1, mainStructMap["Int"].Get(source))
		assert.Equal(t, int8(2), mainStructMap["Int8"].Get(source))
		assert.Equal(t, int16(3), mainStructMap["Int16"].Get(source))
		assert.Equal(t, int32(4), mainStructMap["Int32"].Get(source))
		assert.Equal(t, int64(5), mainStructMap["Int64"].Get(source))
		assert.Equal(t, float32(6.6), mainStructMap["Float32"].Get(source))
		assert.Equal(t, 7.7, mainStructMap["Float64"].Get(source))
		assert.Equal(t, true, mainStructMap["Bool"].Get(source))
		assert.Equal(t, &strVal, mainStructMap["PtrString"].Get(source))
		assert.Equal(t, &intVal, mainStructMap["PtrInt"].Get(source))
		assert.Equal(t, &int8Val, mainStructMap["PtrInt8"].Get(source))
		assert.Equal(t, &int16Val, mainStructMap["PtrInt16"].Get(source))
		assert.Equal(t, &int32Val, mainStructMap["PtrInt32"].Get(source))
		assert.Equal(t, &int64Val, mainStructMap["PtrInt64"].Get(source))
		assert.Equal(t, &f32Val, mainStructMap["PtrFloat32"].Get(source))
		assert.Equal(t, &f64Val, mainStructMap["PtrFloat64"].Get(source))
		assert.Equal(t, &boolVal, mainStructMap["PtrBool"].Get(source))
		assert.Equal(t, "Nested", nestedField)
	})
}

func TestSet(t *testing.T) {
	t.Run("Set fields in struct", func(t *testing.T) {
		dest := &TestStruct{}
		strVal := "NewTest2"
		intVal := 6
		int8Val := int8(7)
		int16Val := int16(8)
		int32Val := int32(9)
		int64Val := int64(10)
		f32Val := float32(11.11)
		f64Val := 12.12
		boolVal := true
		nested := "NewNested"

		fieldMap := GetFrom(dest)
		nestedMap := GetFrom(fieldMap["NestedStruct"].Get(dest))
		fieldMap["String"].Set(dest, "NewTest")
		fieldMap["Int"].Set(dest, 1)
		fieldMap["Int8"].Set(dest, int8(2))
		fieldMap["Int16"].Set(dest, int16(3))
		fieldMap["Int32"].Set(dest, int32(4))
		fieldMap["Int64"].Set(dest, int64(5))
		fieldMap["Float32"].Set(dest, float32(11.11))
		fieldMap["Float64"].Set(dest, 12.12)
		fieldMap["Bool"].Set(dest, true)
		fieldMap["PtrString"].Set(dest, &strVal)
		fieldMap["PtrInt"].Set(dest, &intVal)
		fieldMap["PtrInt8"].Set(dest, &int8Val)
		fieldMap["PtrInt16"].Set(dest, &int16Val)
		fieldMap["PtrInt32"].Set(dest, &int32Val)
		fieldMap["PtrInt64"].Set(dest, &int64Val)
		fieldMap["PtrFloat32"].Set(dest, &f32Val)
		fieldMap["PtrFloat64"].Set(dest, &f64Val)
		fieldMap["PtrBool"].Set(dest, &boolVal)
		nestedMap["String"].Set(fieldMap["NestedStruct"].Get(dest), nested)

		assert.Equal(t, "NewTest", fieldMap["String"].Get(dest))
		assert.Equal(t, 1, fieldMap["Int"].Get(dest))
		assert.Equal(t, int8(2), fieldMap["Int8"].Get(dest))
		assert.Equal(t, int16(3), fieldMap["Int16"].Get(dest))
		assert.Equal(t, int32(4), fieldMap["Int32"].Get(dest))
		assert.Equal(t, int64(5), fieldMap["Int64"].Get(dest))
		assert.Equal(t, f32Val, fieldMap["Float32"].Get(dest))
		assert.Equal(t, f64Val, fieldMap["Float64"].Get(dest))
		assert.Equal(t, boolVal, fieldMap["Bool"].Get(dest))
		assert.Equal(t, &strVal, fieldMap["PtrString"].Get(dest))
		assert.Equal(t, &intVal, fieldMap["PtrInt"].Get(dest))
		assert.Equal(t, &int8Val, fieldMap["PtrInt8"].Get(dest))
		assert.Equal(t, &int16Val, fieldMap["PtrInt16"].Get(dest))
		assert.Equal(t, &int32Val, fieldMap["PtrInt32"].Get(dest))
		assert.Equal(t, &int64Val, fieldMap["PtrInt64"].Get(dest))
		assert.Equal(t, &f32Val, fieldMap["PtrFloat32"].Get(dest))
		assert.Equal(t, &f64Val, fieldMap["PtrFloat64"].Get(dest))
		assert.Equal(t, &boolVal, fieldMap["PtrBool"].Get(dest))
		assert.Equal(t, "NewNested", fieldMap["String"].Get(fieldMap["NestedStruct"].Get(dest)))
	})
}
