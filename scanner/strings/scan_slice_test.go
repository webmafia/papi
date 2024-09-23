package strings

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/webmafia/fast"
)

func Example_createSliceScanner() {
	var ints []int

	f := NewFactory()
	scan, err := f.createSliceScanner(reflect.TypeOf(ints))

	if err != nil {
		panic(err)
	}

	if err = scan(unsafe.Pointer(&ints), "123,456,789"); err != nil {
		panic(err)
	}

	fmt.Println("cap:", cap(ints))
	fmt.Println("len:", len(ints))
	fmt.Println("values:", ints)

	// Output:
	//
	// cap: 3
	// len: 3
	// values: [123 456 789]
}

func Benchmark_createSliceScanner(b *testing.B) {
	var ints []int

	f := NewFactory()
	scan, err := f.createSliceScanner(reflect.TypeOf(ints))

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
