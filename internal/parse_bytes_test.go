package internal

import "fmt"

func ExampleParseBytes() {
	values := []string{"512", "1kb", "2MB", "1.5GiB", "3pb"}
	for _, v := range values {
		n, err := ParseBytes(v)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Println(n)
	}
	// Output:
	// 512
	// 1024
	// 2097152
	// 1610612736
	// 3377699720527872
}
