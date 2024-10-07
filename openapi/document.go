package openapi

import jsoniter "github.com/json-iterator/go"

const Version = "3.0.0"

type Document struct {
	Info    Info
	Servers []Server
	Paths   Paths
	// Components Components
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

	s.WriteObjectEnd()

	return s.Error
}
