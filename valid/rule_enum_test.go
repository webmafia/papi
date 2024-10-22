package valid

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/webmafia/papi/errors"
)

// Generic function to test enum validation for any type.
func testEnumValidator[T any](t *testing.T, value T, allowedValues string, expectedErr bool) {
	validator, err := createEnumValidator(0, reflect.TypeOf(value), "testField", allowedValues)

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

func Test_createEnumValidator(t *testing.T) {
	// Test integer types
	t.Run("IntEnum", func(t *testing.T) {
		allowedValues := "1,2,3"
		testEnumValidator(t, 3, allowedValues, false) // Value is in the enum
		testEnumValidator(t, 5, allowedValues, true)  // Value is not in the enum
		testEnumValidator(t, 0, allowedValues, false) // Zero value, should pass
	})

	// Test float types
	t.Run("Float32Enum", func(t *testing.T) {
		allowedValues := "1.5,2.5,3.5"
		testEnumValidator(t, float32(1.5), allowedValues, false) // Value is in the enum
		testEnumValidator(t, float32(4.5), allowedValues, true)  // Value is not in the enum
		testEnumValidator(t, float32(0), allowedValues, false)   // Zero value, should pass
	})

	// Test string types
	t.Run("StringEnum", func(t *testing.T) {
		allowedValues := "apple,banana,cherry"
		testEnumValidator(t, "apple", allowedValues, false) // Value is in the enum
		testEnumValidator(t, "grape", allowedValues, true)  // Value is not in the enum
		testEnumValidator(t, "", allowedValues, false)      // Zero value (empty string), should pass
	})

	// Test slice of ints
	t.Run("SliceEnum", func(t *testing.T) {
		allowedValues := "1,2,3"
		validSlice := []int{1, 2, 3}
		invalidSlice := []int{1, 2, 4}
		testEnumValidator(t, validSlice, allowedValues, false)  // All values are valid
		testEnumValidator(t, invalidSlice, allowedValues, true) // One invalid value
	})

	// Test array of ints
	t.Run("ArrayEnum", func(t *testing.T) {
		allowedValues := "1,2,3"
		validArray := [3]int{1, 2, 3}
		invalidArray := [3]int{1, 2, 4}
		testEnumValidator(t, validArray, allowedValues, false)  // All values are valid
		testEnumValidator(t, invalidArray, allowedValues, true) // One invalid value
	})

	// Test pointer to an int
	t.Run("PointerEnum", func(t *testing.T) {
		allowedValues := "1,2,3"
		valValid := 3
		valInvalid := 5
		ptrValid := &valValid
		ptrInvalid := &valInvalid
		testEnumValidator(t, ptrValid, allowedValues, false)  // Value is in the enum
		testEnumValidator(t, ptrInvalid, allowedValues, true) // Value is not in the enum
	})
}
