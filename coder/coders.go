package coder

import (
	"reflect"
	"unsafe"
)

type Coders struct {
	paramCoders   map[reflect.Type]paramCoder
	requestCoders map[reflect.Type]requrestCoder
}

func NewCoders() *Coders {
	cs := &Coders{
		paramCoders:   make(map[reflect.Type]paramCoder),
		requestCoders: make(map[reflect.Type]requrestCoder),
	}

	return cs
}

func RegisterParamCoder[T any](cs *Coders, c ParamCoder[T]) {
	typ := reflect.TypeOf((*T)(nil))
	cs.paramCoders[typ] = *(*paramCoder)(unsafe.Pointer(&c))
}

func RegisterRequestCoder[T any](cs *Coders, c RequestCoder[T]) {
	typ := reflect.TypeOf((*T)(nil))
	cs.requestCoders[typ] = *(*requrestCoder)(unsafe.Pointer(&c))
}

func (cs *Coders) registerStandardCoders() {
	RegisterParamCoder(cs, Boolean{})
}
