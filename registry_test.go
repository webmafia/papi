package fastapi

import (
	"fmt"
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/fastapi/pool/json"
	"github.com/webmafia/fastapi/registry"
	"github.com/webmafia/fastapi/registry/request"
)

func ExampleRegistry() {
	json := json.NewPool(jsoniter.ConfigFastest)
	r := registry.NewRegistry(func(r *registry.Registry) (creator registry.RequestScannerCreator) {
		creator, err := request.NewRequestScanner(r, json)

		if err != nil {
			panic(err)
		}
		return
	})

	type req struct {
		Id    int `param:"id"`
		Limit int `query:"limit"`
	}

	s, err := r.Schema(reflect.TypeOf((*req)(nil)).Elem())

	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", *s)

	// Output: TODO
}