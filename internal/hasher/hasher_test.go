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

	// Output: TODO
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

	// Output: TODO
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
		Hash(Array{Min: 123})
	}
}
