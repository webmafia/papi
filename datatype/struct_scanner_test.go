package datatype

import (
	"fmt"
	"testing"
)

func ExampleCreateStructScanner() {
	type Foobar struct {
		Foo string `tag:"foo"`
		Bar int    `tag:"bar"`
		Baz bool   `tag:"baz"`
	}

	d := NewDataTypes()
	scan, err := CreateStructScanner[Foobar](d, "tag")

	if err != nil {
		panic(err)
	}

	var f Foobar
	fmt.Printf("%#v\n", f)

	err = scan(&f, func(tag, val string) string {
		fmt.Printf("scanning value '%d' to field '%s=%s'\n", 1, tag, val)
		return "1"
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", f)

	// Output:
	//
	// datatype.Foobar{Foo:"", Bar:0, Baz:false}
	// scanning value '1' to field 'tag=foo'
	// scanning value '1' to field 'tag=bar'
	// scanning value '1' to field 'tag=baz'
	// datatype.Foobar{Foo:"1", Bar:1, Baz:true}
}

func BenchmarkCreateStructScanner(b *testing.B) {
	type Foobar struct {
		Foo string `tag:"foo"`
		Bar int    `tag:"bar"`
		Baz bool   `tag:"baz"`
	}

	d := NewDataTypes()
	scan, err := CreateStructScanner[Foobar](d, "tag")

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var f Foobar

		err := scan(&f, func(tag, val string) string {
			return "1"
		})

		if err != nil {
			b.Fatal(err)
		}
	}
}
