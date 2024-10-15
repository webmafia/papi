package types

import (
	"errors"
	"mime/multipart"
	"reflect"

	"github.com/webbmaffian/papi/openapi"
	"github.com/webbmaffian/papi/registry"
)

func FileHeaderType() registry.RequestType {
	return multipartType{}
}

type multipartType struct{}

func (t multipartType) Type() reflect.Type {
	return reflect.TypeOf((*multipart.FileHeader)(nil)).Elem()
}

func (t multipartType) CreateRequestDecoder(reflect.StructTag, []string) (registry.RequestDecoder, error) {
	return nil, errors.New("request decoder for 'multipart.FileHeader' is not implemented yet")
}

func (t multipartType) CreateResponseEncoder(*registry.Registry, reflect.StructTag, []string, registry.ResponseEncoder) (registry.ResponseEncoder, error) {
	return nil, errors.New("response encoder for 'multipart.FileHeader' is not implemented yet")
}

func (t multipartType) DescribeOperation(op *openapi.Operation) (err error) {
	return errors.New("scanner for 'multipart.FileHeader' is not implemented yet")
}
