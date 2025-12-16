package internal

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/modern-go/reflect2"
)

type foo struct {
	bar int
	_   [1]byte
}

func BenchmarkNewSetter(b *testing.B) {
	var f foo

	ptr := unsafe.Pointer(&f)
	typ := reflect.TypeOf(f)

	for b.Loop() {
		_ = NewSetter(typ, ptr)
	}
}

func BenchmarkSetterSet(b *testing.B) {
	var f foo
	set := NewSetter(reflect.TypeOf(f), unsafe.Pointer(&f))

	for b.Loop() {
		if err := set.Set(&foo{bar: 123}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetterSetNil(b *testing.B) {
	var f foo
	set := NewSetter(reflect.TypeOf(f), unsafe.Pointer(&f))

	for b.Loop() {
		if err := set.Set(nil); err != nil {
			b.Fatal(err)
		}
	}
}

func ExampleSetter() {
	var f foo

	ptr := unsafe.Pointer(&f)
	typ := reflect.TypeOf(f)

	set := Setter{
		typ: reflect2.Type2(typ),
		ptr: ptr,
	}

	fmt.Println(f)
	set.Set(&foo{bar: 123})
	fmt.Println(f)
	set.Set(nil)
	fmt.Println(f)

	// Output: TODO
}
