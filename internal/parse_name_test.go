package internal

import (
	"fmt"
	"testing"
)

func ExampleParseName() {
	str := "FooBarBazASD_haha"

	fmt.Println(ParseName(str))

	// Output: Foo bar baz ASD haha foo-bar-baz-asd-haha
}

func BenchmarkParseName(b *testing.B) {
	str := "FooBarBazASD_haha"

	b.ResetTimer()

	for range b.N {
		_, _ = ParseName(str)
	}
}

func Example_calcAllow() {
	str := []byte("FooBAR")

	fmt.Println(calcAlloc(str))
	// Output: 7
}

func Benchmark_calcAlloc(b *testing.B) {
	str := []byte("fooBar")

	b.ResetTimer()

	for range b.N {
		_ = calcAlloc(str)
	}
}

func ExampleCallerName() {
	fmt.Println(CallerName(0))

	// Output: ExampleCallerName
}
