package gomapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestingStruct struct {
	Name string
}

type TestingStruct2 struct {
	Name string
}

func Test(t *testing.T) {
	err := AddRoute[TestingStruct, *TestingStruct2](func(source TestingStruct, dest *TestingStruct2) error {
		dest.Name = source.Name
		return nil
	})
	assert.NoError(t, err)

	source1 := &TestingStruct{Name: "Test1"}
	dest1, err := MapTo[TestingStruct2](source1)
	assert.NoError(t, err)
	assert.Equal(t, source1.Name, dest1.Name)

	source2 := TestingStruct{Name: "Test2"}
	dest2, err := MapTo[TestingStruct2](source2)
	assert.NoError(t, err)
	assert.Equal(t, source2.Name, dest2.Name)

	source3 := &TestingStruct{Name: "Test3"}
	dest3 := &TestingStruct2{}
	err = Map(source3, dest3)
	assert.NoError(t, err)
	assert.Equal(t, source3.Name, dest3.Name)

	source4 := TestingStruct{Name: "Test4"}
	dest4 := &TestingStruct2{}
	err = Map(source4, dest4)
	assert.NoError(t, err)
	assert.Equal(t, source4.Name, dest4.Name)
}
