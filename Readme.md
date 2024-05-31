[![codecov](https://codecov.io/github/Insei/gomapper/graph/badge.svg?token=85LGN4NOFA)](https://codecov.io/github/Insei/gomapper)
[![build](https://github.com/insei/gomapper/actions/workflows/go.yml/badge.svg)](https://github.com/Insei/gomapper/actions/workflows/go.yml)
[![Goreport](https://goreportcard.com/badge/github.com/insei/gomapper)](https://goreportcard.com/report/github.com/insei/gomapper)
[![GoDoc](https://godoc.org/github.com/insei/gomapper?status.svg)](https://godoc.org/github.com/insei/gomapper)
# GoMapper 
GoMapper is a library for struct to struct mapping.
There are two use cases: Manual and Auto.<br>
* `Manual` mode allows you to specify a function to convert one structure to another.<br>
* `Auto` mode uses matching field names for automatic conversion; 
it is important that not only the field names match, but also their types. 
This mode also supports structures in structure fields and automatically works by matching field names. 
It's based on [fmap](https://github.com/insei/fmap) switch case and reflect based library.

Both modes are route based. In which the reflect.Type of the source structure and the type of the destination structure 
are specified. If such a route was not found, gomapper will return an error.

Also `gomapper` support slices, you don't need to specify types of slices for mapping.

## Installation

```bash
go get github.com/insei/gomapper@latest
```

## Examples
You can found a lot of examples in tests.<br>
Manual route.

```go
package main

import (
	"fmt"

	"github.com/insei/gomapper"
)

type Source struct {
	Name string
	Age  uint8
}

type Dest struct {
	NameCustom string
	Age        uint8
}

func main() {
	err := gomapper.AddRoute[Source, Dest](func(source Source, dest *Dest) error {
		dest.NameCustom = source.Name
		dest.Age = source.Age
		return nil
	})
	if err != nil {
		panic(err)
	}
	s := Source{
		Name: "DefaultName",
		Age:  16,
	}
	dest, err := gomapper.MapTo[Dest](s) 
	// or gomapper.MapTo[Dest](&s)
	// or d := Dest{}
	// gomapper.Map(&s, &d)
	if err != nil {
		panic(err)
	}
	fmt.Print(dest)
}
```
Auto Route.
```go
package main

import (
	"fmt"

	"github.com/insei/gomapper"
)

type Source struct {
	Name string
	Age  uint8
}

type Dest struct {
	NameCustom string
	Age        uint8
}

func main() {
	err := gomapper.AutoRoute[Source, Dest]()
	if err != nil {
		panic(err)
	}
	s := Source{
		Name: "DefaultName",
		Age:  16,
	}
	dest, err := gomapper.MapTo[Dest](s) 
	// or gomapper.MapTo[Dest](&s)
	// or dest := Dest{}
	// gomapper.Map(&s, &dest)
	if err != nil {
		panic(err)
	}
	fmt.Print(dest)
}
```
Map structs into slices.
```go
package main

import (
	"fmt"

	"github.com/insei/gomapper"
)

type Source struct {
	Name string
	Age  uint8
}

type Dest struct {
	NameCustom string
	Age        uint8
}

func main() {
	err := gomapper.AutoRoute[Source, Dest]() // or manual
	if err != nil {
		panic(err)
	}
	s := Source{
		Name: "DefaultName",
		Age:  16,
	}
	sSlice := []Source{ s }
	sDest, err := gomapper.MapTo[[]Dest](sSlice) 
	// or sDest := []Dest{}
	// sDest, err := gomapper.MapTo(sSlice, &sDest) 
	if err != nil {
		panic(err)
	}
	fmt.Print(sDest)
}
```
