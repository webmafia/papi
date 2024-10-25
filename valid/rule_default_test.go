package valid

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/webmafia/papi/errors"
)

// Utility function to test the default validator
func testDefaultValidator[T any](t *testing.T, value T, s string, expectedErr bool, expectedValue T) {
	validator, err := createDefaultValidator(0, reflect.TypeOf(value), "testField", s)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var fieldErrors errors.Errors
	validator(unsafe.Pointer(&value), &fieldErrors)

	// Check if the expected error state matches
	if expectedErr && !fieldErrors.HasError() {
		t.Log(fieldErrors)
		t.Errorf("expected error but got none for value: %v", value)
	} else if !expectedErr && fieldErrors.HasError() {
		t.Log(fieldErrors)
		t.Errorf("did not expect error but got one for value: %v", value)
	}

	// Check if the value was correctly set to the default or remains unchanged
	if !reflect.DeepEqual(value, expectedValue) {
		t.Errorf("expected value: %v, got: %v", expectedValue, value)
	}
}

func Test_createDefaultValidator(t *testing.T) {
	// Test default validation for integers
	t.Run("IntDefault", func(t *testing.T) {
		var zeroValue int
		nonZeroValue := 123

		// When zero, the default should be applied
		testDefaultValidator(t, zeroValue, "456", false, 456)

		// When non-zero, the default should not be applied
		testDefaultValidator(t, nonZeroValue, "456", false, 123)
	})

	// Test default validation for unsigned integers
	t.Run("UintDefault", func(t *testing.T) {
		var zeroValue uint
		nonZeroValue := uint(123)

		// When zero, the default should be applied
		testDefaultValidator(t, zeroValue, "789", false, 789)

		// When non-zero, the default should not be applied
		testDefaultValidator(t, nonZeroValue, "789", false, 123)
	})

	// Test default validation for floats
	t.Run("FloatDefault", func(t *testing.T) {
		var zeroValue float64
		nonZeroValue := 123.45

		// When zero, the default should be applied
		testDefaultValidator(t, zeroValue, "456.78", false, 456.78)

		// When non-zero, the default should not be applied
		testDefaultValidator(t, nonZeroValue, "456.78", false, 123.45)
	})

	// Test default validation for strings
	t.Run("StringDefault", func(t *testing.T) {
		var emptyString string
		nonEmptyString := "existing value"

		// When empty, the default should be applied
		testDefaultValidator(t, emptyString, "default", false, "default")

		// When non-empty, the default should not be applied
		testDefaultValidator(t, nonEmptyString, "default", false, "existing value")
	})

	// Test default validation for slices
	t.Run("SliceDefault", func(t *testing.T) {
		var nilSlice []int
		emptySlice := []int{}
		nonEmptySlice := []int{1, 2, 3}

		// When nil, the default should be applied
		testDefaultValidator(t, nilSlice, "4,5,6", false, []int{4, 5, 6})

		// When empty, the default should be applied
		testDefaultValidator(t, emptySlice, "4,5,6", false, []int{4, 5, 6})

		// When non-empty, the default should not be applied
		testDefaultValidator(t, nonEmptySlice, "4,5,6", false, []int{1, 2, 3})
	})

	// Test default validation for arrays
	t.Run("ArrayDefault", func(t *testing.T) {
		var zeroArray [3]int
		nonZeroArray := [3]int{1, 2, 3}

		// When all elements are zero, the default should be applied
		testDefaultValidator(t, zeroArray, "4,5,6", false, [3]int{4, 5, 6})

		// When some elements are non-zero, the default should not be applied
		testDefaultValidator(t, nonZeroArray, "4,5,6", false, [3]int{1, 2, 3})
	})

	// Test default validation for pointers
	t.Run("PointerDefault", func(t *testing.T) {
		var nilPtr *int
		nonNilPtr := new(int)
		*nonNilPtr = 123

		// When nil, the default should be applied
		testDefaultValidator(t, nilPtr, "456", false, func() *int {
			val := 456
			return &val
		}())

		// When non-nil, the default should not be applied
		testDefaultValidator(t, nonNilPtr, "456", false, nonNilPtr)
	})
}
