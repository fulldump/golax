package golax

import (
	"fmt"
	"testing"
)

func Test_Operation_Combo1(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.
		Node("users").
		Node("{user_id}").
		Interceptor(&Interceptor{
			Before: func(c *Context) {
				fmt.Fprintln(c.Response, "Interceptor users/{user_id}", c.Parameters)
			},
		}).
		Method("GET", func(c *Context) {
			fmt.Fprintln(c.Response, "Method users/{user_id} ", c.Parameters["user_id"])
		}).
		Operation("list").
		Interceptor(&Interceptor{
			Before: func(c *Context) {
				fmt.Fprintln(c.Response, "Interceptor users/{user_id}:list", c.Parameters)
			},
		}).
		Method("GET", func(c *Context) {
			fmt.Fprintln(c.Response, "Operation 'list' operation over ", c.Parameters["user_id"])
		})

	// Case 1
	r := world.Request("GET", "/users/2").Do()

	expected := `Interceptor users/{user_id} map[user_id:2]
Method users/{user_id}  2
`

	if r.BodyString() != expected {
		t.Error("Expected interceptor + method users")
	}

	// Case 2
	r = world.Request("GET", "/users/2:list").Do()

	expected = `Interceptor users/{user_id} map[user_id:2]
Interceptor users/{user_id}:list map[user_id:2]
Operation 'list' operation over  2
`

	if r.BodyString() != expected {
		t.Error("Expected interceptor (node) + interceptor (operation) + operation")
	}

}

func Test_Operation_Combo2(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.
		Node("users:good").
		Method("GET", func(c *Context) {
			fmt.Fprint(c.Response, "I am /users:good")
		}).
		Operation("list").
		Method("GET", func(c *Context) {
			fmt.Fprint(c.Response, "I am /users:good:list")
		})

	// Case 1
	r := world.Request("GET", "/users:good").Do()

	expected := `I am /users:good`
	body := r.BodyString()

	if body != expected {
		t.Error("Expected interceptor + method users")
	}

	// Case 2
	r = world.Request("GET", "/users:good:list").Do()

	expected = `I am /users:good:list`
	body = r.BodyString()

	if body != expected {
		t.Error("Expected interceptor (node) + interceptor (operation) + operation")
	}

}

func Test_Operation_Combo3(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Root.
		Interceptor(InterceptorError).
		Node("users").
		Interceptor(&Interceptor{
			Before: func(c *Context) {
				c.Error(999, "Unexpected invented error")
			},
		}).
		Node("{user_id}").
		Operation("list").
		Method("GET", func(c *Context) {
			fmt.Fprint(c.Response, "I am /users:good:list")
		})

	// Case 1
	r := world.Request("GET", "/users/23:list").Do()

	if 999 != r.StatusCode {
		t.Error("Expected status code 999")
	}
}
