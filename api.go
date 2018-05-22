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

// Api is a complete API that implements http.Handler interface.
type Api struct {
	Root       *Node
	Prefix     string
	Handler405 Handler
	Handler404 Handler
	Handler500 Handler
}

// NewApi instances and initializes a new *Api.
func NewApi() *Api {
	return &Api{
		Root:       NewNode(),
		Prefix:     "",
		Handler404: defaultHandler404,
		Handler405: defaultHandler405,
		Handler500: defaultHandler500,
	}
}

// Serve start a default server on address 0.0.0.0:8000
func (a *Api) Serve() {
	origin := "0.0.0.0:8000"
	log.Println("Server listening at " + origin)
	http.ListenAndServe(origin, a)
}

// ServeHTTP implements http.Handler interface.
// This code is ugly but... It works! This is a critical part for the
// performance, so it has to be written with love
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

	runInterceptors(a.Root.Interceptors, c)
	addDeepInterceptors(a.Root.InterceptorsDeep, c)

	path := r.URL.Path

	// No prefix, not found
	if !strings.HasPrefix(path, a.Prefix) {
		runHandler(a.Handler404, c)
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
	var operation *Operation       // Resulting operation
	lastPosition := len(parts) - 1 // cache parts length
	fullPath := false
	for i, part := range parts {

		partLast := i == lastPosition // cache is last part

		found := false
		for _, child := range current.Children {

			childPath := child._path // cache child.Path indirection

			if partLast && child._hasOperations {
				subparts := SplitTail(part, ":")
				if 2 == len(subparts) {
					subpart1 := subparts[1]
					operationE := false
					if operation, operationE = child.Operations[subpart1]; operationE {
						part = subparts[0]
					}
				}
			}

			if part == childPath {
				c.Parameter = ""
				c.PathHandlers += "/" + part
				found = true
				current = child
				break
			} else if child._isParameter {
				c.Parameter = part
				c.Parameters[child._parameterKey] = c.Parameter
				c.PathHandlers += "/" + childPath
				found = true
				current = child
				break
			} else if child._isRegex {
				regex := child._parameterKey
				if match, _ := regexp.MatchString(regex, part); match {
					c.Parameter = part
					c.Parameters[child._parameterKey] = c.Parameter
					c.PathHandlers += "/" + childPath
					found = true
					current = child
					break
				}
			} else if child._isFullPath {
				fullPath = true
				c.Parameter = strings.Join(parts[i:], "/")
				c.Parameters["*"] = c.Parameter
				c.PathHandlers += "/" + childPath
				found = true
				current = child
				break
			}
		}

		if found {
			runInterceptors(current.Interceptors, c)
			addDeepInterceptors(current.InterceptorsDeep, c)
			if nil != operation {
				c.PathHandlers += ":" + operation.Path
				runInterceptors(operation.Interceptors, c)
			}
			if fullPath {
				break
			}
		} else {
			c.PathHandlers = "<Handler404>"
			runHandler(a.Handler404, c)
			return
		}
	}

	methods := current.Methods
	if nil != operation {
		methods = operation.Methods
	}

	method := strings.ToUpper(r.Method)
	if f, exists := methods[method]; exists {
		runHandler(f, c)
		return
	}

	if f, exists := methods["*"]; exists {
		runHandler(f, c)
		return
	}

	c.PathHandlers = "<Handler405>"
	runHandler(a.Handler405, c)
}

func defaultHandler404(c *Context) {
	c.Error(404, "Not found")
}

func defaultHandler405(c *Context) {
	c.Error(405, "Method not allowed")
}

func defaultHandler500(c *Context) {
	os.Stderr.WriteString(c.LastError.Description)
	c.Error(http.StatusInternalServerError, "InternalServerError")
}

func addDeepInterceptors(l []*Interceptor, c *Context) {
	if nil != c.LastError {
		return
	}
	for _, i := range l {
		c.deepInterceptors = append([]*Interceptor{i}, c.deepInterceptors...)
	}
}

func runInterceptors(l []*Interceptor, c *Context) {
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

func runHandler(f Handler, c *Context) {

	runInterceptors(c.deepInterceptors, c)

	if nil == c.LastError {
		f(c)
	}
	for i := len(c.afters) - 1; i >= 0; i-- {
		c.afters[i](c)
	}
}
