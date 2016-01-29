package api

import (
	"fmt"
	"net/http"
	"strings"
)

type Server struct {
	Api *Api
}

func is_parameter(path string) bool {
	return '{' == path[0]
}

func (this Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// TODO: process this.Api.Prefix

	path := r.URL.Path
	parts := strings.Split(path, "/")[1:]

	current := this.Api.Root
	for _, part := range parts {

		found := false
		for _, child := range current.Children {
			if is_parameter(part) {
				// Take parameter (or not)
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
			fmt.Println("No existe !!!")
			return
		}
	}

	if f, exists := current.Methods[r.Method]; exists {
		c := NewContext()
		f(w, r, c)
		return
	}

	if f, exists := current.Methods["*"]; exists {
		c := NewContext()
		f(w, r, c)
		return
	}

	fmt.Println("Method not allowed!!!")
}
