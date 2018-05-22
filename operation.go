package golax

import "strings"

// Operation is a terminal node, ready to execute code but exposed as Google
// custom methods (with :operation syntax)
type Operation struct {
	Path         string // Operation name
	Interceptors []*Interceptor
	Methods      map[string]Handler
}

// NewOperation instances and initialize an Operation
func NewOperation() *Operation {
	return &Operation{
		Path:         "",
		Interceptors: []*Interceptor{},
		Methods:      map[string]Handler{},
	}
}

// Method implement an HTTP method for an operation
func (o *Operation) Method(m string, h Handler) *Operation {
	M := strings.ToUpper(m)
	o.Methods[M] = h
	return o
}

// Interceptor attaches an Interceptor to an operation
func (o *Operation) Interceptor(m *Interceptor) *Operation {
	o.Interceptors = append(o.Interceptors, m)
	return o
}
