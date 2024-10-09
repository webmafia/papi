package registry

import (
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/internal"
	"github.com/webbmaffian/papi/registry/structs"
)

func ScanTags[T any](reg *Registry, dst *T, tags reflect.StructTag) (err error) {
	typ := internal.ReflectType[T]()
	scan, ok := reg.tag[typ]

	if !ok {
		if scan, err = structs.CreateTagScanner(typ, reg.CreateValueScanner); err != nil {
			return
		}

		reg.tag[typ] = scan
	}

	return scan(unsafe.Pointer(dst), string(tags))
}

func scanSchemaTags[T any](reg *Registry, dst *T, tags reflect.StructTag) (*T, error) {
	if err := ScanTags(reg, dst, tags); err != nil {
		return nil, err
	}

	return dst, nil
}
