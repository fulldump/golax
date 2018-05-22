package golax

import "net/http"

// Context is a space to store information to be passed between interceptors and
// the final handler.
type Context struct {
	Request          *http.Request
	Response         *ExtendedWriter
	Parameter        string
	Parameters       map[string]string
	LastError        *ContextError
	Scope            map[string]interface{}
	PathHandlers     string
	afters           []Handler
	deepInterceptors []*Interceptor
}

// ContextError is the error passed back when context.Error is called. It can be
// used inside an interceptor or a handler
type ContextError struct {
	StatusCode  int    `json:"status_code"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description_code"`
}

// NewContext instances and initializes a Context
func NewContext() *Context {
	return &Context{
		LastError:        nil,
		Parameters:       map[string]string{},
		Scope:            map[string]interface{}{},
		afters:           []Handler{},
		deepInterceptors: []*Interceptor{},
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

// Set stores a key-value tuple inside a Context
func (c *Context) Set(k string, v interface{}) {
	c.Scope[k] = v
}

// Get retrieves a value from a Context
func (c *Context) Get(k string) (interface{}, bool) {
	a, b := c.Scope[k]
	return a, b
}
