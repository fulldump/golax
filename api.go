package golax

import (
	"log"
	"net/http"
	"strings"
)

type Api struct {
	Root       *Node
	Prefix     string
	Handler405 Handler
	Handler404 Handler
}

func NewApi() *Api {
	return &Api{
		Root:       NewNode(),
		Prefix:     "",
		Handler404: default_handler_404,
		Handler405: default_handler_405,
	}
}

func (a *Api) Serve() {
	log.Println("Server listening at 0.0.0.0:8000")
	http.ListenAndServe("0.0.0.0:8000", a)
}

/**
 * This code is ugly but... It works! This is a critical part for the
 * performance, so it has to be written with love
 */
func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Create the context and populate it
	c := NewContext()
	c.Response = NewExtendedWriter(w)
	c.Request = r
	push_interceptors(a.Root, c)

	path := r.URL.Path

	// No prefix, not found
	if !strings.HasPrefix(path, a.Prefix) {
		run_handler_in_context(a.Handler404, c)
		return
	}

	// Remove prefix
	path = strings.TrimPrefix(path, a.Prefix)

	// Split path
	parts := strings.Split(path, "/")[1:]

	// Remove last part if empty
	last := len(parts) - 1
	if last >= 0 && "" == parts[last] {
		parts = parts[:last]
	}

	// Search the right node
	current := a.Root
	for _, part := range parts {

		found := false
		for _, child := range current.Children {
			if is_parameter(child.Path) {
				c.Parameter = part
				found = true
				current = child
				break
			} else if part == child.Path {
				c.Parameter = ""
				found = true
				current = child
				break
			}
		}

		if found {
			push_interceptors(current, c)
		} else {
			run_handler_in_context(a.Handler404, c)
			return
		}
	}

	method := strings.ToUpper(r.Method)
	if f, exists := current.Methods[method]; exists {
		run_handler_in_context(f, c)
		return
	}

	if f, exists := current.Methods["*"]; exists {
		run_handler_in_context(f, c)
		return
	}

	run_handler_in_context(a.Handler405, c)
}

func default_handler_404(c *Context) {
	c.Error(404, "Not found")
}

func default_handler_405(c *Context) {
	c.Error(405, "Method not allowed")
}

func is_parameter(path string) bool {
	return '{' == path[0] && '}' == path[len(path)-1]
}

func push_interceptors(n *Node, c *Context) {
	for _, m := range n.interceptors {
		c.Interceptors = append(c.Interceptors, m)
	}
}

/**
 * This code is very ugly... it also has a huge impact over the performance so
 * tricks are welcome
 */
func run_handler_in_context(f Handler, c *Context) {
	afters := []Handler{}

	for _, interceptor := range c.Interceptors {
		if nil != interceptor.After {
			afters = append(afters, interceptor.After)
		}
		if nil != interceptor.Before {
			interceptor.Before(c)
			if nil != c.LastError {
				break
			}
		}
	}

	if nil == c.LastError {
		f(c)
	}

	for i := len(afters) - 1; i >= 0; i-- {
		afters[i](c)
	}
}
