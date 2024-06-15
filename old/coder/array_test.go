package coder

import (
	"fmt"
	"reflect"
)

func ExampleArray_scan() {
	var a Array

	if err := a.ScanTags(reflect.StructTag(`min:"10" max:"12" unique:"true" readOnly:"true"`)); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", a)

	// Output: schema.Array{Min:10, Max:12, Items:schema.Schema(nil), UniqueItems:true, Flags:schema.Flags{Nullable:false, ReadOnly:true, WriteOnly:false}}
}
