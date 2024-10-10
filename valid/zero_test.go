package valid

import (
	"reflect"
	"testing"
	"unsafe"
)

func Test_createZeroChecker(t *testing.T) {

	t.Run("Int", func(t *testing.T) {
		testType(t, int(0), true)
		testType(t, int(123), false)
	})

	t.Run("Int8", func(t *testing.T) {
		testType(t, int8(0), true)
		testType(t, int8(127), false)
	})

	t.Run("Int16", func(t *testing.T) {
		testType(t, int16(0), true)
		testType(t, int16(256), false)
	})

	t.Run("Int32", func(t *testing.T) {
		testType(t, int32(0), true)
		testType(t, int32(12345), false)
	})

	t.Run("Int64", func(t *testing.T) {
		testType(t, int64(0), true)
		testType(t, int64(123456789), false)
	})

	t.Run("Uint", func(t *testing.T) {
		testType(t, uint(0), true)
		testType(t, uint(123), false)
	})

	t.Run("Uint8", func(t *testing.T) {
		testType(t, uint8(0), true)
		testType(t, uint8(255), false)
	})

	t.Run("Uint16", func(t *testing.T) {
		testType(t, uint16(0), true)
		testType(t, uint16(12345), false)
	})

	t.Run("Uint32", func(t *testing.T) {
		testType(t, uint32(0), true)
		testType(t, uint32(123456789), false)
	})

	t.Run("Uint64", func(t *testing.T) {
		testType(t, uint64(0), true)
		testType(t, uint64(1234567890123456789), false)
	})

	t.Run("Float32", func(t *testing.T) {
		testType(t, float32(0), true)
		testType(t, float32(123.45), false)
	})

	t.Run("Float64", func(t *testing.T) {
		testType(t, float64(0), true)
		testType(t, float64(123456789.12345), false)
	})

	t.Run("String", func(t *testing.T) {
		testType(t, "", true)
		testType(t, "non-zero", false)
	})

	t.Run("Pointer", func(t *testing.T) {
		var p *int
		testType(t, p, true)
		nonZeroPtr := new(int)
		testType(t, nonZeroPtr, false)
	})

	t.Run("Slice", func(t *testing.T) {
		var s []int
		testType(t, s, true)
		nonZeroSlice := []int{1, 2, 3}
		testType(t, nonZeroSlice, false)
	})

	t.Run("Array", func(t *testing.T) {
		var a [3]int
		testType(t, a, true)
		nonZeroArray := [3]int{1, 0, 0}
		testType(t, nonZeroArray, false)
	})
}

func testType[T any](t *testing.T, v T, expected bool) {
	check, err := createZeroChecker(reflect.TypeOf(v))

	if err != nil {
		t.Fatal(err)
	}

	if result := check(unsafe.Pointer(&v)); result != expected {
		t.Errorf("zero-checker for value %v: expected %v, got %v", v, expected, result)
	}
}

func Benchmark_createZeroChecker(b *testing.B) {
	b.Run("Int", func(b *testing.B) {
		benchmarkType(b, int(0))
		benchmarkType(b, int(123))
	})

	b.Run("Int8", func(b *testing.B) {
		benchmarkType(b, int8(0))
		benchmarkType(b, int8(127))
	})

	b.Run("Int16", func(b *testing.B) {
		benchmarkType(b, int16(0))
		benchmarkType(b, int16(256))
	})

	b.Run("Int32", func(b *testing.B) {
		benchmarkType(b, int32(0))
		benchmarkType(b, int32(12345))
	})

	b.Run("Int64", func(b *testing.B) {
		benchmarkType(b, int64(0))
		benchmarkType(b, int64(123456789))
	})

	b.Run("Uint", func(b *testing.B) {
		benchmarkType(b, uint(0))
		benchmarkType(b, uint(123))
	})

	b.Run("Uint8", func(b *testing.B) {
		benchmarkType(b, uint8(0))
		benchmarkType(b, uint8(255))
	})

	b.Run("Uint16", func(b *testing.B) {
		benchmarkType(b, uint16(0))
		benchmarkType(b, uint16(12345))
	})

	b.Run("Uint32", func(b *testing.B) {
		benchmarkType(b, uint32(0))
		benchmarkType(b, uint32(123456789))
	})

	b.Run("Uint64", func(b *testing.B) {
		benchmarkType(b, uint64(0))
		benchmarkType(b, uint64(1234567890123456789))
	})

	b.Run("Float32", func(b *testing.B) {
		benchmarkType(b, float32(0))
		benchmarkType(b, float32(123.45))
	})

	b.Run("Float64", func(b *testing.B) {
		benchmarkType(b, float64(0))
		benchmarkType(b, float64(123456789.12345))
	})

	b.Run("String", func(b *testing.B) {
		benchmarkType(b, "")
		benchmarkType(b, "non-zero")
	})

	b.Run("Pointer", func(b *testing.B) {
		var p *int
		nonZeroPtr := new(int)
		benchmarkType(b, p)
		benchmarkType(b, nonZeroPtr)
	})

	b.Run("Slice", func(b *testing.B) {
		var s []int
		nonZeroSlice := []int{1, 2, 3}
		benchmarkType(b, s)
		benchmarkType(b, nonZeroSlice)
	})

	b.Run("Array", func(b *testing.B) {
		var a [3]int
		nonZeroArray := [3]int{1, 0, 0}
		benchmarkType(b, a)
		benchmarkType(b, nonZeroArray)
	})
}

func benchmarkType[T any](b *testing.B, v T) {
	check, err := createZeroChecker(reflect.TypeOf(v))

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		check(unsafe.Pointer(&v))
	}
}
