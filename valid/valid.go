package valid

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type String struct {
	Enum      []string `tag:"enum"`
	Format    string   `tag:"format"`
	Pattern   string   `tag:"pattern"`
	Min       int      `tag:"min"`
	Max       int      `tag:"max"`
	Nullable  bool     `tag:"flags:nullable"`
	ReadOnly  bool     `tag:"flags:readonly"`
	WriteOnly bool     `tag:"flags:writeonly"`
	Required  bool     `tag:"flags:required"`
}

type validators []validator

func (valids *validators) append(valid validator, err error) error {
	if err == nil {
		*valids = append(*valids, valid)
	}

	return err
}

func (valids *validators) compile() (valid validator, err error) {
	*valids = compactSlice(*valids)

	return func(ptr unsafe.Pointer) error {
		for _, valid := range *valids {
			if err := valid(ptr); err != nil {
				return err
			}
		}

		return nil
	}, nil
}

func createStructValidator(typ reflect.Type) (valid validator, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	var valids validators

	if err = appendValidators(&valids, 0, typ, ""); err != nil {
		return
	}

	return valids.compile()
}

func appendValidators(valids *validators, offset uintptr, name string, typ reflect.Type, tags reflect.StructTag) (err error) {
	switch kind := typ.Kind(); kind {

	// case reflect.Bool:

	case reflect.Int:
		return valids.append(validNumMin[int](offset, name))

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
