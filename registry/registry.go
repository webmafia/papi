package registry

import (
	"reflect"

	"github.com/webmafia/papi/internal/scanner"
	"github.com/webmafia/papi/security"
)

type Registry struct {
	scanCache      map[reflect.Type]Decoder
	desc           map[reflect.Type]TypeDescription
	scan           scanner.Creator
	securityScheme security.Scheme
	forcePermTag   bool
}

func NewRegistry(securityScheme security.Scheme, forcePermTag bool) (r *Registry, err error) {
	r = &Registry{
		scanCache:      make(map[reflect.Type]Decoder),
		desc:           make(map[reflect.Type]TypeDescription),
		securityScheme: securityScheme,
		forcePermTag:   forcePermTag,
	}

	r.scan = scanner.NewCreator(r.scanner)

	return
}

func (r *Registry) RegisterType(typs ...TypeRegistrar) (err error) {
	for _, typ := range typs {
		desc := typ.TypeDescription(r)

		if desc.IsZero() {
			delete(r.desc, typ.Type())
		} else {
			r.desc[typ.Type()] = desc
		}
	}

	return
}
