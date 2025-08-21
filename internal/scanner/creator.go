package scanner

import (
	"fmt"
	"reflect"
)

type Creator struct {
	custom func(typ reflect.Type) (Scanner, error)
}

func NewCreator(custom ...func(typ reflect.Type) (Scanner, error)) Creator {
	if len(custom) > 0 && custom[0] != nil {
		return Creator{
			custom: custom[0],
		}
	}

	return Creator{
		custom: func(typ reflect.Type) (Scanner, error) { return nil, nil },
	}
}

func (c Creator) CreateScanner(typ reflect.Type) (scan Scanner, err error) {
	if scan, err = c.custom(typ); scan != nil {
		return
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
		return c.createArrayScanner(typ)

	case reflect.Pointer:
		return c.createPointerScanner(typ)

	case reflect.Slice:
		return c.createSliceScanner(typ)

	case reflect.String:
		return scanString, nil

	// case reflect.Struct:
	// 	return createStructScanner(typ)

	default:
		return nil, fmt.Errorf("cannot create a scanner for '%s' of type %s", typ.Name(), kind.String())
	}
}
