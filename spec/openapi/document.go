package openapi

type Document struct {
	Info
	Servers []Server
	Paths
	Components
}
