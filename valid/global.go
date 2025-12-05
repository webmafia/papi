package valid

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/webmafia/papi/errors"
)

var (
	structValids   map[reflect.Type]structValidator
	structValidsMu sync.RWMutex
)

func GetStructValidator[T any]() (StructValidator[T], error) {
	valid, err := getStructValidator(reflect.TypeFor[T]())

	if err != nil {
		return nil, err
	}

	return *(*StructValidator[T])(unsafe.Pointer(&valid)), nil
}

func getStructValidator(typ reflect.Type) (valid structValidator, err error) {
	structValidsMu.RLock()
	if structValids != nil {
		valid = structValids[typ]
	}
	structValidsMu.RUnlock()

	if valid == nil {
		if valid, err = createStructValidator(typ); err != nil {
			return
		}

		structValidsMu.Lock()

		if structValids == nil {
			structValids = make(map[reflect.Type]structValidator)
		}

		structValids[typ] = valid

		structValidsMu.Unlock()
	}

	return
}

func ValidateStruct[T any](ptr *T, errs *errors.Errors) bool {
	typ := reflect.TypeFor[T]()

	var (
		fn  structValidator
		err error
	)

	if fn, err = getStructValidator(typ); err != nil {
		errs.Append(errors.NewError("INVALID_VALIDATOR", err.Error(), 500))
		return false
	}

	return fn(unsafe.Pointer(ptr), errs)
}
