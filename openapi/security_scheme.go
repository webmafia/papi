package openapi

import jsoniter "github.com/json-iterator/go"

type SecurityScheme struct {
	SchemeName       string
	Type             string
	Description      string
	Name             string
	In               string
	Scheme           string
	BearerFormat     string
	Flows            struct{} // TODO
	OpenIdConnectUrl struct{} // TODO
}

func (sec *SecurityScheme) IsZero() bool {
	return sec.SchemeName == ""
}

func (sec *SecurityScheme) JsonEncode(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString(sec.Type)

	if sec.Description != "" {
		s.WriteMore()
		s.WriteObjectField("description")
		s.WriteString(sec.Description)
	}

	s.WriteMore()
	s.WriteObjectField("name")
	s.WriteString(sec.Name)

	s.WriteMore()
	s.WriteObjectField("in")
	s.WriteString(sec.In)

	s.WriteMore()
	s.WriteObjectField("scheme")
	s.WriteString(sec.Scheme)

	if sec.BearerFormat != "" {
		s.WriteMore()
		s.WriteObjectField("bearerFormat")
		s.WriteString(sec.BearerFormat)
	}

	s.WriteObjectEnd()
	return
}
