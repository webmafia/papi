package fastapi

import (
	"fmt"
	"testing"
)

func Example() {
	var r Router[struct{}]
	cb := func(ctx *Ctx[struct{}]) error { return nil }

	r.Add("GET", "/foo/bar/{baz}", cb)
	// r.Add("/foo/bar/mjau", unsafe.Pointer(&v))
	var params Params

	res := r.LookupString("GET", "/foo/bar/456", &params)

	if res == nil {
		fmt.Println("not found")
	} else {
		fmt.Println(params.Get("baz"))
	}

	// Output: Mjaa
}

func ExampleRouter_params() {
	var r Router[struct{}]
	cb := func(ctx *Ctx[struct{}]) error { return nil }

	r.Add("GET", "/foo/{bar}/{baz}/mjau", cb, func(s string) {
		fmt.Println("param:", s)
	})

	// Output: Todo
}

func BenchmarkAdd(b *testing.B) {
	var r Router[struct{}]
	cb := func(ctx *Ctx[struct{}]) error { return nil }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Add("GET", "/foo/bar/{baz}", cb)
		r.Clear()
	}
}

func BenchmarkLookup(b *testing.B) {
	var r Router[struct{}]
	cb := func(ctx *Ctx[struct{}]) error { return nil }

	r.Add("GET", "/foo/bar/{baz}", cb)
	r.Add("GET", "/foo/bar/{baz}/mjau", cb)
	// r.Add("/foo/bar/mjau", cb)
	var params Params

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = r.LookupString("GET", "/foo/bar/123", &params)
		params.Reset()
	}
}
