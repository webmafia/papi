package internal

import "reflect"

func EqualStructs[A, B any]() bool {
	return equalStructs(reflect.TypeFor[A](), reflect.TypeFor[B]())
}

func equalStructs(t1, t2 reflect.Type) bool {

	// Check if both types are structs
	if t1.Kind() != reflect.Struct || t2.Kind() != reflect.Struct {
		return false
	}

	// Check if both have the same number of fields
	if t1.NumField() != t2.NumField() {
		return false
	}

	// Check if each field is identical in type
	for i := 0; i < t1.NumField(); i++ {
		field1 := t1.Field(i)
		field2 := t2.Field(i)

		// Compare field types
		if field1.Type != field2.Type {
			return false
		}
	}

	// If all fields match in name and type, the structs are identical in structure
	return true
}
