package valid

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/registry/structs"
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

		if err = appendFieldValidators(valids, offset, &fld); err != nil {
			return
		}
	}

	return
}

func appendFieldValidators(valids *validators, offset uintptr, fld *reflect.StructField) (err error) {
	for k, v := range structs.IterateStructTags(fld.Tag) {
		switch k {

		case "min":
			if err = appendMinValidators(valids, offset, fld, v); err != nil {
				return
			}

		case "max":
			if err = appendMaxValidators(valids, offset, fld, v); err != nil {
				return
			}

		case "enum":
			if err = appendEnumValidators(valids, offset, fld, v); err != nil {
				return
			}

		case "pattern":
			if err = appendPatternValidators(valids, offset, fld, v); err != nil {
				return
			}

		case "flags":
			if structs.HasFlag(v, "required") {
				if err = appendRequiredValidators(valids, offset, fld); err != nil {
					return
				}
			}
		}
	}

	return
}
