package internal

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type myType struct{}

func (myType) String() string {
	return ""
}

func ExampleImplementsInterface() {
	fmt.Println(ImplementsInterface[myType, fmt.Stringer]())

	// Output:
	//
	// true
}

func BenchmarkImplementsInterface(b *testing.B) {
	b.Run("true", func(b *testing.B) {
		for range b.N {
			_ = ImplementsInterface[myType, fmt.Stringer]()
		}
	})

	b.Run("false", func(b *testing.B) {
		for range b.N {
			_ = ImplementsInterface[int64, fmt.Stringer]()
		}
	})
}

func Example_ptrToInterface() {
	v := &bytes.Buffer{}
	typ := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	inter := reflect.NewAt(typ, unsafe.Pointer(v)).Interface()

	fmt.Printf("%#v\n", inter)

	// Output: TODO
}

func Benchmark_ptrToInterface(b *testing.B) {
	v := &bytes.Buffer{}
	ptr := unsafe.Pointer(v)

	b.ResetTimer()

	for range b.N {
		var str fmt.Stringer
		ptrToInterface(&str, ptr)
		_ = str
	}
}
