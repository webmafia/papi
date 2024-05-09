package scan

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
	"github.com/webmafia/fastapi/internal/jsonpool"
)

var ioReader = reflect.TypeOf((*io.Reader)(nil)).Elem()

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
			iterateStructTags(fld.Tag, func(key, val string) (stop bool) {
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

				case "body":
					sc, err = createBodyScanner(fld.Type)

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
	if typ == ioReader {
		scan = func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, _ []string) error {
			*(*io.Reader)(p) = reqCtx.RequestBodyStream()
			return nil
		}
	} else {
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

type fieldScanner struct {
	offset uintptr
	scan   StructScanner
}

type StructScanner func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, paramVals []string) error

func iterateStructTags(tag reflect.StructTag, cb func(key, val string) bool) bool {
	if tag == "" {
		return false
	}

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if l := len(qvalue); l >= 2 {
			if cb(name, qvalue[1:l-1]) {
				break
			}
		}
	}

	return true
}
