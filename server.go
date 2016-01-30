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

	// TODO: process this.Api.Prefix

	path := r.URL.Path
	parts := strings.Split(path, "/")[1:]

	current := this.Api.Root
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

		if !found {
			this.Handler404(c)
			return
		}
	}

	if f, exists := current.Methods[r.Method]; exists {
		f(c)
		return
	}

	if f, exists := current.Methods["*"]; exists {
		f(c)
		return
	}

	this.Handler405(c)
}
