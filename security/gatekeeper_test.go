package security

import (
	"context"
	"fmt"
	"testing"
)

func ExampleGatekeeper() {
	s, err := GenerateSecret()

	if err != nil {
		panic(err)
	}

	g, err := NewGatekeeper(s, dummyStore{})

	if err != nil {
		panic(err)
	}

	tok, err := g.CreateToken(context.Background(), nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(tok)
	fmt.Println(g.ValidateToken(context.Background(), tok))
	fmt.Println(tok.Payload())

	// Output: TODO
}

func BenchmarkGatekeeper(b *testing.B) {
	s, err := GenerateSecret()

	if err != nil {
		b.Fatal(err)
	}

	g, err := NewGatekeeper(s, dummyStore{})

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	b.Run("CreateToken", func(b *testing.B) {
		for range b.N {
			_, _ = g.CreateToken(context.Background(), nil)
		}
	})

	b.Run("ValidateToken", func(b *testing.B) {
		tok, _ := g.CreateToken(context.Background(), nil)
		b.ResetTimer()

		for range b.N {
			_, _ = g.ValidateToken(context.Background(), tok)
		}
	})
}
