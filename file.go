package papi

import (
	"io"
	"reflect"
	"strings"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
)

var _ registry.TypeDescriber = (*File[fileType])(nil)

type FileType interface {
	ContentType() string
	Binary() bool
}

var _ FileType = fileType{}

type fileType struct{}

func (fileType) Binary() bool        { return false }
func (fileType) ContentType() string { return "" }

type File[T FileType] struct {
	w        io.Writer
	filename string
}

func (f *File[T]) Writer() io.Writer {
	return f.w
}

func (f *File[T]) SetFilename(name string) {
	f.filename = name
}

// TypeDescription implements registry.TypeDescriber.
func (File[T]) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	var fileType T

	return registry.TypeDescription{
		Schema: func(tags reflect.StructTag) (schema openapi.Schema, err error) {
			str := &openapi.String{}

			if fileType.Binary() {
				str.Format = "binary"
			}

			return &openapi.Custom{
				ContentType: fileType.ContentType(),
				Schema: &openapi.String{
					Format: "binary",
				},
			}, nil
		},
		Handler: func(tags reflect.StructTag, handler registry.Handler) (registry.Handler, error) {
			return func(c *fasthttp.RequestCtx, in, out unsafe.Pointer) (err error) {
				c.Response.Header.SetContentType(fileType.ContentType())
				c.Response.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
				c.Response.Header.Set("Pragma", "no-cache")
				c.Response.Header.Set("Expires", "0")

				f := (*File[T])(out)
				f.w = c.Response.BodyWriter()

				if err = handler(c, in, out); err != nil {
					return
				}

				if f.filename != "" {
					c.Response.Header.Set("Content-Disposition", strings.Join([]string{`inline; filename="`, f.filename, `"`}, ""))
				}

				return
			}, nil
		},
	}
}
