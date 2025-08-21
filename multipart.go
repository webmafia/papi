package papi

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"unsafe"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
	"github.com/webmafia/papi/internal"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
)

var _ registry.TypeDescriber = (*MultipartFile)(nil)
var _ io.WriterTo = (*MultipartFile)(nil)

type MultipartFile struct {
	file     *multipart.FileHeader
	filetype types.Type
}

// WriteTo implements io.WriterTo.
func (m *MultipartFile) WriteTo(w io.Writer) (n int64, err error) {
	f, err := m.file.Open()

	if err != nil {
		return
	}

	defer f.Close()

	// For known file types, let's check the file header to ensure
	// that it's nothing suspicious.
	if m.filetype != filetype.Unknown {
		var head [262]byte

		if _, err = io.ReadFull(f, head[:]); err != nil && err != io.ErrUnexpectedEOF {
			return
		}

		if !filetype.IsType(head[:], m.filetype) {
			return 0, errors.New("bad file header")
		}

		if _, err = f.Seek(0, io.SeekStart); err != nil {
			return
		}
	}

	return io.CopyN(w, f, m.file.Size)
}

// TypeDescription implements registry.TypeDescriber.
func (m MultipartFile) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(_ reflect.StructTag) (schema openapi.Schema, err error) {
			return &openapi.String{
				Format: "binary",
			}, nil
		},
		Parser: registry.NoParser,
		Binder: func(fieldName string, tags reflect.StructTag) (registry.Binder, error) {
			var allow []string
			var maxSize int64
			var err error

			if v, ok := tags.Lookup("allow"); ok {
				allow = strings.Split(v, ",")
			}

			if v, ok := tags.Lookup("size"); ok {
				maxSize, err = internal.ParseBytes(v)

				if err != nil {
					return nil, err
				}
			}

			return func(c *fasthttp.RequestCtx, ptr unsafe.Pointer) (err error) {
				form, err := c.MultipartForm()

				if err != nil {
					return
				}

				files, ok := form.File[fieldName]

				// Fail silently, as the file might not be mandatory.
				if !ok || len(files) == 0 {
					return
				}

				file := files[0]
				ext := strings.ToLower(strings.TrimLeft(filepath.Ext(file.Filename), "."))

				if !slices.Contains(allow, ext) {
					return fmt.Errorf("filetype '%s' is not allowed", ext)
				}

				if file.Size > maxSize {
					return fmt.Errorf("file too large; max %d bytes is allowed", maxSize)
				}

				if file.Size == 0 {
					file.Size = maxSize
				}

				typ := filetype.GetType(ext)

				// For known file types, let's check the provided mime type to ensure
				// that it's nothing suspicious.
				if typ != filetype.Unknown {
					if contentType := file.Header.Get("Content-Type"); contentType != "" && contentType != typ.MIME.Value {
						return fmt.Errorf("invalid mime type; expected '%s', got '%s'", typ.MIME.Value, contentType)
					}
				}

				v := (*MultipartFile)(ptr)
				v.file = fast.Noescape(file)
				v.filetype = typ

				return
			}, nil
		},
	}
}
