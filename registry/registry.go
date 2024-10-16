package registry

import (
	"reflect"
	"sync"

	"github.com/webbmaffian/papi/pool/json"
	"github.com/webbmaffian/papi/registry/scanner"
)

type Registry struct {
	scanCache map[reflect.Type]Decoder
	desc      map[reflect.Type]TypeDescription
	json      *json.Pool
	scan      scanner.Creator
	mu        sync.RWMutex
}

func NewRegistry(json *json.Pool) (r *Registry, err error) {
	r = &Registry{
		scanCache: make(map[reflect.Type]Decoder),
		desc:      make(map[reflect.Type]TypeDescription),
		json:      json,
	}

	r.scan = scanner.NewCreator(r.scanner)

	return
}

//go:inline
func (s *Registry) JSON() *json.Pool {
	return s.json
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
