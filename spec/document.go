package spec

type Document struct {
	OpenAPI           string
	Info              Info
	JsonSchemaDialect string
	Servers           []Server
	Paths             Paths
}
