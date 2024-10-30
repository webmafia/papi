package token

import (
	"fmt"
	"testing"
)

func ExampleGenerator() {
	var key [32]byte

	g, err := NewGenerator(key[:])

	if err != nil {
		panic(err)
	}

	tok, err := g.CreateToken(nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(tok)
	fmt.Println(g.ValidateToken(tok))
	fmt.Println(tok.Payload())
	view, _ := g.GetValidatedTokenView(tok.bytes())
	fmt.Println(view.Payload())

	// Output: TODO
}

func BenchmarkGenerator(b *testing.B) {
	var key [32]byte
	g, err := NewGenerator(key[:])

	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	b.Run("CreateToken", func(b *testing.B) {
		for range b.N {
			_, _ = g.CreateToken(nil)
		}
	})

	b.Run("ValidateToken", func(b *testing.B) {
		tok, _ := g.CreateToken(nil)
		b.ResetTimer()

		for range b.N {
			_ = g.ValidateToken(tok)
		}
	})

	b.Run("GetValidatedTokenView", func(b *testing.B) {
		tok, _ := g.CreateToken(nil)
		buf := tok.bytes()
		b.ResetTimer()

		for range b.N {
			if _, err := g.GetValidatedTokenView(buf); err != nil {
				b.Fatal(err)
			}
		}
	})
}
