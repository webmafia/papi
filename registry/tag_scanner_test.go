package registry

import (
	"fmt"
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/webbmaffian/papi/pool/json"
)

func Example_createTagScanner() {
	type Foo struct {
		A string  `tag:"a"`
		B int     `tag:"b"`
		C float64 `tag:"c"`
		D bool    `tag:"d"`
		E bool    `tag:"flags:e"`
		F bool    `tag:"flags:f"`
		G bool    `tag:"flags:g"`
	}

	reg, err := NewRegistry(json.NewPool(jsoniter.ConfigFastest))

	if err != nil {
		panic(err)
	}

	scan, err := reg.createTagScanner(reflect.TypeFor[Foo]())

	if err != nil {
		panic(err)
	}

	tags := reflect.StructTag(`a:"foobar" b:"123" c:"456.789" d:"true" h:"nothing" flags:"f,g,x"`)

	var foo Foo

	if err = scan(unsafe.Pointer(&foo), string(tags)); err != nil {
		return
	}

	fmt.Printf("%#v\n", foo)

	// Output: registry.Foo{A:"foobar", B:123, C:456.789, D:true, E:false, F:true, G:true}
}
