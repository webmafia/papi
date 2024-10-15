package internal

import "fmt"

func ExampleIterateFlags() {
	for flag := range IterateFlags("foo,bar,baz") {
		fmt.Println(flag)
	}

	// Output:
	//
	// foo
	// bar
	// baz
}
