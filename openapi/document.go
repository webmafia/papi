package openapi

type Document struct {
	Info       Info
	Servers    []Server
	Paths      Paths
	Components Components
}

func NewDocument() *Document {
	return &Document{
		Paths: make(Paths),
	}
}
