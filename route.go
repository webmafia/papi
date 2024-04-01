package fastapi

import (
	"errors"
	"reflect"
)

func (api *API[U]) RegisterRoutes(types ...any) (err error) {
	for i := range types {
		val := reflect.ValueOf(types[i])
		numMethods := val.NumMethod()

		for i := 0; i < numMethods; i++ {
			cb, ok := val.Method(i).Interface().(func(api *API[U]))

			if !ok {
				return errors.New("invalid handler")
			}

			cb(api)
		}
	}

	return
}

func AddRoute[U, I, O any](api *API[U], r Route[U, I, O]) (err error) {
	// TODO: Create input scanner and handle output.
	return
}
