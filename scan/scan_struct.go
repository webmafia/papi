package scan

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"slices"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
	"github.com/webmafia/fastapi/internal"
	"github.com/webmafia/fastapi/internal/jsonpool"
)

const defaultMaxInMemoryFileSize = 16 * 1024 * 1024

var (
	ioReader            = reflect.TypeOf((*io.Reader)(nil)).Elem()
	multipartForm       = reflect.TypeOf((*multipart.Form)(nil))
	multipartFileHeader = reflect.TypeOf((*multipart.FileHeader)(nil))
)

func CreateStructScanner(typ reflect.Type, params []string) (scan StructScanner, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	numFields := typ.NumField()
	flds := make([]fieldScanner, 0, numFields)

	for i := 0; i < numFields; i++ {
		var sc StructScanner
		fld := typ.Field(i)

		if fld.Name == "Body" {
			sc, err = createBodyScanner(fld.Type)
		} else {
			internal.IterateStructTags(fld.Tag, func(key, val string) (stop bool) {
				switch key {

				case "param":
					idx := slices.Index(params, val)

					if idx < 0 {
						err = fmt.Errorf("unknown param '%s'", val)
					} else {
						sc, err = createParamScanner(fld.Type, idx)
					}

				case "query":
					sc, err = createQueryScanner(fld.Type, val)

				default:
					return false
				}

				return true
			})
		}

		if err != nil {
			return
		}

		if sc != nil {
			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}
	}

	return func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, paramVals []string) (err error) {
		for _, fld := range flds {
			if err = fld.scan(unsafe.Add(p, fld.offset), reqCtx, paramVals); err != nil {
				return
			}
		}

		return
	}, nil
}

func createBodyScanner(typ reflect.Type) (scan StructScanner, err error) {
	switch typ {

	case ioReader:
		scan = func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, _ []string) error {
			*(*io.Reader)(p) = reqCtx.RequestBodyStream()
			return nil
		}

	case multipartForm:
		scan = func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, _ []string) (err error) {
			form, err := parseMultipartForm(reqCtx)

			if err != nil {
				return
			}

			*(*multipart.Form)(p) = *form
			return
		}

	case multipartFileHeader:
		scan = func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, _ []string) (err error) {
			form, err := parseMultipartForm(reqCtx)

			if err != nil {
				return
			}

			for _, files := range form.File {
				if len(files) == 0 {
					continue
				}

				*(*multipart.FileHeader)(p) = *files[0]
				return
			}

			return errors.New("missing file")
		}

	default:
		dec := jsonpool.DecoderOf(typ)
		scan = func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, _ []string) error {
			iter := jsonpool.AcquireIterator(reqCtx.Request.BodyStream())
			defer jsonpool.ReleaseIterator(iter)

			dec.Decode(p, iter)
			return iter.Error
		}
	}

	return
}

func createParamScanner(typ reflect.Type, idx int) (scan StructScanner, err error) {
	sc, err := CreateScanner(typ)

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, paramVals []string) error {
		return sc(p, paramVals[idx])
	}, nil
}

func createQueryScanner(typ reflect.Type, key string) (scan StructScanner, err error) {
	sc, err := CreateScanner(typ)

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, paramVals []string) error {
		val := reqCtx.QueryArgs().Peek(key)

		if len(val) > 0 {
			return sc(p, fast.BytesToString(val))
		}

		return nil
	}, nil
}

func parseMultipartForm(reqCtx *fasthttp.RequestCtx) (form *multipart.Form, err error) {
	bounds := fast.BytesToString(reqCtx.Request.Header.MultipartFormBoundary())

	if len(bounds) > 0 && len(reqCtx.Request.Header.Peek(fasthttp.HeaderContentEncoding)) == 0 {
		form, err = readMultipartForm(reqCtx.RequestBodyStream(), bounds, reqCtx.Request.Header.ContentLength(), defaultMaxInMemoryFileSize)
	} else {
		err = errors.New("expected multipart upload")
	}

	return
}

func readMultipartForm(r io.Reader, boundary string, size, maxInMemoryFileSize int) (*multipart.Form, error) {
	// Do not care about memory allocations here, since they are tiny
	// compared to multipart data (aka multi-MB files) usually sent
	// in multipart/form-data requests.

	if size <= 0 {
		return nil, fmt.Errorf("form size must be greater than 0. Given %d", size)
	}
	lr := io.LimitReader(r, int64(size))
	mr := multipart.NewReader(lr, boundary)
	f, err := mr.ReadForm(int64(maxInMemoryFileSize))
	if err != nil {
		return nil, fmt.Errorf("cannot read multipart/form-data body: %w", err)
	}
	return f, nil
}

type fieldScanner struct {
	offset uintptr
	scan   StructScanner
}

type StructScanner func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, paramVals []string) error
