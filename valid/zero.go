package valid

import (
	"reflect"
	"unsafe"
)

func createZeroChecker(t reflect.Type) (func(ptr unsafe.Pointer) bool, error) {
	switch kind := t.Kind(); kind {
	case reflect.Int:
		return func(ptr unsafe.Pointer) bool {
			return *(*int)(ptr) == 0
		}, nil
	case reflect.Int8:
		return func(ptr unsafe.Pointer) bool {
			return *(*int8)(ptr) == 0
		}, nil
	case reflect.Int16:
		return func(ptr unsafe.Pointer) bool {
			return *(*int16)(ptr) == 0
		}, nil
	case reflect.Int32:
		return func(ptr unsafe.Pointer) bool {
			return *(*int32)(ptr) == 0
		}, nil
	case reflect.Int64:
		return func(ptr unsafe.Pointer) bool {
			return *(*int64)(ptr) == 0
		}, nil
	case reflect.Uint:
		return func(ptr unsafe.Pointer) bool {
			return *(*uint)(ptr) == 0
		}, nil
	case reflect.Uint8:
		return func(ptr unsafe.Pointer) bool {
			return *(*uint8)(ptr) == 0
		}, nil
	case reflect.Uint16:
		return func(ptr unsafe.Pointer) bool {
			return *(*uint16)(ptr) == 0
		}, nil
	case reflect.Uint32:
		return func(ptr unsafe.Pointer) bool {
			return *(*uint32)(ptr) == 0
		}, nil
	case reflect.Uint64:
		return func(ptr unsafe.Pointer) bool {
			return *(*uint64)(ptr) == 0
		}, nil
	case reflect.Float32:
		return func(ptr unsafe.Pointer) bool {
			return *(*float32)(ptr) == 0
		}, nil
	case reflect.Float64:
		return func(ptr unsafe.Pointer) bool {
			return *(*float64)(ptr) == 0
		}, nil
	case reflect.String:
		return func(ptr unsafe.Pointer) bool {
			return *(*string)(ptr) == ""
		}, nil
	case reflect.Pointer:
		return func(ptr unsafe.Pointer) bool {
			return *(*unsafe.Pointer)(ptr) == nil
		}, nil
	case reflect.Slice:
		return func(ptr unsafe.Pointer) bool {
			return sliceLen(ptr) == 0
		}, nil
	case reflect.Array:
		arrayLen := uintptr(t.Len())
		elemType := t.Elem()
		checkElem, err := createZeroChecker(elemType)

		if err != nil {
			return nil, err
		}

		return func(ptr unsafe.Pointer) bool {
			for i := range arrayLen {
				// Calculate the offset of the i-th element
				elemPtr := unsafe.Add(ptr, i*elemType.Size())
				if !checkElem(elemPtr) {
					return false
				}
			}
			return true
		}, nil
	default:
		return nil, notImplemented("required", kind)
	}
}
