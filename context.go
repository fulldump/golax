package golax

import "net/http"

type Context struct {
	Request     *http.Request
	Response    http.ResponseWriter
	Parameter   string
	LastError   *ContextError
	Scope       map[string]interface{}
	Middlewares []*Middleware
}

type ContextError struct {
	Status      int
	Code        int
	Description string
}

func NewContext() *Context {
	return &Context{
		LastError:   nil,
		Scope:       map[string]interface{}{},
		Middlewares: []*Middleware{},
	}
}

func (this *Context) Error(s int, c int, d string) {
	this.LastError = &ContextError{
		Status:      s,
		Code:        c,
		Description: d,
	}
}

func (this *Context) Set(k string, v interface{}) {
	this.Scope[k] = v
}

func (this *Context) Get(k string) (interface{}, bool) {
	a, b := this.Scope[k]
	return a, b
}
