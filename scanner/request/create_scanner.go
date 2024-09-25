package request

import (
	"fmt"
	"reflect"
	"slices"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
	"github.com/webmafia/fastapi/scanner/strings"
)

type inputTags struct {
	Body  string `tag:"body" eq:"json"`
	Param string `tag:"param"`
	Query string `tag:"query"`
}
type fieldScanner struct {
	offset uintptr
	scan   RequestScanner
}

func createInputScanner(sf *strings.Factory, typ reflect.Type, params []string, structTags reflect.StructTag) (scan RequestScanner, err error) {
	var sc RequestScanner
	var tags inputTags

	if err = strings.ScanString(sf, &tags, string(structTags)); err != nil {
		return
	}

	if tags.Body == "json" {
		if sc, err = createJsonScanner(api, fld.Type); err != nil {
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

		if sc, err = createParamScanner(api, fld.Type, idx); err != nil {
			return
		}

		flds = append(flds, fieldScanner{
			offset: fld.Offset,
			scan:   sc,
		})
	}

	if tags.Query != "" {
		if sc, err = createQueryScanner(api, fld.Type, tags.Query); err != nil {
			return
		}

		flds = append(flds, fieldScanner{
			offset: fld.Offset,
			scan:   sc,
		})
	}
}

func createJsonScanner(api *API, typ reflect.Type) (scan RequestScanner, err error) {
	dec := api.opt.JsonPool.DecoderOf(typ)
	scan = func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		iter := api.opt.JsonPool.AcquireIterator(c.Request.BodyStream())
		defer api.opt.JsonPool.ReleaseIterator(iter)

		dec.Decode(p, iter)
		return iter.Error
	}

	return
}

func createParamScanner(api *API, typ reflect.Type, idx int) (scan RequestScanner, err error) {
	sc, err := api.opt.StringScan.Scanner(typ)

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		params := RequestParams(c)
		return sc(p, params.Value(idx))
	}, nil
}

func createQueryScanner(api *API, typ reflect.Type, key string) (scan RequestScanner, err error) {
	sc, err := api.opt.StringScan.Scanner(typ)

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		val := c.QueryArgs().Peek(key)

		if len(val) > 0 {
			return sc(p, fast.BytesToString(val))
		}

		return nil
	}, nil
}
