package valid

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/webmafia/papi/errors"
)

// Generic function to test pattern validation for any type.
func testPatternValidator[T any](t *testing.T, value T, pattern string, expectedErr bool) {
	validator, err := createPatternValidator(0, reflect.TypeOf(value), "testField", pattern)

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

func Test_createPatternValidator(t *testing.T) {
	// Test pattern matching for strings
	t.Run("StringPattern", func(t *testing.T) {
		pattern := "^[a-z]+$"
		testPatternValidator(t, "hello", pattern, false)   // Matches pattern
		testPatternValidator(t, "Hello123", pattern, true) // Does not match pattern
		testPatternValidator(t, "", pattern, false)        // Empty string should pass
	})

	// Test pattern matching for arrays of strings
	t.Run("ArrayPattern", func(t *testing.T) {
		pattern := "^[a-z]+$"
		validArray := [3]string{"hello", "world", "test"}
		invalidArray := [3]string{"hello", "WORLD", "123"}
		testPatternValidator(t, validArray, pattern, false)  // All elements match
		testPatternValidator(t, invalidArray, pattern, true) // One element does not match
	})

	// Test pattern matching for slices of strings
	t.Run("SlicePattern", func(t *testing.T) {
		pattern := "^[a-z]+$"
		validSlice := []string{"hello", "world", "test"}
		invalidSlice := []string{"hello", "WORLD", "123"}
		testPatternValidator(t, validSlice, pattern, false)  // All elements match
		testPatternValidator(t, invalidSlice, pattern, true) // One element does not match
	})

	// Test pattern matching for pointers to strings
	t.Run("PointerPattern", func(t *testing.T) {
		pattern := "^[a-z]+$"
		valValid := "hello"
		valInvalid := "Hello123"
		ptrValid := &valValid
		ptrInvalid := &valInvalid
		testPatternValidator(t, ptrValid, pattern, false)  // Matches pattern
		testPatternValidator(t, ptrInvalid, pattern, true) // Does not match pattern
	})
}
