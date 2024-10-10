package valid

import (
	"reflect"
	"regexp"
	"unsafe"
)

func appendPatternValidators(valids *validators, offset uintptr, fld *reflect.StructField, s string) (err error) {
	switch kind := fld.Type.Kind(); kind {

	// case reflect.Array:
	// case reflect.Slice:

	case reflect.String:
		return valids.append(validStringPattern(offset, fld.Name, s))
	}

	return
}

func validStringPattern(offset uintptr, field string, s string) (validator, error) {
	pattern, err := regexp.Compile(s)

	if err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		if val := *(*string)(unsafe.Add(ptr, offset)); val != "" && !pattern.MatchString(val) {
			errs.Append(FieldError{
				err:    ErrInvalidPattern,
				field:  field,
				expect: s,
			})
		}
	}, nil
}
