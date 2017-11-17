package golax

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
)

type Api struct {
	Root       *Node
	Prefix     string
	Handler405 Handler
	Handler404 Handler
	Handler500 Handler
}

func NewApi() *Api {
	return &Api{
		Root:       NewNode(),
		Prefix:     "",
		Handler404: default_handler_404,
		Handler405: default_handler_405,
		Handler500: default_handler_500,
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

	defer func(c *Context) {
		if r := recover(); r != nil {
			c.Error(http.StatusInternalServerError, fmt.Sprintln(r)+string(debug.Stack()))
			a.Handler500(c)
		}
	}(c)

	run_interceptors(a.Root.Interceptors, c)
	add_deepinterceptors(a.Root.InterceptorsDeep, c)

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
	var operation *Operation = nil  // Resulting operation
	last_position := len(parts) - 1 // cache parts length
	fullpath := false
	for i, part := range parts {

		part_last := i == last_position // cache is last part

		found := false
		for _, child := range current.Children {

			child_path := child._path // cache child.Path indirection

			if part_last && child._has_operations {
				subparts := SplitTail(part, ":")
				if 2 == len(subparts) {
					subpart_1 := subparts[1]
					operation_e := false
					if operation, operation_e = child.Operations[subpart_1]; operation_e {
						part = subparts[0]
					}
				}
			}

			if part == child_path {
				c.Parameter = ""
				c.PathHandlers += "/" + part
				found = true
				current = child
				break
			} else if child._is_parameter {
				c.Parameter = part
				c.Parameters[child._parameter_key] = c.Parameter
				c.PathHandlers += "/" + child_path
				found = true
				current = child
				break
			} else if child._is_regex {
				regex := child._parameter_key
				if match, _ := regexp.MatchString(regex, part); match {
					c.Parameter = part
					c.Parameters[child._parameter_key] = c.Parameter
					c.PathHandlers += "/" + child_path
					found = true
					current = child
					break
				}
			} else if child._is_fullpath {
				fullpath = true
				c.Parameter = strings.Join(parts[i:], "/")
				c.Parameters["*"] = c.Parameter
				c.PathHandlers += "/" + child_path
				found = true
				current = child
				break
			}
		}

		if found {
			run_interceptors(current.Interceptors, c)
			add_deepinterceptors(current.InterceptorsDeep, c)
			if nil != operation {
				c.PathHandlers += ":" + operation.Path
				run_interceptors(operation.Interceptors, c)
			}
			if fullpath {
				break
			}
		} else {
			c.PathHandlers = "<Handler404>"
			run_handler(a.Handler404, c)
			return
		}
	}

	methods := current.Methods
	if nil != operation {
		methods = operation.Methods
	}

	method := strings.ToUpper(r.Method)
	if f, exists := methods[method]; exists {
		run_handler(f, c)
		return
	}

	if f, exists := methods["*"]; exists {
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

func default_handler_500(c *Context) {
	os.Stderr.WriteString(c.LastError.Description)
	c.Error(http.StatusInternalServerError, "InternalServerError")
}

func add_deepinterceptors(l []*Interceptor, c *Context) {
	if nil != c.LastError {
		return
	}
	for _, i := range l {
		c.deep_interceptors = append([]*Interceptor{i}, c.deep_interceptors...)
	}
}

func run_interceptors(l []*Interceptor, c *Context) {
	if nil != c.LastError {
		return
	}
	for _, i := range l {
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

	run_interceptors(c.deep_interceptors, c)

	if nil == c.LastError {
		f(c)
	}
	for i := len(c.afters) - 1; i >= 0; i-- {
		c.afters[i](c)
	}
}
