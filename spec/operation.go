package spec

import (
	jsoniter "github.com/json-iterator/go"
	// "github.com/webmafia/fastapi/spec/schema"
)

type Operation struct {
	Id          string
	Path        string
	Method      string
	Summary     string
	Description string
	Parameters  []Parameter
	// RequestBodyRef schema.Schema
	Tags []*Tag
}

// func (op *Operation) ParamsFromStruct(v any) (err error) {
// 	typ := reflect.TypeOf(v)

// 	if typ.Kind() != reflect.Struct {
// 		return errors.New("expected struct")
// 	}

// 	numFld := typ.NumField()

// 	for i := 0; i < numFld; i++ {
// 		fld := typ.Field(i)

// 		internal.IterateStructTags(fld.Tag, func(key, val string) bool {
// 			switch key {

// 			case "param":
// 				op.Parameters = append(op.Parameters, Parameter{
// 					Name:        val,
// 					In:          InPath,
// 					Description: fld.Tag.Get("docs"),
// 					Required:    true,
// 				})

// 			case "query":
// 				op.Parameters = append(op.Parameters, Parameter{
// 					Name:        val,
// 					In:          InQuery,
// 					Description: fld.Tag.Get("docs"),
// 					// TODO: Required
// 				})

// 			}
// 			return true
// 		})
// 	}

// 	return
// }

func (op *Operation) JsonEncode(ctx *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("summary")
	s.WriteString(op.Summary)

	if op.Description != "" {
		s.WriteMore()
		s.WriteObjectField("description")
		s.WriteString(op.Description)
	}

	s.WriteMore()
	s.WriteObjectField("operationId")
	s.WriteString(op.Id)

	if len(op.Parameters) > 0 {
		s.WriteMore()
		s.WriteObjectField("parameters")

		s.WriteArrayStart()

		for i := range op.Parameters {
			if i != 0 {
				s.WriteMore()
			}

			op.Parameters[i].JsonEncode(ctx, s)
		}

		s.WriteArrayEnd()
	}

	if len(op.Tags) > 0 {
		s.WriteMore()
		s.WriteObjectField("tags")
		s.WriteArrayStart()

		for i := range op.Tags {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteString(op.Tags[i].Name)
			ctx.addTag(op.Tags[i])
		}

		s.WriteArrayEnd()
	}

	// if op.RequestBodyRef != nil {
	// 	s.WriteMore()
	// 	s.WriteObjectField("requestBody")
	// 	s.WriteObjectStart()

	// 	s.WriteObjectField("content")
	// 	s.WriteObjectStart()

	// 	s.WriteObjectField("application/json")
	// 	s.WriteObjectStart()

	// 	s.WriteObjectField("schema")
	// 	s.WriteObjectStart()

	// 	s.WriteObjectField("$ref")
	// 	s.WriteRaw(`"#/components/schemas/`)
	// 	s.WriteRaw(op.RequestBodyRef.Name())
	// 	s.WriteRaw(`"`)

	// 	s.WriteObjectEnd()

	// 	s.WriteObjectEnd()

	// 	s.WriteObjectEnd()

	// 	s.WriteObjectEnd()
	// }

	s.WriteObjectEnd()
}
