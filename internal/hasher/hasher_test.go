package hasher

import (
	"fmt"
	"testing"
)

func ExampleHasher() {
	var h Hasher
	h.Reset()

	fmt.Println(h.Hash())
	h.WriteString("foobar")
	fmt.Println(h.Hash())

	// Output:
	//
	// 17241709254077376921
	// 11721187498075204345
}

func ExampleHash() {
	type Array struct {
		Title       string
		Description string
		Min         int
		Max         int
		Nullable    bool
		ReadOnly    bool
		WriteOnly   bool
		UniqueItems bool
	}

	fmt.Println(Hash(Array{Min: 1}))

	// Output: 8852669417329623389
}

func BenchmarkHasher(b *testing.B) {
	var h Hasher
	h.Reset()

	for i := range b.N {
		_ = i
		h.WriteString("hello")
		_ = h.Hash()
	}
}

func BenchmarkHash(b *testing.B) {
	type Array struct {
		Title       string
		Description string
		Min         int
		Max         int
		Nullable    bool
		ReadOnly    bool
		WriteOnly   bool
		UniqueItems bool
	}

	var h Hasher
	h.Reset()

	for range b.N {
		_ = Hash(Array{Min: 123})
	}
}
