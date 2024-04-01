package spec

type Operation struct {
	Id             string
	Path           string
	Method         string
	Summary        string
	Description    string
	Parameters     []Parameter
	RequestBodyRef string
}
