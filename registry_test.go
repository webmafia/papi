package papi

import (
	"fmt"
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/webbmaffian/papi/pool/json"
	"github.com/webbmaffian/papi/registry"
)

func ExampleRegistry() {
	r, err := registry.NewRegistry(json.NewPool(jsoniter.ConfigFastest))

	if err != nil {
		return
	}

	type req struct {
		Id    int `param:"id"`
		Limit int `query:"limit"`
	}

	s, err := r.Schema(reflect.TypeOf((*req)(nil)).Elem())

	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", s)

	// Output: TODO
}
