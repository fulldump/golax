package golax

import "strings"

// Operation is a terminal node
type Operation struct {
	Path         string // Operation name
	Interceptors []*Interceptor
	Methods      map[string]Handler
}

func NewOperation() *Operation {
	return &Operation{
		Path:         "",
		Interceptors: []*Interceptor{},
		Methods:      map[string]Handler{},
	}
}

func (o *Operation) Method(m string, h Handler) *Operation {
	M := strings.ToUpper(m)
	o.Methods[M] = h
	return o
}

func (o *Operation) Interceptor(m *Interceptor) *Operation {
	o.Interceptors = append(o.Interceptors, m)
	return o
}
