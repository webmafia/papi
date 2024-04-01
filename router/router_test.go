package router

import (
	"fmt"
	"testing"
	"unsafe"
)

func Example() {
	r := New()
	v := 123

	r.Add("/foo/bar/{baz}", unsafe.Pointer(&v))
	// r.Add("/foo/bar/mjau", unsafe.Pointer(&v))
	var params Params

	res := r.Lookup("/foo/bar/456", &params)

	if res == nil {
		fmt.Println("not found")
	} else {
		val := *(*int)(res)
		fmt.Println(val)
		fmt.Println(params.Get("baz"))
	}

	// Output: Mjaa
}

func BenchmarkAdd(b *testing.B) {
	r := New()
	v := 123

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Add("/foo/bar/{baz}", unsafe.Pointer(&v))
		r.Clear()
	}
}

func BenchmarkLookup(b *testing.B) {
	r := New()
	v := 123

	r.Add("/foo/bar/{baz}", unsafe.Pointer(&v))
	r.Add("/foo/bar/{baz}/mjau", unsafe.Pointer(&v))
	// r.Add("/foo/bar/mjau", unsafe.Pointer(&v))
	var params Params

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = r.Lookup("/foo/bar/123", &params)
		params.reset()
	}
}
