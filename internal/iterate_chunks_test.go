package internal

import (
	"fmt"
	"testing"
)

func ExampleIterateChunks() {
	var i int
	for s := range IterateChunks("123,456,789,,,,,", ',') {
		fmt.Println(s)
		i++
	}

	fmt.Println("count:", i)

	// Output:
	//
	// 123
	// 456
	// 789
	// count: 3
}

func BenchmarkIterateChunks(b *testing.B) {
	for range b.N {
		for s := range IterateChunks("123,456,789", ',') {
			_ = s
		}
	}
}
