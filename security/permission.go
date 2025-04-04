package security

import "strings"

type Permission string

func Perm(action, resource string) (p Permission) {
	p.set(action, resource)
	return
}

func (p Permission) String() string {
	return string(p)
}

func (p Permission) HasAction() bool {
	action, _ := p.cut()
	return action != ""
}

func (p Permission) HasResource() bool {
	_, resource := p.cut()
	return resource != ""
}

func (p Permission) HasWildcard() bool {
	action, resource := p.cut()
	return action == "*" || resource == "*"
}

func (p1 Permission) Match(p2 Permission) Permission {
	if p1 == p2 {
		return p1
	}

	a1, r1 := p1.cut()
	a2, r2 := p2.cut()

	if a1 == "*" {
		if r1 == "*" {
			return p2
		}

		return Perm(a2, r1)
	}

	if a2 == "*" {
		if r2 == "*" {
			return p1
		}

		return Perm(a1, r2)
	}

	return ""
}

func (p Permission) Action() string {
	action, _ := p.cut()
	return action
}

func (p Permission) Resource() string {
	_, resource := p.cut()
	return resource
}

func (p *Permission) SetAction(action string) {
	_, resource := p.cut()
	p.set(action, resource)
}

func (p *Permission) SetResource(resource string) {
	action, _ := p.cut()
	p.set(action, resource)
}

func (p *Permission) set(action, resource string) {
	var b strings.Builder
	b.Grow(len(action) + 1 + len(resource))
	b.WriteString(action)
	b.WriteByte(':')
	b.WriteString(resource)
	*p = Permission(b.String())
}

func (p Permission) cut() (string, string) {
	s := string(p)

	if i := strings.IndexByte(s, ':'); i >= 0 {
		return s[:i], s[i+1:]
	}

	return s, ""
}
