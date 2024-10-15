package registry

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/openapi"
	"github.com/webbmaffian/papi/pool/json"
)

type inputTags struct {
	Body  string `tag:"body" eq:"json"`
	Param string `tag:"param"`
	Query string `tag:"query"`
}
type fieldScanner struct {
	offset uintptr
	scan   RequestDecoder
}

type requestScanner struct {
	reg  *Registry
	json *json.Pool
}

func (r *requestScanner) CreateRequestDecoder(typ reflect.Type, tags reflect.StructTag, paramKeys []string) (scan RequestDecoder, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	numFields := typ.NumField()
	flds := make([]fieldScanner, 0, numFields)

	for i := 0; i < numFields; i++ {
		var sc RequestDecoder
		var tags inputTags

		fld := typ.Field(i)

		if sc, err := r.reg.CreateRequestDecoder(fld.Type, fld.Tag, paramKeys); err == nil {
			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if err = ScanTags(r.reg, &tags, fld.Tag); err != nil {
			return
		}

		if tags.Body == "json" {
			if sc, err = r.createJsonScanner(fld.Type); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if tags.Param != "" {
			idx := slices.Index(paramKeys, tags.Param)

			if idx < 0 {
				err = fmt.Errorf("unknown param '%s'", tags.Param)
				return
			}

			if sc, err = r.createParamScanner(fld.Type, idx); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if tags.Query != "" {
			if sc, err = r.createQueryScanner(fld.Type, tags.Query); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) (err error) {
		for _, fld := range flds {
			if err = fld.scan(unsafe.Add(p, fld.offset), c); err != nil {
				return
			}
		}

		return
	}, nil
}

func (r *requestScanner) CreateResponseEncoder(typ reflect.Type, tags reflect.StructTag, paramKeys []string, handler ResponseEncoder) (scan ResponseEncoder, err error) {
	enc := r.json.EncoderOf(typ)

	return func(c *fasthttp.RequestCtx, in, out unsafe.Pointer) error {
		if err := handler(c, in, out); err != nil {
			return err
		}

		c.SetContentType("application/json")

		s := r.json.AcquireStream(c.Response.BodyWriter())
		defer r.json.ReleaseStream(s)

		enc.Encode(out, s)
		return s.Error
	}, nil
}

func (r *requestScanner) DescribeOperation(op *openapi.Operation, typ reflect.Type) (err error) {
	if typ.Kind() != reflect.Struct {
		return errors.New("invalid struct")
	}

	numFields := typ.NumField()

	for i := 0; i < numFields; i++ {
		var tags inputTags
		fld := typ.Field(i)

		if err = ScanTags(r.reg, &tags, fld.Tag); err != nil {
			return
		}

		if tags.Body == "json" {
			schema, err := r.reg.Schema(fld.Type, fld.Tag)

			if err != nil {
				return err
			}

			op.RequestBody = schema
		}

		if tags.Param != "" {
			schema, err := r.reg.Schema(fld.Type, fld.Tag)

			if err != nil {
				return err
			}

			param := openapi.Parameter{
				Name:     tags.Param,
				In:       openapi.InPath,
				Required: true,
				Schema:   schema,
			}

			op.Parameters = append(op.Parameters, param)
		}

		if tags.Query != "" {
			schema, err := r.reg.Schema(fld.Type, fld.Tag)

			if err != nil {
				return err
			}

			param := openapi.Parameter{
				Name:   tags.Query,
				In:     openapi.InQuery,
				Schema: schema,
			}

			op.Parameters = append(op.Parameters, param)
		}
	}

	return
}
