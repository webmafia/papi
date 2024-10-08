package openapi

import jsoniter "github.com/json-iterator/go"

const Version = "3.0.0"

type Document struct {
	Info    Info
	Servers []Server
	Paths   Paths
}

func NewDocument() *Document {
	return &Document{
		Paths: make(Paths),
	}
}

func (doc *Document) JsonEncode(s *jsoniter.Stream) (err error) {
	if err = s.Error; err != nil {
		return
	}

	ctx := newEncoderContext()

	s.WriteObjectStart()

	s.WriteObjectField("openapi")
	s.WriteString(Version)

	s.WriteMore()
	s.WriteObjectField("info")
	doc.Info.JsonEncode(ctx, s)

	if len(doc.Servers) > 0 {
		s.WriteMore()
		s.WriteObjectField("servers")
		s.WriteArrayStart()

		for i := range doc.Servers {
			if i != 0 {
				s.WriteMore()
			}

			doc.Servers[i].JsonEncode(ctx, s)
		}

		s.WriteArrayEnd()
	}

	s.WriteMore()
	s.WriteObjectField("paths")
	doc.Paths.JsonEncode(ctx, s)

	doc.encodeReferences(s, ctx)

	s.WriteObjectEnd()

	return s.Error
}

func (doc *Document) encodeReferences(s *jsoniter.Stream, ctx *encoderContext) {
	if len(ctx.refs) > 0 {
		s.WriteMore()
		s.WriteObjectField("components")
		s.WriteObjectStart()

		s.WriteObjectField("schemas")
		s.WriteObjectStart()

		var written bool

		for {
			var changed bool

			for ref := range ctx.refs {
				if ref.written {
					continue
				}

				if written {
					s.WriteMore()
				} else {
					written = true
				}

				s.WriteObjectField(ref.Name)
				ref.Schema.encodeSchema(ctx, s)
				ref.written = true
				changed = true
			}

			if !changed {
				break
			}
		}

		s.WriteObjectEnd()

		s.WriteObjectEnd()
	}
}
