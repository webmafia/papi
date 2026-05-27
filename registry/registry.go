package registry

import (
	"iter"
	"reflect"

	"github.com/webmafia/papi/internal"
	"github.com/webmafia/papi/internal/scanner"
	"github.com/webmafia/papi/security"
)

type Registry struct {
	scanCache  map[reflect.Type]Parser
	desc       map[reflect.Type]TypeDescription
	scan       scanner.Creator
	gatekeeper security.Gatekeeper
	perms      internal.Set[security.Permission]
}

func NewRegistry(gatekeeper ...security.Gatekeeper) (r *Registry) {
	r = &Registry{
		scanCache: make(map[reflect.Type]Parser),
		desc:      make(map[reflect.Type]TypeDescription),
	}

	if len(gatekeeper) > 0 && gatekeeper[0] != nil {
		r.gatekeeper = gatekeeper[0]
	}

	r.scan = scanner.NewCreator(func(typ reflect.Type) (scan scanner.Scanner, err error) {

		// 1. If there is an explicit registered parser, use it
		if desc, ok := r.desc[typ]; ok && desc.Parser != nil {
			return desc.Parser("")
		}

		// 2. If the type can describe itself, let it
		if typ.Implements(typeDescriber) {
			if v, ok := reflect.New(typ).Interface().(TypeDescriber); ok {
				if desc := v.TypeDescription(r); desc.Schema != nil {
					return desc.Parser("")
				}
			}
		}

		return
	})

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

// Get the register's gatekeeper. Could be nil.
func (r *Registry) Gatekeeper() security.Gatekeeper {
	return r.gatekeeper
}

// Iterate through all unique permissions of all registred routes without any guaranteed order.
func (t *Registry) Permissions() iter.Seq[security.Permission] {
	return t.perms.Values()
}

// Returns whether the permission tag is optional or not on routes.
func (r *Registry) OptionalPermTag() bool {
	return r.gatekeeper == nil || r.gatekeeper.OptionalPermTag()
}

func (r *Registry) describe(typ reflect.Type) (desc TypeDescription, ok bool) {

	// 1. If there is an explicit registered decoder, use it
	if desc, ok = r.desc[typ]; ok {
		return
	}

	// 2. If the type can describe itself, let it
	if typ.Implements(typeDescriber) {
		if v, ok := reflect.New(typ).Interface().(TypeDescriber); ok {
			return v.TypeDescription(r), true
		}
	}

	return
}
