package valid

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/webbmaffian/papi/errors"
)

// Generic function to test min validation for any type.
func testMinValidator[T any](t *testing.T, value T, min string, expectedErr bool) {
	validator, err := createMinValidator(0, reflect.TypeOf(value), "testField", min)

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

func Test_createMinValidator(t *testing.T) {
	// Test integer types
	t.Run("IntMin", func(t *testing.T) {
		min := "10"
		testMinValidator(t, 5, min, true)   // Value is below the min
		testMinValidator(t, 10, min, false) // Value is at the min
		testMinValidator(t, 15, min, false) // Value is above the min
		testMinValidator(t, 0, min, false)  // Zero value, should pass
	})

	// Test unsigned integer types
	t.Run("UintMin", func(t *testing.T) {
		min := "10"
		testMinValidator(t, uint(5), min, true)   // Value is below the min
		testMinValidator(t, uint(10), min, false) // Value is at the min
		testMinValidator(t, uint(15), min, false) // Value is above the min
		testMinValidator(t, uint(0), min, false)  // Zero value, should pass
	})

	// Test float types
	t.Run("FloatMin", func(t *testing.T) {
		min := "10.5"
		testMinValidator(t, float32(5.0), min, true)   // Value is below the min
		testMinValidator(t, float32(10.5), min, false) // Value is at the min
		testMinValidator(t, float32(15.0), min, false) // Value is above the min
		testMinValidator(t, float32(0), min, false)    // Zero value, should pass
	})

	// Test slice min length
	t.Run("SliceMin", func(t *testing.T) {
		min := "3"
		testMinValidator(t, []int{1, 2}, min, true)        // Length is below the min
		testMinValidator(t, []int{1, 2, 3}, min, false)    // Length is at the min
		testMinValidator(t, []int{1, 2, 3, 4}, min, false) // Length is above the min
	})

	// Test string min length
	t.Run("StringMin", func(t *testing.T) {
		min := "5"
		testMinValidator(t, "abc", min, true)     // Length is below the min
		testMinValidator(t, "abcde", min, false)  // Length is at the min
		testMinValidator(t, "abcdef", min, false) // Length is above the min
		testMinValidator(t, "", min, false)       // Empty string, should pass
	})

	// Test pointer min validation
	t.Run("PointerMin", func(t *testing.T) {
		min := "10"
		valValid := 15
		valInvalid := 5
		ptrValid := &valValid
		ptrInvalid := &valInvalid
		testMinValidator(t, ptrValid, min, false)  // Value is above the min
		testMinValidator(t, ptrInvalid, min, true) // Value is below the min
	})
}
