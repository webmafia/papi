package valid

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/internal"
	"github.com/webbmaffian/papi/registry/structs"
)

func CreateStructValidator[T any]() (Validator[T], error) {
	valid, err := createStructValidator(internal.ReflectType[T]())

	if err != nil {
		return nil, err
	}

	return *(*Validator[T])(unsafe.Pointer(&valid)), nil
}

func createStructValidator(typ reflect.Type) (valid validator, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	var valids validators

	if err = appendStructValidators(&valids, typ, 0); err != nil {
		return
	}

	return valids.compile()
}

func appendStructValidators(valids *validators, typ reflect.Type, offset uintptr) (err error) {
	numFields := typ.NumField()

	for i := range numFields {
		fld := typ.Field(i)

		if !fld.IsExported() {
			continue
		}

		if err = appendFieldValidators(valids, fld.Type, offset+fld.Offset, fld.Name, fld.Tag); err != nil {
			return
		}

		if fld.Type.Kind() == reflect.Struct {
			if err = appendStructValidators(valids, fld.Type, offset+fld.Offset); err != nil {
				return
			}
		}
	}

	return
}

func appendFieldValidators(valids *validators, typ reflect.Type, offset uintptr, field string, tag reflect.StructTag) (err error) {
	for k, v := range structs.IterateStructTags(tag) {
		var valid validator

		switch k {

		case "min":
			if valid, err = createMinValidator(offset, typ, field, v); err != nil {
				return
			}

			valids.append(valid)

		case "max":
			if valid, err = createMaxValidator(offset, typ, field, v); err != nil {
				return
			}

			valids.append(valid)

		case "enum":
			if valid, err = createEnumValidator(offset, typ, field, v); err != nil {
				return
			}

			valids.append(valid)

		case "pattern":
			if valid, err = createPatternValidator(offset, typ, field, v); err != nil {
				return
			}

			valids.append(valid)

		// case "unique":

		case "flags":
			if structs.HasFlag(v, "required") {
				if valid, err = createRequiredValidator(offset, typ, field); err != nil {
					return
				}

				valids.append(valid)
			}
		}
	}

	return
}
