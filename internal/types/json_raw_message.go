package types

import (
	"encoding/json"
	"reflect"

	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
)

func JsonRawMessage() registry.TypeRegistrar {
	return jsonRawMessage{}
}

type jsonRawMessage struct{}

func (t jsonRawMessage) Type() reflect.Type {
	return reflect.TypeOf((*json.RawMessage)(nil)).Elem()
}

func (t jsonRawMessage) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(_ reflect.StructTag) (openapi.Schema, error) {
			return &openapi.Object{
				AdditionalProperties: &openapi.Raw{Schema: "true"},
			}, nil
		},
	}
}
