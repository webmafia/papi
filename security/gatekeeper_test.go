package security

import (
	"fmt"
	"testing"
)

func ExampleGatekeeper() {
	s, err := GenerateSecret()

	if err != nil {
		panic(err)
	}

	g, err := NewGatekeeper(s)

	if err != nil {
		panic(err)
	}

	tok, err := g.CreateToken()

	if err != nil {
		panic(err)
	}

	fmt.Println(tok)
	fmt.Println(g.ValidateToken(tok))
	fmt.Println(tok.Payload())

	// Output: TODO
}

func BenchmarkGatekeeper(b *testing.B) {
	s, err := GenerateSecret()

	if err != nil {
		b.Fatal(err)
	}

	g, err := NewGatekeeper(s)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	b.Run("CreateToken", func(b *testing.B) {
		for range b.N {
			_, _ = g.CreateToken()
		}
	})

	b.Run("ValidateToken", func(b *testing.B) {
		tok, _ := g.CreateToken()
		b.ResetTimer()

		for range b.N {
			_ = g.ValidateToken(tok)
		}
	})
}
