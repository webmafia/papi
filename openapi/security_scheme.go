package openapi

import (
	jsoniter "github.com/json-iterator/go"
)

type SecurityScheme struct {
	SchemeName       string
	Type             string
	Description      string
	Name             string
	In               string
	Scheme           string
	BearerFormat     string
	Flows            SecuritySchemeFlows
	OpenIdConnectUrl string
	Extensions       map[string]any
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

	if sec.Name != "" {
		s.WriteMore()
		s.WriteObjectField("name")
		s.WriteString(sec.Name)
	}

	if sec.In != "" {
		s.WriteMore()
		s.WriteObjectField("in")
		s.WriteString(sec.In)
	}

	if sec.Scheme != "" {
		s.WriteMore()
		s.WriteObjectField("scheme")
		s.WriteString(sec.Scheme)
	}

	if sec.BearerFormat != "" {
		s.WriteMore()
		s.WriteObjectField("bearerFormat")
		s.WriteString(sec.BearerFormat)
	}

	if !sec.Flows.IsZero() {
		s.WriteMore()
		s.WriteObjectField("flows")
		sec.Flows.JsonEncode(s)
	}

	if sec.OpenIdConnectUrl != "" {
		s.WriteMore()
		s.WriteObjectField("openIdConnectUrl")
		s.WriteString(sec.OpenIdConnectUrl)
	}

	if len(sec.Extensions) > 0 {
		for k, v := range sec.Extensions {
			if k[:2] != "x-" {
				continue
			}

			s.WriteMore()
			s.WriteObjectField(k)
			s.WriteVal(v)
		}
	}

	s.WriteObjectEnd()
}
