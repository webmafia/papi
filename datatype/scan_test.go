package datatype

import (
	"fmt"
	"testing"
)

func ExampleScan() {
	d := NewDataTypes()

	var b bool

	err := Scan(d, &b, "false")

	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", b)

	// Output: Mjaaa
}

func BenchmarkScan(b *testing.B) {
	d := NewDataTypes()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var foo bool
		err := Scan(d, &foo, "true")

		if err != nil {
			b.Fatal(err)
		}
	}
}
