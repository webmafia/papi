package scanner

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Scanner func(unsafe.Pointer, string) error
type CreateValueScanner func(typ reflect.Type, createElemScanner CreateValueScanner) (scan Scanner, err error)

func CreateScanner(typ reflect.Type) (scan Scanner, err error) {
	return CreateCustomScanner(typ, CreateCustomScanner)
}

func CreateCustomScanner(typ reflect.Type, createElemScanner CreateValueScanner) (scan Scanner, err error) {
	switch kind := typ.Kind(); kind {

	case reflect.Bool:
		return scanBool, nil

	case reflect.Int:
		return scanInt, nil

	case reflect.Int8:
		return scanInt8, nil

	case reflect.Int16:
		return scanInt16, nil

	case reflect.Int32:
		return scanInt32, nil

	case reflect.Int64:
		return scanInt64, nil

	case reflect.Uint:
		return scanUint, nil

	case reflect.Uint8:
		return scanUint8, nil

	case reflect.Uint16:
		return scanUint16, nil

	case reflect.Uint32:
		return scanUint32, nil

	case reflect.Uint64:
		return scanUint64, nil

	case reflect.Float32:
		return scanFloat32, nil

	case reflect.Float64:
		return scanFloat64, nil

	case reflect.Complex64:
		return scanComplex64, nil

	case reflect.Complex128:
		return scanComplex128, nil

	case reflect.Array:
		return createArrayScanner(typ, createElemScanner)

	case reflect.Pointer:
		return createPointerScanner(typ, createElemScanner)

	case reflect.Slice:
		return createSliceScanner(typ, createElemScanner)

	case reflect.String:
		return scanString, nil

	// case reflect.Struct:
	// 	return createStructScanner(typ)

	default:
		return nil, fmt.Errorf("cannot scan to type: %s", kind.String())
	}
}
