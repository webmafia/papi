package security

import "fmt"

func ExampleSecret() {
	s, err := GenerateSecret()

	if err != nil {
		panic(err)
	}

	fmt.Println(s)

	// Output: TODO
}
