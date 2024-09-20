package fastapi

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
	"github.com/webmafia/fastapi/internal"
	"github.com/webmafia/fastapi/internal/jsonpool"
)

type inputTags struct {
	Body  string `tag:"body" eq:"json"`
	Param string `tag:"param"`
	Query string `tag:"query"`
}

var scanInputTags internal.StructTagScanner = internal.Must(internal.CreateStructTagScanner(internal.ReflectType[inputTags]()))

type fieldScanner struct {
	offset uintptr
	scan   StructScanner
}

type StructScanner func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, paramVals []string) error

func createInputScanner(typ reflect.Type, params []string) (scan StructScanner, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	numFields := typ.NumField()
	flds := make([]fieldScanner, 0, numFields)

	for i := 0; i < numFields; i++ {
		var sc StructScanner
		var tags inputTags

		fld := typ.Field(i)

		if err = scanInputTags(unsafe.Pointer(&tags), fld.Tag); err != nil {
			return
		}

		if tags.Body == "json" {
			if sc, err = createJsonScanner(fld.Type); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if tags.Param != "" {
			idx := slices.Index(params, tags.Param)

			if idx < 0 {
				err = fmt.Errorf("unknown param '%s'", tags.Param)
				return
			}

			if sc, err = createParamScanner(fld.Type, idx); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if tags.Query != "" {
			if sc, err = createQueryScanner(fld.Type, tags.Query); err != nil {
				return
			}

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

func createJsonScanner(typ reflect.Type) (scan StructScanner, err error) {
	dec := jsonpool.DecoderOf(typ)
	scan = func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, _ []string) error {
		iter := jsonpool.AcquireIterator(reqCtx.Request.BodyStream())
		defer jsonpool.ReleaseIterator(iter)

		dec.Decode(p, iter)
		return iter.Error
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
