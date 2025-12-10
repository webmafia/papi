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
var _ io.Reader = (*MultipartFile)(nil)
var _ io.WriterTo = (*MultipartFile)(nil)

type MultipartFile struct {
	file     *multipart.FileHeader
	filetype types.Type
	r        multipart.File
}

func (m *MultipartFile) openAndValidate() (multipart.File, error) {
	f, err := m.file.Open()
	if err != nil {
		return nil, err
	}

	if m.filetype != filetype.Unknown {
		var head [1024]byte

		if _, err = io.ReadFull(f, head[:]); err != nil && err != io.ErrUnexpectedEOF {
			f.Close()
			return nil, err
		}

		if !filetype.IsType(head[:], m.filetype) {
			f.Close()
			return nil, errors.New("bad file header")
		}

		if _, err = f.Seek(0, io.SeekStart); err != nil {
			f.Close()
			return nil, err
		}
	}
	return f, nil
}

func (m *MultipartFile) reader() (multipart.File, error) {
	if m.r != nil {
		return m.r, nil
	}
	f, err := m.openAndValidate()
	if err != nil {
		return nil, err
	}
	m.r = f
	return f, nil
}

func (m *MultipartFile) Read(p []byte) (int, error) {
	r, err := m.reader()
	if err != nil {
		return 0, err
	}
	return r.Read(p)
}

func (m *MultipartFile) WriteTo(w io.Writer) (int64, error) {
	f, err := m.openAndValidate()
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return io.CopyN(w, f, m.file.Size)
}

func (m *MultipartFile) Filename() string {
	if m.file == nil {
		return ""
	}

	return m.file.Filename
}

func (m *MultipartFile) Type() types.Type {
	return m.filetype
}

func (m *MultipartFile) IsType(ext string) bool {
	if m.file == nil {
		return false
	}

	return filetype.GetType(ext) == m.filetype
}

func (m *MultipartFile) Size() int64 {
	if m.file == nil {
		return 0
	}

	return m.file.Size
}

func (m *MultipartFile) Header(key string) string {
	if m.file == nil {
		return ""
	}

	return m.file.Header.Get(key)
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
			allow, maxSize, err := parseTags(tags)

			if err != nil {
				return nil, err
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
				typ, err := validateFile(file, allow, maxSize)

				if err != nil {
					return
				}

				v := (*MultipartFile)(ptr)
				v.file = fast.Noescape(file)
				v.filetype = typ

				return
			}, nil
		},
	}
}

func validateFile(file *multipart.FileHeader, allow []string, maxSize int64) (typ types.Type, err error) {
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(file.Filename), "."))

	if !slices.Contains(allow, ext) {
		err = fmt.Errorf("filetype '%s' is not allowed", ext)
		return
	}

	if file.Size > maxSize {
		err = fmt.Errorf("file too large; max %d bytes is allowed", maxSize)
		return
	}

	if file.Size == 0 {
		file.Size = maxSize
	}

	typ = filetype.GetType(ext)

	// For known file types, let's check the provided mime type to ensure
	// that it's nothing suspicious.
	if typ != filetype.Unknown {
		if contentType := file.Header.Get("Content-Type"); contentType != "" && contentType != typ.MIME.Value {
			err = fmt.Errorf("invalid mime type; expected '%s', got '%s'", typ.MIME.Value, contentType)
			return
		}
	}

	return
}

func parseTags(tags reflect.StructTag) (allow []string, maxSize int64, err error) {
	if v, ok := tags.Lookup("allow"); ok {
		allow = strings.Split(v, ",")
	}

	if v, ok := tags.Lookup("size"); ok {
		maxSize, err = internal.ParseBytes(v)

		if err != nil {
			return nil, 0, err
		}
	}

	return
}

type multipartFiles struct{}

func (t multipartFiles) Type() reflect.Type {
	return reflect.TypeOf((*[]MultipartFile)(nil)).Elem()
}

func (t multipartFiles) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(_ reflect.StructTag) (openapi.Schema, error) {
			return &openapi.Array{
				Items: &openapi.String{
					Format: "binary",
				},
			}, nil
		},
		Binder: func(fieldName string, tags reflect.StructTag) (registry.Binder, error) {
			allow, maxSize, err := parseTags(tags)

			if err != nil {
				return nil, err
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

				v := (*[]MultipartFile)(ptr)

				*v = make([]MultipartFile, len(files))

				for i, file := range files {
					typ, err := validateFile(file, allow, maxSize)

					if err != nil {
						return err
					}

					(*v)[i].file = fast.Noescape(file)
					(*v)[i].filetype = typ
				}

				return
			}, nil
		},
	}
}
