package main

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func BenchmarkContext(b *testing.B) {
	var ctx fasthttp.RequestCtx
	params := []string{"foo", "bar", "baz"}

	b.Run("RequestCtx_Set", func(b *testing.B) {
		for range b.N {
			ctx.SetUserValue("params", &params)
		}
	})

	b.Run("RequestCtx_Get", func(b *testing.B) {
		for range b.N {
			_ = ctx.Value("params")
		}
	})

	b.Run("RequestCtx_Reslice", func(b *testing.B) {
		ctx.SetUserValue("params", make([]string, 0, 8))
		b.ResetTimer()

		for range b.N {
			params := ctx.Value("params").([]string)
			params = append(params[:0], "foo", "bar", "baz")
			ctx.SetUserValue("params", params)
		}
	})

	b.Run("RequestCtx_Reslice2", func(b *testing.B) {
		para := make([]string, 0, 8)
		ctx.SetUserValue("params", &para)
		b.ResetTimer()

		for range b.N {
			params := ctx.Value("params").(*[]string)
			*params = append((*params)[:0], "foo", "bar", "baz")
		}
	})
}
