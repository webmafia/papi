package types

import (
	"errors"
	"mime/multipart"
	"reflect"

	"github.com/webbmaffian/papi/openapi"
)

func FileHeaderType() RequestType {
	return multipartType{}
}

type multipartType struct{}

func (t multipartType) Type() reflect.Type {
	return reflect.TypeOf((*multipart.FileHeader)(nil)).Elem()
}

func (t multipartType) CreateRequestScanner(reflect.StructTag, []string) (RequestScanner, error) {
	return nil, errors.New("scanner for 'multipart.FileHeader' is not implemented yet")
}

func (t multipartType) DescribeOperation(op *openapi.Operation) (err error) {
	return errors.New("scanner for 'multipart.FileHeader' is not implemented yet")
}
