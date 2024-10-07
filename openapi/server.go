package openapi

import jsoniter "github.com/json-iterator/go"

type Server struct {
	Description string
	Url         string
}

func (serv *Server) JsonEncode(_ *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("url")
	s.WriteString(serv.Url)

	s.WriteMore()
	s.WriteObjectField("description")
	s.WriteString(serv.Description)

	s.WriteObjectEnd()
}
