package registry

import (
	"reflect"

	"github.com/webmafia/papi/internal"
	"github.com/webmafia/papi/internal/scanner"
	"github.com/webmafia/papi/policy"
	"github.com/webmafia/papi/token"
)

type Registry struct {
	scanCache map[reflect.Type]Decoder
	desc      map[reflect.Type]TypeDescription
	json      *internal.JSONPool
	scan      scanner.Creator
	policies  *policy.Store
	tokGen    *token.Generator
}

func NewRegistry(tokGen *token.Generator, json *internal.JSONPool) (r *Registry, err error) {
	r = &Registry{
		scanCache: make(map[reflect.Type]Decoder),
		desc:      make(map[reflect.Type]TypeDescription),
		json:      json,
		policies:  policy.NewStore(json),
		tokGen:    tokGen,
	}

	r.scan = scanner.NewCreator(r.scanner)

	return
}

//go:inline
func (s *Registry) JSON() *internal.JSONPool {
	return s.json
}

//go:inline
func (r *Registry) Policies() *policy.Store {
	return r.policies
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
