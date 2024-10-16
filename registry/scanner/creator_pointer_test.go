package scanner

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func Example_createPointerScanner() {
	var foo *int

	c := NewCreator()
	scan, err := c.createPointerScanner(reflect.TypeOf(foo))

	if err != nil {
		panic(err)
	}

	// Pointer to a pointer
	if err = scan(unsafe.Pointer(&foo), "123"); err != nil {
		panic(err)
	}

	fmt.Println(*foo)

	// Pointer to a pointer
	if err = scan(unsafe.Pointer(&foo), "456"); err != nil {
		panic(err)
	}

	fmt.Println(*foo)

	// Output:
	//
	// 123
	// 456
}

func Test_createPointerScanner(t *testing.T) {
	var foo *int
	c := NewCreator()
	scan, err := c.createPointerScanner(reflect.TypeOf(foo))

	if err != nil {
		panic(err)
	}

	if uintptr(unsafe.Pointer(foo)) != 0 {
		t.Fatal("expected nil pointer")
	}

	// Pointer to a pointer
	if err = scan(unsafe.Pointer(&foo), "123"); err != nil {
		panic(err)
	}

	ptr := uintptr(unsafe.Pointer(foo))

	if ptr == 0 {
		t.Fatal("expected non-nil pointer")
	}

	if *foo != 123 {
		t.Fatal("expected value 123")
	}

	// Pointer to a pointer
	if err = scan(unsafe.Pointer(&foo), "456"); err != nil {
		panic(err)
	}

	if uintptr(unsafe.Pointer(foo)) != ptr {
		t.Fatal("expected unchanged pointer")
	}

	if *foo != 456 {
		t.Fatal("expected value 456")
	}
}
