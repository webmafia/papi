package oauth2

import (
	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/security"
)

var _ security.Gatekeeper = (*Gatekeeper)(nil)

type Gatekeeper struct {
	flow            openapi.SecuritySchemeFlow
	optionalPermTag bool
}

// OAuth 2.0 (Authorization Code flow)
func NewGatekeeper(flow openapi.SecuritySchemeFlow, optionalPermTag ...bool) *Gatekeeper {
	g := &Gatekeeper{
		flow: flow,
	}

	if len(optionalPermTag) > 0 {
		g.optionalPermTag = optionalPermTag[0]
	}

	return g
}

// OptionalPermTag implements security.Gatekeeper.
func (g *Gatekeeper) OptionalPermTag() bool {
	return g.optionalPermTag
}

// OperationSecurityDocs implements security.Scheme.
func (s *Gatekeeper) SecurityRequirement(perm security.Permission) openapi.SecurityRequirement {
	sec := openapi.SecurityRequirement{
		Name: "oauth2",
	}

	if !perm.IsZero() {
		sec.Scopes = []string{perm.String()}
	}

	return sec
}

// SecurityDocs implements security.Scheme.
func (s *Gatekeeper) SecurityScheme() openapi.SecurityScheme {
	return openapi.SecurityScheme{
		SchemeName:  "oauth2",
		Type:        "oauth2",
		Description: "OAuth 2.0",
		Flows: openapi.SecuritySchemeFlows{
			AuthorizationCode: s.flow,
		},
	}
}

// PreRequest implements security.Gatekeeper.
func (g *Gatekeeper) PreRequest(c *fasthttp.RequestCtx) error {
	return nil
}

// UserRoles implements security.Scheme.
func (s *Gatekeeper) UserRoles(c *fasthttp.RequestCtx) (roles []string, err error) {
	panic("unimplemented")
}
