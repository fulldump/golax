package golax

import (
	"net/http"
	"strings"
)

type Server struct {
	Api        *Api
	Handler405 Handler
	Handler404 Handler
}

func default_handler_404(c *Context) {
	http.Error(c.Response, "Error 404: Not found",
		http.StatusNotFound)
}

func default_handler_405(c *Context) {
	http.Error(c.Response, "Error 405: Method not allowed",
		http.StatusMethodNotAllowed)
}

func NewServer() *Server {
	return &Server{
		Handler404: default_handler_404,
		Handler405: default_handler_405,
	}
}

func is_parameter(path string) bool {
	return '{' == path[0] && '}' == path[len(path)-1]
}

func (this *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c := NewContext()
	c.Response = w
	c.Request = r

	path := r.URL.Path

	// Remove prefix
	if strings.HasPrefix(path, this.Api.Prefix) {
		path = strings.TrimPrefix(path, this.Api.Prefix)
	}

	// Split path
	parts := strings.Split(path, "/")[1:]

	// Remove last part if empty
	last := len(parts) - 1
	if last >= 0 && "" == parts[last] {
		parts = parts[:last]
	}

	current := this.Api.Root
	push_middlewares(current, c)
	for _, part := range parts {

		found := false
		for _, child := range current.Children {
			if is_parameter(child.Path) {
				c.Parameter = part
				found = true
				current = child
				break
			} else if part == child.Path {
				found = true
				current = child
				break
			}
		}

		if found {
			push_middlewares(current, c)
		} else {
			run_handler_in_context(this.Handler404, c)
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

	run_handler_in_context(this.Handler405, c)
}

func run_handler_in_context(f Handler, c *Context) {

	afters := []Handler{}

	for _, middleware := range c.Middlewares {
		if nil != middleware.After {
			afters = append(afters, middleware.After)
		}
		if nil != middleware.Before {
			middleware.Before(c)
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

func push_middlewares(n *Node, c *Context) {
	for _, m := range n.middlewares {
		c.Middlewares = append(c.Middlewares, m)
	}
}
