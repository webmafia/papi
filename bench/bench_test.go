package main

import (
	"net/http"
	"testing"

	"github.com/fasthttp/router"
	"github.com/julienschmidt/httprouter"
	"github.com/valyala/fasthttp"
)

func BenchmarkRouter(b *testing.B) {
	r := router.New()
	r.GET("/foo/{bar}/{baz}", func(ctx *fasthttp.RequestCtx) {})
	r.GET("/mjau", func(ctx *fasthttp.RequestCtx) {})
	r.GET("/mjau/abc", func(ctx *fasthttp.RequestCtx) {})
	r.GET("/mjau/abc/123", func(ctx *fasthttp.RequestCtx) {})
	var ctx fasthttp.RequestCtx

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h, _ := r.Lookup("GET", "/mjau/abc/123", &ctx)

		if h == nil {
			b.Fatal("route not found")
		}

		ctx.ResetUserValues()
	}
}

func BenchmarkHTTPRouter(b *testing.B) {
	r := httprouter.New()
	r.GET("/foo/{bar}/{baz}", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {})
	// r.GET("/mjau/abc/123", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {})
	// var ctx fasthttp.RequestCtx

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h, p, _ := r.Lookup("GET", "/foo/abc/123")

		if h == nil {
			b.Fatal("route not found")
		}

		_ = p
	}
}
