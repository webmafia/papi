package openapi

import jsoniter "github.com/json-iterator/go"

type Info struct {
	Title          string
	Description    string
	TermsOfService string
	Contact        Contact
	License        License
	Version        string
}

func (i *Info) JsonEncode(ctx *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("title")
	s.WriteString(i.Title)

	s.WriteMore()
	s.WriteObjectField("description")
	s.WriteString(i.Description)

	s.WriteMore()
	s.WriteObjectField("termsOfService")
	s.WriteString(i.TermsOfService)

	if i.Contact.Name != "" {
		s.WriteMore()
		s.WriteObjectField("contact")
		i.Contact.JsonEncode(ctx, s)
	}

	if i.License.Name != "" {
		s.WriteMore()
		s.WriteObjectField("license")
		i.License.JsonEncode(ctx, s)
	}

	s.WriteMore()
	s.WriteObjectField("version")
	s.WriteString(i.Version)

	s.WriteObjectEnd()
}
