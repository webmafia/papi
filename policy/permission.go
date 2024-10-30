package policy

import "strings"

type Permission string

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
