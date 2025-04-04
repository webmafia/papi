package openapi

import jsoniter "github.com/json-iterator/go"

type SecuritySchemeFlows struct {
	Implicit          struct{} // Deprecated
	Password          struct{} // TODO
	ClientCredentials struct{} // TODO
	AuthorizationCode SecuritySchemeFlow
}

func (flows *SecuritySchemeFlows) IsZero() bool {
	return flows.AuthorizationCode.IsZero()
}

func (flows *SecuritySchemeFlows) JsonEncode(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("authorizationCode")
	flows.AuthorizationCode.JsonEncode(s)

	s.WriteObjectEnd()
}

type SecuritySchemeFlow struct {
	AuthorizationUrl string
	TokenUrl         string
	RefreshUrl       string
	scopes           struct{} // TODO
}

func (flow *SecuritySchemeFlow) IsZero() bool {
	return flow.AuthorizationUrl == ""
}

func (flow *SecuritySchemeFlow) JsonEncode(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("authorizationUrl")
	s.WriteString(flow.AuthorizationUrl)

	s.WriteMore()
	s.WriteObjectField("tokenUrl")
	s.WriteString(flow.TokenUrl)

	s.WriteMore()
	s.WriteObjectField("refreshUrl")
	s.WriteString(flow.RefreshUrl)

	s.WriteMore()
	s.WriteObjectField("scopes")
	s.WriteEmptyObject()

	s.WriteObjectEnd()
}
