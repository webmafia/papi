package valid

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/webbmaffian/papi/internal"
)

func Example() {
	type foo struct {
		A uint16    `min:"123" flags:"required"`
		B *[]uint64 `enum:"1,2,3"`
		C struct {
			Foo   string `enum:"foo"`
			Bar   string `enum:"bar"`
			Baz   [4]int `enum:"7,8,9" flags:"required"`
			Bazzo []int  `enum:"7,8,9" flags:"required"`
		}
	}

	valid, err := createStructValidator(internal.ReflectType[foo]())

	if err != nil {
		panic(err)
	}

	f := foo{
		A: 122,
		B: &[]uint64{1, 2, 3, 4},
	}
	f.C.Foo = "baz"
	f.C.Bar = "baz"
	// f.C.Baz = [4]int{1}

	var errs FieldErrors

	valid(unsafe.Pointer(&f), &errs)

	// fmt.Printf("%#v\n", errs)
	fmt.Println(errs)

	// Output: TODO
}

func Benchmark_validation(b *testing.B) {
	type foo struct {
		MyValue int `min:"123" flags:"required"`
	}

	valid, err := createStructValidator(internal.ReflectType[foo]())

	if err != nil {
		return
	}

	f := foo{
		MyValue: 122,
	}

	var errs FieldErrors

	b.ResetTimer()

	for range b.N {
		valid(unsafe.Pointer(&f), &errs)
		errs = errs[:0]
	}
}
