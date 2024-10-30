package policy

type PolicyKey struct {
	Role     string
	Action   string
	Resource string
}

func (k PolicyKey) accessKey() accessKey {
	return accessKey{
		action:   k.Action,
		resource: k.Resource,
	}
}

type accessKey struct {
	action   string
	resource string
}
