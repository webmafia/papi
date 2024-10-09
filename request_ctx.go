package fastapi

import (
	"unsafe"

	"github.com/valyala/fasthttp"
)

type RequestCtx fasthttp.RequestCtx

func papiToFasthttp(ctx *RequestCtx) *fasthttp.RequestCtx {
	return (*fasthttp.RequestCtx)(unsafe.Pointer(ctx))
}

func fasthttpToPapi(ctx *fasthttp.RequestCtx) *RequestCtx {
	return (*RequestCtx)(unsafe.Pointer(ctx))
}
