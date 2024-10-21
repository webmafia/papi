package types

import (
	"mime/multipart"
	"reflect"

	"github.com/webbmaffian/papi/registry"
)

func FileHeaderType() registry.TypeRegistrar {
	return multipartType{}
}

type multipartType struct{}

func (t multipartType) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	panic("unimplemented")
}

func (t multipartType) Type() reflect.Type {
	return reflect.TypeOf((*multipart.FileHeader)(nil)).Elem()
}
