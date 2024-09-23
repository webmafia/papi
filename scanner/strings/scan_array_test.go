package strings

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/webmafia/fast"
)

func Example_createArrayScanner() {
	var ints [3]int

	f := NewFactory()
	scan, err := f.createArrayScanner(reflect.TypeOf(ints))

	if err != nil {
		panic(err)
	}

	if err = scan(unsafe.Pointer(&ints), "123,456,789,999,,,,"); err != nil {
		panic(err)
	}

	fmt.Println("values:", ints)

	// Output: values: [123 456 789]
}

func Benchmark_createArrayScanner(b *testing.B) {
	var ints [3]int

	f := NewFactory()
	scan, err := f.createArrayScanner(reflect.TypeOf(ints))

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for range b.N {
		var ints []int

		if err = scan(fast.Noescape(unsafe.Pointer(&ints)), "123,456,789"); err != nil {
			b.Fatal(err)
		}
	}
}
