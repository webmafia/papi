package security

// A user representation. The representation SHOULD be brief and only contain values that are necessary to
// comply with policies, e.g. user ID, tenant ID, roles, etc.
type User interface {
	UserRoles() []string
}
