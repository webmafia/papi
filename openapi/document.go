package openapi

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

const Version = "3.0.0"

type Document struct {
	info            Info
	servers         []Server
	paths           Paths
	securitySchemes []SecurityScheme
}

// Create a new OpenAPI root document that is ready to be used in the API service.
func NewDocument(info Info, servers ...Server) *Document {
	return &Document{
		info:    info,
		servers: servers,
		paths:   make(Paths),
	}
}

func (doc *Document) NumOperations() int {
	return len(doc.paths)
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
	doc.info.JsonEncode(ctx, s)

	if len(doc.servers) > 0 {
		s.WriteMore()
		s.WriteObjectField("servers")
		s.WriteArrayStart()

		for i := range doc.servers {
			if i != 0 {
				s.WriteMore()
			}

			doc.servers[i].JsonEncode(ctx, s)
		}

		s.WriteArrayEnd()
	}

	s.WriteMore()
	s.WriteObjectField("paths")
	doc.paths.JsonEncode(ctx, s)

	doc.encodeReferences(s, ctx)

	s.WriteObjectEnd()

	return s.Error
}

func (doc *Document) AddSecurityScheme(scheme SecurityScheme) (err error) {
	for i := range doc.securitySchemes {
		if doc.securitySchemes[i].SchemeName == scheme.SchemeName {
			return fmt.Errorf("a security scheme with name '%s' does already exist", scheme.SchemeName)
		}
	}

	doc.securitySchemes = append(doc.securitySchemes, scheme)

	return
}

func (doc *Document) encodeReferences(s *jsoniter.Stream, ctx *encoderContext) {
	s.WriteMore()
	s.WriteObjectField("components")
	s.WriteObjectStart()

	/*
		"securitySchemes": {
			"token": {
				"description": "API token",
				"type": "http",
				"scheme": "bearer",
				"bearerFormat": "base32hex"
			}
		}
	*/

	if len(doc.securitySchemes) > 0 {
		s.WriteObjectField("securitySchemes")
		s.WriteObjectStart()

		for i := range doc.securitySchemes {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteObjectField(doc.securitySchemes[i].SchemeName)
			doc.securitySchemes[i].JsonEncode(s)
		}

		// s.WriteObjectField("token")
		// s.WriteObjectStart()

		// s.WriteObjectField("description")
		// s.WriteString("API token")

		// s.WriteMore()
		// s.WriteObjectField("type")
		// s.WriteString("http")

		// s.WriteMore()
		// s.WriteObjectField("scheme")
		// s.WriteString("bearer")

		// s.WriteMore()
		// s.WriteObjectField("bearerFormat")
		// s.WriteString("base32hex")

		// s.WriteObjectEnd()

		s.WriteObjectEnd()

		if len(ctx.refs) > 0 {
			s.WriteMore()
		}
	}

	if len(ctx.refs) > 0 {
		s.WriteObjectField("schemas")
		s.WriteObjectStart()

		for i, ref := range ctx.allRefs() {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteObjectField(ref.Name)
			encodeSchema(ctx, s, ref.Schema)
		}

		s.WriteObjectEnd()
	}

	s.WriteObjectEnd()

	if len(ctx.tags) > 0 {
		s.WriteMore()
		s.WriteObjectField("tags")
		s.WriteArrayStart()

		var written bool

		tags := slices.Collect(maps.Values(ctx.tags))
		slices.SortFunc(tags, func(a, b Tag) int {
			return strings.Compare(a.Name, b.Name)
		})

		for _, tag := range tags {
			if written {
				s.WriteMore()
			} else {
				written = true
			}

			tag.JsonEncode(ctx, s)
		}

		s.WriteArrayEnd()
	}
}
