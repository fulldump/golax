package golax

import "net/http"

type Context struct {
	Request      *http.Request
	Response     *ExtendedWriter
	Parameter    string
	Parameters   map[string]string
	LastError    *ContextError
	Scope        map[string]interface{}
	PathHandlers string
	afters       []Handler
}

type ContextError struct {
	StatusCode  int    `json: "status_code"`
	ErrorCode   int    `json: "error_code"`
	Description string `json: "description_code"`
}

func NewContext() *Context {
	return &Context{
		LastError:  nil,
		Parameters: map[string]string{},
		Scope:      map[string]interface{}{},
		afters:     []Handler{},
	}
}

func (c *Context) Error(s int, d string) *ContextError {
	c.Response.WriteHeader(s)
	e := &ContextError{
		StatusCode:  s,
		Description: d,
	}
	c.LastError = e
	return e
}

func (c *Context) Set(k string, v interface{}) {
	c.Scope[k] = v
}

func (c *Context) Get(k string) (interface{}, bool) {
	a, b := c.Scope[k]
	return a, b
}
