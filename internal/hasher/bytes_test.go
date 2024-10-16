package hasher

import (
	"fmt"
	"testing"
)

func Example_toBytes() {
	var i int16 = 1234

	fmt.Println(i)
	fmt.Println(toBytes(&i))

	// Output:
	//
	// 1234
	// [210 4]
}

func Benchmark_toBytes(b *testing.B) {
	b.ResetTimer()

	for i := range b.N {
		_ = toBytes(&i)
	}
}
