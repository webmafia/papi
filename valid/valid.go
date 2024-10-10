package valid

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/registry/structs"
)

type validators []validator

func (valids *validators) append(valid validator) {
	*valids = append(*valids, valid)
}

func (valids *validators) compile() (valid validator, err error) {
	*valids = compactSlice(*valids)

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		for _, valid := range *valids {
			valid(ptr, errs)
		}
	}, nil
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
			if valid, err = appendMinValidators(offset, typ, field, v); err != nil {
				return
			}

			valids.append(valid)

		case "max":
			if valid, err = appendMaxValidators(offset, typ, field, v); err != nil {
				return
			}

			valids.append(valid)

		case "enum":
			if valid, err = appendEnumValidators(offset, typ, field, v); err != nil {
				return
			}

			valids.append(valid)

		case "pattern":
			if valid, err = appendPatternValidators(offset, typ, field, v); err != nil {
				return
			}

			valids.append(valid)

		// case "unique":

		case "flags":
			if structs.HasFlag(v, "required") {
				if valid, err = appendRequiredValidators(offset, typ, field); err != nil {
					return
				}

				valids.append(valid)
			}
		}
	}

	return
}
