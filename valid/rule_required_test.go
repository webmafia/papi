package valid

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/webbmaffian/papi/errors"
)

func testRequiredValidator[T any](t *testing.T, value T, expectedErr bool) {
	validator, err := createRequiredValidator(0, reflect.TypeOf(value), "testField")

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

func Test_createRequiredValidator(t *testing.T) {
	// Test required validation for integers
	t.Run("IntRequired", func(t *testing.T) {
		testRequiredValidator(t, 0, true)    // Zero value, should fail
		testRequiredValidator(t, 123, false) // Non-zero value, should pass
	})

	// Test required validation for unsigned integers
	t.Run("UintRequired", func(t *testing.T) {
		testRequiredValidator(t, uint(0), true)    // Zero value, should fail
		testRequiredValidator(t, uint(123), false) // Non-zero value, should pass
	})

	// Test required validation for floats
	t.Run("FloatRequired", func(t *testing.T) {
		testRequiredValidator(t, 0.0, true)     // Zero value, should fail
		testRequiredValidator(t, 123.45, false) // Non-zero value, should pass
	})

	// Test required validation for strings
	t.Run("StringRequired", func(t *testing.T) {
		testRequiredValidator(t, "", true)           // Empty string, should fail
		testRequiredValidator(t, "non-empty", false) // Non-empty string, should pass
	})

	// Test required validation for slices
	t.Run("SliceRequired", func(t *testing.T) {
		testRequiredValidator(t, []int(nil), true)      // Nil slice, should fail
		testRequiredValidator(t, []int{}, true)         // Empty slice, should fail
		testRequiredValidator(t, []int{1, 2, 3}, false) // Non-empty slice, should pass
	})

	// Test required validation for arrays
	t.Run("ArrayRequired", func(t *testing.T) {
		testRequiredValidator(t, [3]int{}, true)         // Zero array, should fail
		testRequiredValidator(t, [3]int{1, 2, 3}, false) // Non-zero array, should pass
	})

	// Test required validation for pointers
	t.Run("PointerRequired", func(t *testing.T) {
		var ptr *int
		nonZeroPtr := new(int)
		*nonZeroPtr = 123
		testRequiredValidator(t, ptr, true)         // Nil pointer, should fail
		testRequiredValidator(t, nonZeroPtr, false) // Non-nil pointer, should pass
	})
}
