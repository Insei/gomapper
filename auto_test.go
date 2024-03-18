package gomapper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type AutoMappingStructSource struct {
	Name         string
	NestedStruct NestedStructSource
}

type AutoMappingStructDest struct {
	Name         string
	SecondName   string
	NestedStruct NestedStructDest
}

type NestedStructSource struct {
	FirstNestedName  string
	DeepNestedStruct DeepNestedStructSource
}

type NestedStructDest struct {
	FirstNestedName       string
	FirstNestedSecondName string
	DeepNestedStruct      DeepNestedStructDest
}

type DeepNestedStructSource struct {
	SecondNestedName string
}

type DeepNestedStructDest struct {
	SecondNestedName string
}

func TestAutoRoute(t *testing.T) {
	_ = AutoRoute[AutoMappingStructSource, AutoMappingStructDest]()
	t.Run("Auto route without options", func(t *testing.T) {
		source := &AutoMappingStructSource{Name: "Test1"}
		dest, err := MapTo[AutoMappingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, source.Name, dest.Name)
	})
	_ = AutoRoute[AutoMappingStructSource, AutoMappingStructDest](WithFieldRoute("Name", "SecondName"))
	t.Run("Auto route with options", func(t *testing.T) {
		source := &AutoMappingStructSource{Name: "Test1"}
		dest, err := MapTo[AutoMappingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, "", dest.Name)
		assert.Equal(t, source.Name, dest.SecondName)
	})
	t.Run("Auto mapping struct fields", func(t *testing.T) {
		source := &AutoMappingStructSource{
			NestedStruct: NestedStructSource{
				FirstNestedName: "Test1",
				DeepNestedStruct: DeepNestedStructSource{
					SecondNestedName: "Test2",
				},
			},
		}
		dest, err := MapTo[AutoMappingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, source.NestedStruct.FirstNestedName, dest.NestedStruct.FirstNestedName)
		assert.Equal(t, source.NestedStruct.DeepNestedStruct.SecondNestedName, dest.NestedStruct.DeepNestedStruct.SecondNestedName)
	})

	_ = AddRoute[NestedStructSource, NestedStructDest](func(source NestedStructSource, dest *NestedStructDest) error {
		dest.FirstNestedSecondName = source.FirstNestedName
		return nil
	})
	t.Run("Auto mapping using existing route on nested struct", func(t *testing.T) {
		source := &AutoMappingStructSource{
			NestedStruct: NestedStructSource{
				FirstNestedName: "Test1",
			},
		}
		dest, err := MapTo[AutoMappingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, "", dest.NestedStruct.FirstNestedName)
		assert.Equal(t, source.NestedStruct.FirstNestedName, dest.NestedStruct.FirstNestedSecondName)
	})
}
