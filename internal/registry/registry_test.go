package registry

import (
	"fmt"
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/webbmaffian/papi/internal"
)

func ExampleRegistry() {
	r, err := NewRegistry(internal.NewJSONPool(jsoniter.ConfigFastest))

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
