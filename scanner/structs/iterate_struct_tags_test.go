package structs

import (
	"fmt"
	"reflect"
)

func Example_iterateStructTags() {
	tags := reflect.StructTag(`json:"foo" db:"bar" number:"123"`)

	for k, v := range iterateStructTags(tags) {
		fmt.Println(k, "is set to:", v)
	}

	// Output:
	//
	// json is set to: foo
	// db is set to: bar
	// number is set to: 123
}