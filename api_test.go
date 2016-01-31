package golax

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func body_bytes(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func body_string(r *http.Request) string {
	return string(body_bytes(r))
}

func Test_404_ok(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	response := world.Request("GET", "/hello").Do()

	if http.StatusNotFound != response.StatusCode {
		t.Error("Status code '404' is expected")
	}
}

func Test_405_ok(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.Method("POST", func(c *Context) {
		// Do nothing
	})

	response := world.Request("GET", "/").Do()

	if http.StatusMethodNotAllowed != response.StatusCode {
		t.Error("Status code '405' is expected")
	}
}

/**
 * Whant happens if path is empty string
 */
func Test_border_case_1(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	response := world.Request("GET", "").Do()

	if http.StatusMethodNotAllowed != response.StatusCode {
		t.Error("Status code '405' is expected")
	}
}

func Test_Prefix(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Prefix = "/my/prefix/v3"

	world.Api.Root.AddNode("resource").Method(
		"GET",
		func(c *Context) {
			fmt.Fprint(c.Response, "My resource")
		},
	)

	response := world.Request("GET", "/my/prefix/v3/resource").Do()

	if "My resource" != response.BodyString() {
		t.Error("Body 'My resource' is expected")
	}
}

/**
 * Test if standard methods (and a invented one) are handleable.
 * A valid response should return the non standard `432` status code.
 */
func Test_Methods_ok(t *testing.T) {

	methods := []string{
		"OPTIONS", "GET", "HEAD",
		"POST", "PUT", "DELETE",
		"TRACE", "CONNECT", "PATCH",
		"INVENTED",
	}

	for _, method := range methods {
		world := NewWorld()

		world.Api.Root.AddNode("hello").Method(method, func(c *Context) {
			c.Response.WriteHeader(432)
		})

		response := world.Request(method, "/hello").Do()

		if 432 != response.StatusCode {
			t.Error("Method '" + method + "': Status code '432' is expected")
		}

		world.Destroy()
	}

}

/**
 * Test if standard methods (and a invented one) are handleable if are not
 * defined but the asterisk method is defined.
 * A valid response should return the non standard `432` status code.
 */
func Test_Method_asterisk_ok(t *testing.T) {

	methods := []string{
		"OPTIONS", "GET", "HEAD",
		"POST", "PUT", "DELETE",
		"TRACE", "CONNECT", "PATCH",
		"INVENTED",
	}

	for _, method := range methods {
		world := NewWorld()

		world.Api.Root.AddNode("hello").Method("*", func(c *Context) {
			c.Response.WriteHeader(432)
		})

		response := world.Request(method, "/hello").Do()

		if 432 != response.StatusCode {
			t.Error("Method '" + method + "': Status code '432' is expected")
		}

		world.Destroy()
	}

}

/**
 * Test method precedence (all methods over asterisk)
 * Status code `432` should be returned
 */
func Test_Method_not_asterisk_ok(t *testing.T) {

	methods := []string{
		"OPTIONS", "GET", "HEAD",
		"POST", "PUT", "DELETE",
		"TRACE", "CONNECT", "PATCH",
		"INVENTED",
	}

	for _, method := range methods {
		world := NewWorld()

		world.Api.Root.AddNode("hello").Method(method, func(c *Context) {
			c.Response.WriteHeader(432)
		})

		world.Api.Root.AddNode("hello").Method("*", func(c *Context) {
			c.Response.WriteHeader(431)
		})

		response := world.Request(method, "/hello").Do()

		if 432 != response.StatusCode {
			t.Error("Method '" + method + "': Status code '432' is expected")
		}

		world.Destroy()
	}

}

/**
 * methods defined as lower case should be also handled
 */
func Test_Method_lowercase_ok(t *testing.T) {

	methods := []string{
		"options", "get", "head",
		"post", "put", "delete",
		"trace", "connect", "patch",
		"invented", "opTionS", "Put", "pOst", "dELETE",
	}

	for _, method := range methods {
		world := NewWorld()

		world.Api.Root.AddNode("hello").Method(method, func(c *Context) {
			c.Response.WriteHeader(432)
		})

		METHOD := strings.ToUpper(method)
		response := world.Request(METHOD, "/hello").Do()

		if 432 != response.StatusCode {
			t.Error("Method '" + method + "': Status code '432' is expected")
		}

		world.Destroy()
	}

}

/**
 * methods defined as upper case but the http request is lowercase
 */
func Test_Method_uppercase_ok(t *testing.T) {

	methods := []string{
		"options", "get", "head",
		"post", "put", "delete",
		"trace", "connect", "patch",
		"invented", "opTionS", "Put", "pOst", "dELETE",
	}

	for _, method := range methods {
		world := NewWorld()

		world.Api.Root.AddNode("hello").Method(method, func(c *Context) {
			c.Response.WriteHeader(432)
		})

		METHOD := strings.ToLower(method)
		response := world.Request(METHOD, "/hello").Do()

		if 432 != response.StatusCode {
			t.Error("Method '" + method + "': Status code '432' is expected")
		}

		world.Destroy()
	}
}

/**
 * Call to context.Error `555`
 */
func Test_Method_error_555(t *testing.T) {

	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.Middleware(&Middleware{
		After: func(c *Context) {
			if nil != c.LastError {
				c.Response.WriteHeader(c.LastError.StatusCode)
			}
		},
	})

	world.Api.Root.AddNode("hello").Method("GET", func(c *Context) {
		c.Error(555, "Sample error")
	})

	response := world.Request("GET", "/hello").Do()

	if 555 != response.StatusCode {
		t.Error("Status code '555' is expected")
	}

}

func Test_Parameter(t *testing.T) {

	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.AddNode("users").AddNode("{id}").Method("GET", func(c *Context) {
		fmt.Fprintln(c.Response, "The user is "+c.Parameter)
	})

	response := world.Request("GET", "/users/42").Do()

	if 200 != response.StatusCode {
		t.Error("Status code '200' is expected")
	}

	if "The user is 42\n" != response.BodyString() {
		t.Error("Body 'The user is 42\\n' is expected")
	}

}

/**
 * The users node has two nodes in order:
 * - stats
 * - {user_id}
 * GET /users/stats should return 200 `There are 2000 users`
 * GET /users/1231 should return 200 `User 1231`
 * Get /users/9999 should return 404 `User 9999 does not exist`
 */
func Test_Parameter_precedence(t *testing.T) {

	world := NewWorld()
	defer world.Destroy()

	root := world.Api.Root
	root.Middleware(&Middleware{
		After: func(c *Context) {
			if nil != c.LastError {
				c.Response.WriteHeader(c.LastError.StatusCode)
				fmt.Fprint(c.Response, c.LastError.Description)
			}
		},
	})

	users := root.AddNode("users")

	stats := users.AddNode("stats")
	stats.Method("GET", func(c *Context) {
		fmt.Fprint(c.Response, "There are 2000 users")
	})

	user := users.AddNode("{user_id}")
	user.Method("GET", func(c *Context) {
		user_id, _ := strconv.Atoi(c.Parameter)
		if user_id > 2000 {
			c.Error(404, "User "+c.Parameter+" does not exist")
			return
		}

		fmt.Fprint(c.Response, "User "+c.Parameter)
	})

	response1 := world.Request("GET", "/users/stats").Do()
	if 200 != response1.StatusCode {
		t.Error("Status code `200` is expected")
	}
	if "There are 2000 users" != response1.BodyString() {
		t.Error("Body `There are 2000 users` is expected")
	}

	response2 := world.Request("GET", "/users/1231").Do()
	if 200 != response2.StatusCode {
		t.Error("Status code `200` is expected")
	}
	if "User 1231" != response2.BodyString() {
		t.Error("Body `User 1231` is expected")
	}

	response3 := world.Request("GET", "/users/9999").Do()
	if 404 != response3.StatusCode {
		t.Error("Status code `404` is expected")
	}
	if "User 9999 does not exist" != response3.BodyString() {
		t.Error("Body `User 9999 does not exist` is expected")
	}

}
