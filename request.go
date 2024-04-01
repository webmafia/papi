package fastapi

type Request[P, Q, B any] struct {
	Params      P
	QueryParams Q
	Body        B
}
