package coder

import (
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type Coder interface {
	Encoder(reflect.StructTag) func(*jsoniter.Stream)
}

type paramCoder interface {
	Coder
	ScanParam(ptr unsafe.Pointer, str string) error
}

type ParamCoder[T any] interface {
	Coder
	ScanParam(ptr *T, str string) error
}

type requrestCoder interface {
	Coder
	ScanRequest(ctx *fasthttp.RequestCtx, ptr unsafe.Pointer, pathParams []string) error
}

type RequestCoder[T any] interface {
	Coder
	ScanRequest(ctx *fasthttp.RequestCtx, ptr *T, pathParams []string) error
}
