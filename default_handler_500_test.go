package golax

import (
	"net/http"
	"testing"
)

func Test_Default500_PathHandlers(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	PanicHandler := func(c *Context) {
		panic("Something bad happened...")
	}

	world.Api.Root.
		Interceptor(InterceptorError).
		Node("a").
		Node("b").
		Node("c").
		Method("GET", PanicHandler)

	{ // Case 1
		r := world.Request("GET", "/a/b/c").Do()

		if http.StatusInternalServerError != r.StatusCode {
			t.Error("StatusCode should be 500")
		}

		body := r.BodyString()
		if "" != body {
			t.Error("Body for default 500 handler should be empty")
		}

	}

}
