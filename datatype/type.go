package datatype

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type Type interface {
	ScanTags(tags reflect.StructTag) error
	EncodeSchema(*jsoniter.Stream)
}

type ParamType[T any] interface {
	Type
	Name() string
	ScanParam(ptr *T, str string) error
}

type BodyType[T any] interface {
	Type
	ScanRequest(ctx *fasthttp.RequestCtx, ptr *T, pathParams []string) error
}
