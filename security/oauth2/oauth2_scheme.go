package oauth2

import (
	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/security"
)

var _ security.Scheme = (*Scheme)(nil)

type Scheme struct {
	flow openapi.SecuritySchemeFlow
}

func NewScheme(flow openapi.SecuritySchemeFlow) *Scheme {
	return &Scheme{
		flow: flow,
	}
}

// OperationSecurityDocs implements security.Scheme.
func (s *Scheme) OperationSecurityDocs(permTag string) openapi.SecurityRequirement {
	sec := openapi.SecurityRequirement{
		Name: "oauth2",
	}

	if permTag != "" && permTag != "-" {
		sec.Scopes = []string{permTag}
	}

	return sec
}

// SecurityDocs implements security.Scheme.
func (s *Scheme) SecurityDocs() openapi.SecurityScheme {
	return openapi.SecurityScheme{
		SchemeName:  "oauth2",
		Type:        "oauth2",
		Description: "OAuth 2.0",
	}
}

// UserRoles implements security.Scheme.
func (s *Scheme) UserRoles(c *fasthttp.RequestCtx) (roles []string, err error) {
	panic("unimplemented")
}
