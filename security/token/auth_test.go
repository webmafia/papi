package token

import (
	"fmt"
	"testing"
)

func Example_auth() {
	s, err := GenerateSecret()

	if err != nil {
		panic(err)
	}

	a := auth{secret: s}
	tok, err := a.CreateToken()

	if err != nil {
		panic(err)
	}

	fmt.Println(tok)
	fmt.Println(a.ValidateToken(tok))
	fmt.Println(tok.Payload())

	// Output: TODO
}

func Benchmark_auth(b *testing.B) {
	s, err := GenerateSecret()

	if err != nil {
		b.Fatal(err)
	}

	a := auth{secret: s}
	b.ResetTimer()

	b.Run("CreateToken", func(b *testing.B) {
		for range b.N {
			_, _ = a.CreateToken()
		}
	})

	b.Run("ValidateToken", func(b *testing.B) {
		tok, _ := a.CreateToken()
		b.ResetTimer()

		for range b.N {
			_ = a.ValidateToken(tok)
		}
	})
}
