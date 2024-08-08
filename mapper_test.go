package gomapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestingStructSource struct {
	Name string
}

type TestingStructDest struct {
	Name string
}

var converterFunc = func(source TestingStructSource, dest *TestingStructDest) error {
	dest.Name = source.Name
	return nil
}

func TestAddRoute(t *testing.T) {
	err := AddRoute[TestingStructSource, TestingStructDest](converterFunc)
	assert.NoError(t, err)
}

func TestMapTo(t *testing.T) {
	_ = AddRoute[TestingStructSource, TestingStructDest](converterFunc)
	t.Run("Source is a pointer to struct", func(t *testing.T) {
		source := &TestingStructSource{Name: "Test1"}
		dest, err := MapTo[TestingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, source.Name, dest.Name)
	})
	t.Run("Source is a struct", func(t *testing.T) {
		source := TestingStructSource{Name: "Test1"}
		dest, err := MapTo[TestingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, source.Name, dest.Name)
	})
	t.Run("Source is a slice", func(t *testing.T) {
		source := []TestingStructSource{
			{
				Name: "ArrayTest1",
			},
			{
				Name: "ArrayTest2",
			},
			{
				Name: "ArrayTest3",
			},
		}
		dest, err := MapTo[[]TestingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, len(source), len(dest))
		for i := range source {
			assert.Equal(t, source[i].Name, dest[i].Name)
		}
	})
	t.Run("Source is an empty slice", func(t *testing.T) {
		source := make([]TestingStructSource, 0)
		dest, err := MapTo[[]TestingStructDest](source)
		assert.NoError(t, err)
		assert.NotNil(t, source)
		assert.NotNil(t, dest)
		assert.Equal(t, len(source), len(dest))
	})
	t.Run("Source is a slice with pointer elements", func(t *testing.T) {
		source := []*TestingStructSource{
			{
				Name: "ArrayTest1",
			},
			{
				Name: "ArrayTest2",
			},
			{
				Name: "ArrayTest3",
			},
		}
		dest, err := MapTo[[]TestingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, len(source), len(dest))
		for i := range source {
			assert.Equal(t, source[i].Name, dest[i].Name)
		}
	})
	t.Run("Source is an empty slice with pointer elements", func(t *testing.T) {
		source := make([]*TestingStructSource, 0)
		dest, err := MapTo[[]TestingStructDest](source)
		assert.NoError(t, err)
		assert.NotNil(t, source)
		assert.NotNil(t, dest)
		assert.Equal(t, len(source), len(dest))
	})
	t.Run("Source is a slice with pointer elements, dest is a slice with pointer elements", func(t *testing.T) {
		source := []*TestingStructSource{
			{
				Name: "ArrayTest1",
			},
			{
				Name: "ArrayTest2",
			},
			{
				Name: "ArrayTest3",
			},
		}
		dest, err := MapTo[[]*TestingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, len(source), len(dest))
		for i := range source {
			assert.Equal(t, source[i].Name, dest[i].Name)
		}
	})
	t.Run("Source is an empty slice with pointer elements, dest is a slice with pointer elements", func(t *testing.T) {
		source := make([]*TestingStructSource, 0)
		dest, err := MapTo[[]*TestingStructDest](source)
		assert.NoError(t, err)
		assert.NotNil(t, source)
		assert.NotNil(t, dest)
		assert.Equal(t, len(source), len(dest))
	})
	t.Run("Source is a slice pointer", func(t *testing.T) {
		source := &[]TestingStructSource{
			{
				Name: "ArrayTest1",
			},
			{
				Name: "ArrayTest2",
			},
			{
				Name: "ArrayTest3",
			},
		}
		dest, err := MapTo[[]TestingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, len(*source), len(dest))
		for i := range *source {
			assert.Equal(t, (*source)[i].Name, dest[i].Name)
		}
	})
	t.Run("Source is an empty slice pointer", func(t *testing.T) {
		arr := make([]TestingStructSource, 0)
		source := &arr
		dest, err := MapTo[[]TestingStructDest](source)
		assert.NoError(t, err)
		assert.NotNil(t, source)
		assert.NotNil(t, dest)
		assert.Equal(t, len(*source), len(dest))
	})
	t.Run("Dest is a slice with pointer elements", func(t *testing.T) {
		source := []TestingStructSource{
			{
				Name: "ArrayTest1",
			},
			{
				Name: "ArrayTest2",
			},
			{
				Name: "ArrayTest3",
			},
		}
		dest, err := MapTo[[]*TestingStructDest](source)
		assert.NoError(t, err)
		assert.Equal(t, len(source), len(dest))
		for i := range source {
			assert.Equal(t, (source)[i].Name, dest[i].Name)
		}
	})
	t.Run("Dest is an empty slice with pointer elements", func(t *testing.T) {
		source := make([]TestingStructSource, 0)
		dest, err := MapTo[[]*TestingStructDest](source)
		assert.NoError(t, err)
		assert.NotNil(t, source)
		assert.NotNil(t, dest)
		assert.Equal(t, len(source), len(dest))
	})
	t.Run("Source is a pointer to pointer", func(t *testing.T) {
		source := &TestingStructSource{Name: "Test1"}
		_, err := MapTo[TestingStructDest](&source)
		assert.Error(t, err)
	})
	t.Run("Dest is a pointer to pointer", func(t *testing.T) {
		source := &TestingStructSource{Name: "Test1"}
		_, err := MapTo[*TestingStructDest](source)
		assert.Error(t, err)
	})
}
func TestMap(t *testing.T) {
	_ = AddRoute[TestingStructSource, TestingStructDest](converterFunc)
	t.Run("Source is a pointer to struct", func(t *testing.T) {
		source := &TestingStructSource{Name: "Test1"}
		dest := &TestingStructDest{}
		err := Map(source, dest)
		assert.NoError(t, err)
		assert.Equal(t, source.Name, dest.Name)
	})
	t.Run("Source is a struct", func(t *testing.T) {
		source := TestingStructSource{Name: "Test1"}
		dest := &TestingStructDest{}
		err := Map(source, dest)
		assert.NoError(t, err)
		assert.Equal(t, source.Name, dest.Name)
	})
	t.Run("Source is a slice", func(t *testing.T) {
		source := []TestingStructSource{
			{
				Name: "ArrayTest1",
			},
			{
				Name: "ArrayTest2",
			},
			{
				Name: "ArrayTest3",
			},
		}
		var dest []TestingStructDest
		err := Map(source, &dest)
		assert.NoError(t, err)
		assert.Equal(t, len(source), len(dest))
		for i := range source {
			assert.Equal(t, source[i].Name, dest[i].Name)
		}
	})
	t.Run("Source is a pointer to slice", func(t *testing.T) {
		source := &[]TestingStructSource{
			{
				Name: "ArrayTest1",
			},
			{
				Name: "ArrayTest2",
			},
			{
				Name: "ArrayTest3",
			},
		}
		var dest []TestingStructDest
		err := Map(source, &dest)
		assert.NoError(t, err)
		assert.Equal(t, len(*source), len(dest))
		for i := range *source {
			assert.Equal(t, (*source)[i].Name, dest[i].Name)
		}
	})
	t.Run("Dest is a not a pointer", func(t *testing.T) {
		source := TestingStructSource{Name: "Test1"}
		dest := TestingStructDest{}
		err := Map(source, dest)
		assert.Error(t, err)
	})
	t.Run("Dest is nil", func(t *testing.T) {
		source := TestingStructSource{Name: "Test1"}
		err := Map(source, nil)
		assert.Error(t, err)
	})
	t.Run("Source is nil", func(t *testing.T) {
		dest := &TestingStructDest{}
		err := Map(nil, dest)
		assert.Error(t, err)
	})
}
