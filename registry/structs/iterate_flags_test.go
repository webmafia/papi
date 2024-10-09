package structs

import "fmt"

func Example_iterateFlags() {
	for flag := range iterateFlags("foo,bar,baz") {
		fmt.Println(flag)
	}

	// Output:
	//
	// foo
	// bar
	// baz
}
