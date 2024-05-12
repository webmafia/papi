package spec

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

type Document struct {
	OpenAPI string
	Info    Info
	Servers []Server
	Paths   Paths
	Schemas map[reflect.Type]*Schema
}

func (d *Document) JsonEncode(s *jsoniter.Stream) {
	ctx := newEncoderContext()

	s.WriteObjectStart()

	s.WriteObjectField("openapi")
	s.WriteString(d.OpenAPI)

	s.WriteMore()
	s.WriteObjectField("info")
	d.Info.JsonEncode(ctx, s)

	if len(d.Servers) > 0 {
		s.WriteMore()
		s.WriteObjectField("servers")
		s.WriteArrayStart()

		for i := range d.Servers {
			if i != 0 {
				s.WriteMore()
			}

			d.Servers[i].JsonEncode(ctx, s)
		}

		s.WriteArrayEnd()
	}

	s.WriteMore()
	s.WriteObjectField("paths")
	d.Paths.JsonEncode(ctx, s)

	// TODO: Add "security"

	if len(ctx.tags) > 0 {
		s.WriteMore()
		s.WriteObjectField("tags")
		s.WriteArrayStart()

		var written bool

		for tag := range ctx.tags {
			if written {
				s.WriteMore()
			} else {
				written = true
			}

			tag.JsonEncode(ctx, s)
		}

		s.WriteArrayEnd()
	}

	s.WriteObjectEnd()
}
