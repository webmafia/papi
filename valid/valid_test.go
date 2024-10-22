package valid

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/webmafia/papi/errors"
)

// Define a complex struct with nested structs, slices, and pointers
type Address struct {
	City    string `pattern:"^[A-Za-z ]+$" flags:"required"`
	ZipCode string `pattern:"^[0-9]{5}$" flags:"required"`
}

type Contact struct {
	Email string  `pattern:"^[^@]+@[^@]+\\.[^@]+$" flags:"required"`
	Phone *string `pattern:"^[0-9]{10}$"`
}

type Person struct {
	Name     string    `pattern:"^[A-Za-z]+$" flags:"required"`
	Age      int       `min:"18" max:"100"`
	Address  Address   `flags:"required"`
	Contacts []Contact `min:"1" flags:"required"`
	Notes    *string   `max:"500"`
}

// Helper function to test struct validation
func testStructValidator[T any](t *testing.T, value T, expectedErr bool, testCaseName string) {
	validator, err := createStructValidator(reflect.TypeOf(value))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var fieldErrors errors.Errors
	validator(unsafe.Pointer(&value), &fieldErrors)

	if expectedErr && !fieldErrors.HasError() {
		t.Errorf("expected error but got none for test case %s with struct: %+v", testCaseName, value)
	} else if !expectedErr && fieldErrors.HasError() {
		t.Errorf("did not expect error but got one for test case %s with struct: %+v", testCaseName, value)
	}
}

func Test_createStructValidator_Complex(t *testing.T) {
	// Valid phone and note pointers
	phoneValid := "1234567890"
	noteValid := "This is a note."

	// Valid struct instance
	t.Run("ValidStruct", func(t *testing.T) {
		validStruct := Person{
			Name: "Alice",
			Age:  25,
			Address: Address{
				City:    "Los Angeles",
				ZipCode: "90210",
			},
			Contacts: []Contact{
				{Email: "alice@example.com", Phone: &phoneValid},
			},
			Notes: &noteValid,
		}
		testStructValidator(t, validStruct, false, "ValidStruct")
	})

	// Struct with invalid name (does not match pattern)
	t.Run("InvalidName", func(t *testing.T) {
		invalidStruct := Person{
			Name: "Alice123", // Invalid pattern
			Age:  30,
			Address: Address{
				City:    "Los Angeles",
				ZipCode: "90210",
			},
			Contacts: []Contact{
				{Email: "alice@example.com", Phone: &phoneValid},
			},
			Notes: &noteValid,
		}
		testStructValidator(t, invalidStruct, true, "InvalidName")
	})

	// Struct with missing required address (empty city)
	t.Run("MissingCity", func(t *testing.T) {
		invalidStruct := Person{
			Name: "Bob",
			Age:  30,
			Address: Address{
				City:    "", // Empty city should fail
				ZipCode: "12345",
			},
			Contacts: []Contact{
				{Email: "bob@example.com", Phone: &phoneValid},
			},
			Notes: &noteValid,
		}
		testStructValidator(t, invalidStruct, true, "MissingCity")
	})

	// Struct with invalid ZipCode pattern
	t.Run("InvalidZipCode", func(t *testing.T) {
		invalidStruct := Person{
			Name: "Charlie",
			Age:  30,
			Address: Address{
				City:    "San Francisco",
				ZipCode: "ABCDE", // Invalid ZipCode
			},
			Contacts: []Contact{
				{Email: "charlie@example.com", Phone: &phoneValid},
			},
			Notes: &noteValid,
		}
		testStructValidator(t, invalidStruct, true, "InvalidZipCode")
	})

	// Struct with missing email (required)
	t.Run("MissingContactEmail", func(t *testing.T) {
		invalidStruct := Person{
			Name: "Dave",
			Age:  40,
			Address: Address{
				City:    "Boston",
				ZipCode: "02108",
			},
			Contacts: []Contact{
				{Email: "", Phone: &phoneValid}, // Missing required email
			},
			Notes: &noteValid,
		}
		testStructValidator(t, invalidStruct, true, "MissingContactEmail")
	})

	// Struct with invalid slice length for Contacts (below min)
	t.Run("EmptyContacts", func(t *testing.T) {
		invalidStruct := Person{
			Name: "Eve",
			Age:  35,
			Address: Address{
				City:    "New York",
				ZipCode: "10001",
			},
			Contacts: []Contact{}, // No contacts provided
			Notes:    &noteValid,
		}
		testStructValidator(t, invalidStruct, true, "EmptyContacts")
	})

	// Struct with invalid pointer (Notes too long)
	t.Run("NotesTooLong", func(t *testing.T) {
		invalidNote := "This note is way too long." + string(make([]byte, 501)) // 501 characters
		invalidStruct := Person{
			Name: "Frank",
			Age:  50,
			Address: Address{
				City:    "Miami",
				ZipCode: "33101",
			},
			Contacts: []Contact{
				{Email: "frank@example.com", Phone: &phoneValid},
			},
			Notes: &invalidNote, // Notes exceed 500 characters
		}
		testStructValidator(t, invalidStruct, true, "NotesTooLong")
	})
}

func Test_createStructValidator_Simple(t *testing.T) {
	// Valid phone and note pointers
	phoneValid := "1234567890"
	type foo struct {
		Contacts []Contact `min:"1" flags:"required"`
	}

	t.Run("MissingContactEmail", func(t *testing.T) {
		invalidStruct := Contact{
			Email: "", // Missing required email
			Phone: &phoneValid,
		}
		testStructValidator(t, invalidStruct, true, "MissingContactEmail")
	})

	t.Run("MissingContactEmail_Slice", func(t *testing.T) {
		invalidStruct := foo{
			Contacts: []Contact{
				{Email: "", Phone: &phoneValid}, // Missing required email
			},
		}
		testStructValidator(t, invalidStruct, true, "MissingContactEmail")
	})
}
