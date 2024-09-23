package strings

import (
	"fmt"
	"reflect"
	"sync"
)

type Factory struct {
	types sync.Map
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Scanner(typ reflect.Type) (scan Scanner, err error) {
	if val, ok := f.types.Load(typ); ok {
		if scan, ok := val.(Scanner); ok {
			return scan, nil
		}
	}

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
		return f.createArrayScanner(typ)

	case reflect.Pointer:
		return f.createPointerScanner(typ)

	case reflect.Slice:
		return f.createSliceScanner(typ)

	case reflect.String:
		return scanString, nil

	// case reflect.Struct:
	// 	return createStructScanner(typ)

	default:
		return nil, fmt.Errorf("cannot scan to type: %s", kind.String())
	}
}

func (f *Factory) RegisterScanner(typ reflect.Type, scan Scanner) {
	if scan == nil {
		f.types.Delete(typ)
		return
	}

	f.types.Store(typ, scan)
}
