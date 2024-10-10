package valid

import (
	"reflect"
	"regexp"
	"unsafe"
)

func appendPatternValidators(offset uintptr, typ reflect.Type, field string, s string) (validator, error) {
	switch kind := typ.Kind(); kind {

	case reflect.Array:
		return validArray(offset, typ, field, s, appendPatternValidators)

	case reflect.Slice:
		return validSlice(offset, typ, field, s, appendPatternValidators)

	case reflect.Pointer:
		return validPointer(offset, typ, field, s, appendPatternValidators)

	case reflect.String:
		return validStringPattern(offset, field, s)

	default:
		return nil, notImplemented("pattern", kind)
	}
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
