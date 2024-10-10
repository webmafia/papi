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
		B [4]uint64 `enum:"1,2,3"`
	}

	valid, err := createStructValidator(internal.ReflectType[foo]())

	if err != nil {
		panic(err)
	}

	f := foo{
		A: 122,
		B: [4]uint64{1, 2, 2},
	}

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
