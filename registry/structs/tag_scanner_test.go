package structs

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/internal"
	"github.com/webbmaffian/papi/registry"
)

func ExampleCreateTagScanner() {
	type Foo struct {
		A string  `tag:"a"`
		B int     `tag:"b"`
		C float64 `tag:"c"`
		D bool    `tag:"d"`
	}

	scan, err := CreateTagScanner(registry.NewRegistry(), internal.ReflectType[Foo]())

	if err != nil {
		panic(err)
	}

	tags := reflect.StructTag(`a:"foobar" b:"123" c:"456.789" d:"true" e:"nothing"`)

	var foo Foo

	if err = scan(unsafe.Pointer(&foo), string(tags)); err != nil {
		return
	}

	fmt.Printf("%#v\n", foo)

	// Output: structs.Foo{A:"foobar", B:123, C:456.789, D:true}
}
