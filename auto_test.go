package gomapper

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type AutoMappingStructSource struct {
	Name         string
	Time         time.Time
	UUID         uuid.UUID
	PtrTime      *time.Time
	PtrUUID      *uuid.UUID
	NestedStruct NestedStructSource
}

type AutoMappingStructDest struct {
	Name         string
	SecondName   string
	Time         time.Time
	UUID         uuid.UUID
	PtrTime      *time.Time
	PtrUUID      *uuid.UUID
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
	_ = AutoRoute[TestingStructSource, TestingStructDest]()
	t.Run("Auto route without options", func(t *testing.T) {
		ptrTime := time.Now()
		ptrUuid := uuid.New()
		source := &AutoMappingStructSource{Name: "Test1", Time: time.Now(), UUID: uuid.New(), PtrTime: &ptrTime, PtrUUID: &ptrUuid}
		dest, err := MapTo[AutoMappingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, source.Name, dest.Name)
		assert.Equal(t, source.Time, dest.Time)
		assert.Equal(t, source.UUID, dest.UUID)
		assert.Equal(t, source.PtrUUID, dest.PtrUUID)
		assert.Equal(t, source.PtrTime, dest.PtrTime)
	})
	timeNow := time.Now()
	_ = AutoRoute[AutoMappingStructSource, AutoMappingStructDest](WithFunc(func(source AutoMappingStructSource, dest *AutoMappingStructDest) {
		if source.Name == "Test1" {
			dest.SecondName = "Test2"
		}
		if source.PtrTime == nil {
			dest.PtrTime = &timeNow
		}
		dest.Time = timeNow
	}))
	t.Run("Auto route with options", func(t *testing.T) {
		source := &AutoMappingStructSource{Name: "Test1"}
		dest, err := MapTo[AutoMappingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, source.Name, dest.Name)
		assert.Equal(t, "Test2", dest.SecondName)
		assert.Equal(t, timeNow, *dest.PtrTime)
		assert.Equal(t, timeNow, dest.Time)
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
		assert.Equal(t, source.NestedStruct.FirstNestedName, dest.NestedStruct.FirstNestedName)
	})
}
