package registry

import (
	"reflect"

	"github.com/webbmaffian/papi/internal/scanner"
)

type Decoder = scanner.Scanner

func (r *Registry) Decoder(typ reflect.Type, tags reflect.StructTag) (dec Decoder, err error) {

	// 1. If there is an explicit registered decoder, use it
	if desc, ok := r.desc[typ]; ok && desc.Decoder != nil {
		return desc.Decoder(tags)
	}

	// 2. If the type can describe itself, let it
	if typ.Implements(typeDescriber) {
		if v, ok := reflect.New(typ).Interface().(TypeDescriber); ok {
			if desc := v.TypeDescription(r); desc.Schema != nil {
				return desc.Decoder(tags)
			}
		}
	}

	// 3. In all other cases, generate a decoder
	return r.scan.CreateScanner(typ)
}

func (r *Registry) scanner(typ reflect.Type) (scan scanner.Scanner, err error) {
	if desc, ok := r.desc[typ]; ok && desc.Decoder != nil {
		return desc.Decoder("")
	}

	return
}
