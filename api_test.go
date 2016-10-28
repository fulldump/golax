package golax

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
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

	world.Api.Root.Node("resource").Method(
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

		world.Api.Root.Node("hello").Method(method, func(c *Context) {
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

		world.Api.Root.Node("hello").Method("*", func(c *Context) {
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

		world.Api.Root.Node("hello").Method(method, func(c *Context) {
			c.Response.WriteHeader(432)
		})

		world.Api.Root.Node("hello").Method("*", func(c *Context) {
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

		world.Api.Root.Node("hello").Method(method, func(c *Context) {
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

		world.Api.Root.Node("hello").Method(method, func(c *Context) {
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

	world.Api.Root.Interceptor(&Interceptor{
		After: func(c *Context) {
			if nil != c.LastError {
				c.Response.WriteHeader(c.LastError.StatusCode)
			}
		},
	})

	world.Api.Root.Node("hello").Method("GET", func(c *Context) {
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

	world.Api.Root.Node("users").Node("{id}").Method("GET", func(c *Context) {
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

func Test_Parameters(t *testing.T) {

	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.
		Node("{a}").
		Node("{b}").
		Node("(^[0-9]+[A-Z]$)").
		Node("{{*}}").
		Method("GET", func(c *Context) {
		json.NewEncoder(c.Response).Encode(c.Parameters)
	})

	response := world.Request("GET", "/1/2/33Z/444/555/666").Do()

	obtainedBody := response.BodyJson()
	expectedBody := map[string]interface{}{
		"a":             "1",
		"b":             "2",
		"^[0-9]+[A-Z]$": "33Z",
		"*":             "444/555/666",
	}

	if !reflect.DeepEqual(obtainedBody, expectedBody) {
		t.Error("Parameters are not being collected")
	}

}

func Test_Parameters_collision(t *testing.T) {

	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.Node("{a}").Node("{a}").Method("GET", func(c *Context) {
		json.NewEncoder(c.Response).Encode(c.Parameters)
	})

	response := world.Request("GET", "/1/2").Do()

	obtainedBody := response.BodyJson()
	expectedBody := map[string]interface{}{
		"a": "2",
	}

	if !reflect.DeepEqual(obtainedBody, expectedBody) {
		t.Error("Parameters are not being collected")
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
	root.Interceptor(&Interceptor{
		After: func(c *Context) {
			if nil != c.LastError {
				c.Response.WriteHeader(c.LastError.StatusCode)
				fmt.Fprint(c.Response, c.LastError.Description)
			}
		},
	})

	users := root.Node("users")

	stats := users.Node("stats")
	stats.Method("GET", func(c *Context) {
		fmt.Fprint(c.Response, "There are 2000 users")
	})

	user := users.Node("{user_id}")
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

/**
 * https://github.com/fulldump/golax/issues/5
 * If a parameter is not the last one, it is not possible getting its value
 */
func Test_ParameterBug_issue_5(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	my_interceptor := &Interceptor{
		Before: func(c *Context) {
			c.Set("my_parameter", c.Parameter)
		},
	}

	get_profile := func(c *Context) {
		my_parameter, _ := c.Get("my_parameter")
		fmt.Fprint(c.Response, "parameter: "+my_parameter.(string))
	}

	world.Api.Root.
		Node("users").
		Node("{aa}").
		Interceptor(my_interceptor).
		Node("profile").Method("GET", get_profile)

	response := world.Request("GET", "/users/-the-value-/profile").Do()

	body := response.BodyString()

	if "parameter: -the-value-" != body {
		t.Error("Body does not match")
	}
}

func Test_handling(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	wrapper := func(text string) *Interceptor {
		return &Interceptor{
			Before: func(c *Context) {
				fmt.Println(text)
				fmt.Fprintf(c.Response, "%s(", text)
			},
			After: func(c *Context) {
				fmt.Println("/" + text)
				fmt.Fprintf(c.Response, ")%s", text)
			},
		}
	}

	root := world.Api.Root
	root.Interceptor(wrapper("root"))

	a := root.Node("a")
	a.Interceptor(wrapper("a"))

	b := a.Node("b")
	b.Interceptor(wrapper("b"))

	c := b.Node("c")
	c.Interceptor(wrapper("c"))
	c.Method("GET", func(c *Context) {
		fmt.Println("Hello world, I am C")
		fmt.Fprint(c.Response, "Hello world, I am C")
	})

	response := world.Request("GET", "/a/b/c").Do()

	body := response.BodyString()

	if "root(a(b(c(Hello world, I am C)c)b)a)root" != body {
		t.Error("Body does not match")
	}
}

func Test_RegexParameter_ok(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.Node("(^a+$)").Method("GET", func(c *Context) {
		fmt.Fprint(c.Response, "a+:", c.Parameter)
	})

	world.Api.Root.Node("(^b+$)").Method("GET", func(c *Context) {
		fmt.Fprint(c.Response, "b+:", c.Parameter)
	})

	world.Api.Root.Node("abba").Method("GET", func(c *Context) {
		fmt.Fprint(c.Response, "static:abba")
	})

	r1 := world.Request("GET", "/a").Do()
	if "a+:a" != r1.BodyString() {
		t.Error("r1: Body a+ does not match")
	}

	r2 := world.Request("GET", "/bbbbb").Do()
	if "b+:bbbbb" != r2.BodyString() {
		t.Error("Body b+ does not match")
	}

	r3 := world.Request("GET", "/abba").Do()
	if "static:abba" != r3.BodyString() {
		t.Error("Body abba does not match")
	}

	r4 := world.Request("GET", "/").Do()
	if 405 != r4.StatusCode {
		t.Error("r4: Status code does not match")
	}
}

func Test_FullPath(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Prefix = "/service"

	world.Api.Root.Node("files").Node("{{*}}").Method(
		"GET",
		func(c *Context) {
			fmt.Fprint(c.Response, c.Parameter)
		},
	)

	response := world.Request("GET", "/service/files/static/docs/document.txt").Do()

	body := response.BodyString()

	if "static/docs/document.txt" != body {
		t.Error("Parameter does not match")
	}

}

func always_break(name string) *Interceptor {
	return &Interceptor{
		Before: func(c *Context) {
			c.Response.Header().Add("Interceptors", "[BREAK "+name+"]")
			c.Error(666, "Break "+name)
		},
	}
}

func always_work(name string) *Interceptor {
	return &Interceptor{
		Before: func(c *Context) {
			c.Response.Header().Add("Interceptors", "[WORK "+name+"]")
		},
	}
}

func Test_Interceptors_ErrorChain0(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.
		Interceptor(always_work("1")).
		Interceptor(always_work("2")).
		Interceptor(always_break("Z")).
		Method("GET", func(c *Context) {
		c.Response.Header().Add("Interceptors", "[ROOT]")
	}).
		Node("node").
		Interceptor(always_work("3")).
		Interceptor(always_work("4")).
		Method("GET", func(c *Context) {
		c.Response.Header().Add("Interceptors", "[NODE]")
	})

	r := world.Request("GET", "/node").Do()

	status := r.StatusCode
	chain := strings.Join(r.Header["Interceptors"], "")
	fmt.Println(r.StatusCode, chain)

	if 666 != status {
		t.Error("Status code should be 666")
	}

	if "[WORK 1][WORK 2][BREAK Z]" != chain {
		t.Error("Chain does not match")
	}
}

func Test_Interceptors_ErrorChain1(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.
		Interceptor(always_work("1")).
		Interceptor(always_break("Z")).
		Interceptor(always_work("2")).
		Method("GET", func(c *Context) {
		c.Response.Header().Add("Interceptors", "[ROOT]")
	}).
		Node("node").
		Interceptor(always_work("3")).
		Interceptor(always_work("4")).
		Method("GET", func(c *Context) {
		c.Response.Header().Add("Interceptors", "[NODE]")
	})

	r := world.Request("GET", "/node").Do()

	status := r.StatusCode
	chain := strings.Join(r.Header["Interceptors"], "")
	fmt.Println(r.StatusCode, chain)

	if 666 != status {
		t.Error("Status code should be 666")
	}

	if "[WORK 1][BREAK Z]" != chain {
		t.Error("Chain does not match")
	}
}

func Test_Interceptors_ErrorChain2(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.
		Interceptor(always_work("1")).
		Interceptor(always_work("2")).
		Method("GET", func(c *Context) {
		c.Response.Header().Add("Interceptors", "[ROOT]")
	}).
		Node("node").
		Interceptor(always_work("3")).
		Interceptor(always_work("4")).
		Method("GET", func(c *Context) {
		c.Response.Header().Add("Interceptors", "[NODE]")
	})

	r := world.Request("GET", "/node").Do()

	status := r.StatusCode
	chain := strings.Join(r.Header["Interceptors"], "")
	fmt.Println(r.StatusCode, chain)

	if 200 != status {
		t.Error("Status code should be 666")
	}

	if "[WORK 1][WORK 2][WORK 3][WORK 4][NODE]" != chain {
		t.Error("Chain does not match")
	}
}

func Test_Interceptors_ErrorChain3(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.
		Interceptor(always_work("1")).
		Interceptor(always_work("2")).
		Method("GET", func(c *Context) {
		c.Response.Header().Add("Interceptors", "[ROOT]")
	}).
		Node("node").
		Interceptor(always_work("3")).
		Interceptor(always_break("Z")).
		Interceptor(always_work("4")).
		Method("GET", func(c *Context) {
		c.Response.Header().Add("Interceptors", "[NODE]")
	})

	r := world.Request("GET", "/node").Do()

	status := r.StatusCode
	chain := strings.Join(r.Header["Interceptors"], "")
	fmt.Println(r.StatusCode, chain)

	if 666 != status {
		t.Error("Status code should be 666")
	}

	if "[WORK 1][WORK 2][WORK 3][BREAK Z]" != chain {
		t.Error("Chain does not match")
	}

}
