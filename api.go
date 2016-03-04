package golax

import (
	"log"
	"net/http"
	"regexp"
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
	origin := "0.0.0.0:8000"
	log.Println("Server listening at " + origin)
	http.ListenAndServe(origin, a)
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
	run_interceptors(a.Root, c)

	path := r.URL.Path

	// No prefix, not found
	if !strings.HasPrefix(path, a.Prefix) {
		run_handler(a.Handler404, c)
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

	// Calculate PathHttp and PathHandler
	c.PathHandlers = a.Prefix

	// Search the right node
	current := a.Root
	for i, part := range parts {

		fullpath := false
		found := false
		for _, child := range current.Children {
			if is_parameter(child.Path) {
				if "{{*}}" == child.Path {
					// Fullpath
					fullpath = true
					c.Parameter = strings.Join(parts[i:], "/")
					c.PathHandlers += "/" + child.Path
					found = true
					current = child
					break
				}
				c.Parameter = part
				c.PathHandlers += "/" + child.Path
				found = true
				current = child
				break
			} else if is_regex(child.Path) {
				regex := child.Path[1 : len(child.Path)-1]
				if match, _ := regexp.MatchString(regex, part); match {
					c.Parameter = part
					c.PathHandlers += "/" + child.Path
					found = true
					current = child
					break
				}
			} else if part == child.Path {
				c.Parameter = ""
				c.PathHandlers += "/" + part
				found = true
				current = child
				break
			}
		}

		if fullpath {
			run_interceptors(current, c)
			break
		} else if found {
			run_interceptors(current, c)
		} else {
			c.PathHandlers = "<Handler404>"
			run_handler(a.Handler404, c)
			return
		}
	}

	method := strings.ToUpper(r.Method)
	if f, exists := current.Methods[method]; exists {
		run_handler(f, c)
		return
	}

	if f, exists := current.Methods["*"]; exists {
		run_handler(f, c)
		return
	}

	c.PathHandlers = "<Handler405>"
	run_handler(a.Handler405, c)
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

func is_regex(path string) bool {
	return '(' == path[0] && ')' == path[len(path)-1]
}

func run_interceptors(n *Node, c *Context) {
	if nil != c.LastError {
		return
	}
	for _, i := range n.Interceptors {
		if nil != i.After {
			c.afters = append(c.afters, i.After)
		}

		if nil != i.Before {
			i.Before(c)
			if nil != c.LastError {
				break
			}
		}
	}
}

func run_handler(f Handler, c *Context) {
	if nil == c.LastError {
		f(c)
	}
	for i := len(c.afters) - 1; i >= 0; i-- {
		c.afters[i](c)
	}
}
