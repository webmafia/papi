package structs

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/internal"
	"github.com/webbmaffian/papi/registry/scanner"
	"github.com/webbmaffian/papi/registry/types"
)

func ExampleCreateTagScanner() {
	type Foo struct {
		A string  `tag:"a"`
		B int     `tag:"b"`
		C float64 `tag:"c"`
		D bool    `tag:"d"`
		E bool    `tag:"flags:e"`
		F bool    `tag:"flags:f"`
		G bool    `tag:"flags:g"`
	}

	scan, err := CreateTagScanner(internal.ReflectType[Foo](), func(typ reflect.Type, _ reflect.StructTag) (scan types.ParamDecoder, err error) {
		return scanner.CreateScanner(typ)
	})

	if err != nil {
		panic(err)
	}

	tags := reflect.StructTag(`a:"foobar" b:"123" c:"456.789" d:"true" h:"nothing" flags:"f,g,x"`)

	var foo Foo

	if err = scan(unsafe.Pointer(&foo), string(tags)); err != nil {
		return
	}

	fmt.Printf("%#v\n", foo)

	// Output: structs.Foo{A:"foobar", B:123, C:456.789, D:true, E:false, F:true, G:true}
}
