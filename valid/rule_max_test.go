package valid

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/webmafia/papi/errors"
)

func testMaxValidator[T any](t *testing.T, value T, max string, expectedErr bool) {
	validator, err := createMaxValidator(0, reflect.TypeOf(value), "testField", max)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var fieldErrors errors.Errors
	validator(unsafe.Pointer(&value), &fieldErrors)

	if expectedErr && !fieldErrors.HasError() {
		t.Errorf("expected error but got none for value: %v", value)
	} else if !expectedErr && fieldErrors.HasError() {
		t.Errorf("did not expect error but got one for value: %v", value)
	}
}

func Test_createMaxValidator(t *testing.T) {
	// Test integer types
	t.Run("IntMax", func(t *testing.T) {
		max := "100"
		testMaxValidator(t, 50, max, false)  // Value is below the max
		testMaxValidator(t, 100, max, false) // Value is at the max
		testMaxValidator(t, 150, max, true)  // Value is above the max
		testMaxValidator(t, 0, max, false)   // Zero value, should pass
	})

	// Test unsigned integer types
	t.Run("UintMax", func(t *testing.T) {
		max := "100"
		testMaxValidator(t, uint(50), max, false)  // Value is below the max
		testMaxValidator(t, uint(100), max, false) // Value is at the max
		testMaxValidator(t, uint(150), max, true)  // Value is above the max
		testMaxValidator(t, uint(0), max, false)   // Zero value, should pass
	})

	// Test float types
	t.Run("FloatMax", func(t *testing.T) {
		max := "100.5"
		testMaxValidator(t, float32(50.0), max, false)  // Value is below the max
		testMaxValidator(t, float32(100.5), max, false) // Value is at the max
		testMaxValidator(t, float32(150.0), max, true)  // Value is above the max
		testMaxValidator(t, float32(0), max, false)     // Zero value, should pass
	})

	// Test slice max length
	t.Run("SliceMax", func(t *testing.T) {
		max := "3"
		testMaxValidator(t, []int{1, 2}, max, false)      // Length is below the max
		testMaxValidator(t, []int{1, 2, 3}, max, false)   // Length is at the max
		testMaxValidator(t, []int{1, 2, 3, 4}, max, true) // Length is above the max
	})

	// Test string max length
	t.Run("StringMax", func(t *testing.T) {
		max := "5"
		testMaxValidator(t, "abc", max, false)   // Length is below the max
		testMaxValidator(t, "abcde", max, false) // Length is at the max
		testMaxValidator(t, "abcdef", max, true) // Length is above the max
		testMaxValidator(t, "", max, false)      // Empty string, should pass
	})

	// Test pointer max validation
	t.Run("PointerMax", func(t *testing.T) {
		max := "100"
		valValid := 50
		valInvalid := 150
		ptrValid := &valValid
		ptrInvalid := &valInvalid
		testMaxValidator(t, ptrValid, max, false)  // Value is below the max
		testMaxValidator(t, ptrInvalid, max, true) // Value is above the max
	})
}
