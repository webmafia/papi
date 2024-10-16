package internal

import "testing"

// Structs for testing
type StructA struct {
	Field1 int
	Field2 string
}

type StructB struct {
	Field1 int
	Field2 string
}

type StructD struct {
	Field1 int
	Field2 int
}

type NonStructType int

func TestEqualStructs(t *testing.T) {

	// Test case 1: StructA and StructB are identical in structure
	if !EqualStructs[StructA, StructB]() {
		t.Errorf("Expected StructA and StructB to be identical")
	}

	// Test case 2: StructA and StructD differ by a field type (string vs int for Field2)
	if EqualStructs[StructA, StructD]() {
		t.Errorf("Expected StructA and StructD to be different")
	}

	// Test case 3: Non-struct types should return false (StructA vs NonStructType)
	if EqualStructs[StructA, NonStructType]() {
		t.Errorf("Expected StructA and NonStructType to be different")
	}

	// Test case 4: Identical struct compared with itself
	if !EqualStructs[StructA, StructA]() {
		t.Errorf("Expected StructA to be identical to itself")
	}
}
