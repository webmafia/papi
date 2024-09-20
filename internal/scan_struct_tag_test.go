package internal

import (
	"fmt"
	"reflect"
	"unsafe"
)

func ExampleCreateStructTagScanner() {
	type Foo struct {
		A string  `tag:"a"`
		B int     `tag:"b"`
		C float64 `tag:"c"`
		D bool    `tag:"d"`
	}

	scan, err := CreateStructTagScanner(ReflectType[Foo]())

	if err != nil {
		panic(err)
	}

	tags := reflect.StructTag(`a:"foobar" b:"123" c:"456.789" d:"true" e:"nothing"`)

	var foo Foo

	if err = scan(unsafe.Pointer(&foo), tags); err != nil {
		return
	}

	fmt.Printf("%#v\n", foo)

	// Output: internal.Foo{A:"foobar", B:123, C:456.789, D:true}
}
