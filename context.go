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
	StatusCode  int    `json: "status_code"`
	ErrorCode   int    `json: "error_code"`
	Description string `json: "description_code"`
}

func NewContext() *Context {
	return &Context{
		LastError:   nil,
		Scope:       map[string]interface{}{},
		Middlewares: []*Middleware{},
	}
}

func (this *Context) Error(s int, d string) *ContextError {
	this.Response.WriteHeader(s)
	e := &ContextError{
		StatusCode:  s,
		Description: d,
	}
	this.LastError = e
	return e
}

func (this *Context) Set(k string, v interface{}) {
	this.Scope[k] = v
}

func (this *Context) Get(k string) (interface{}, bool) {
	a, b := this.Scope[k]
	return a, b
}
