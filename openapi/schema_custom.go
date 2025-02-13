package openapi

var _ Schema = (*Custom)(nil)

type Custom struct {
	ContentType string
	Schema
}
