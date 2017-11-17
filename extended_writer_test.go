package golax

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_ExtendedWriter_WriteTwice(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	world.Api.Handler500 = func(c *Context) {
		json.NewEncoder(c.Response).Encode(map[string]interface{}{
			"status":      c.LastError.StatusCode,
			"description": c.LastError.Description,
		})
	}

	world.Api.Root.
		Interceptor(InterceptorError).
		Node("example").
		Method("GET", func(c *Context) {
			c.Response.WriteHeader(400)
			c.Response.WriteHeader(401)

			c.Response.Write([]byte("Hello\n"))
			panic("This is a panic!")
			c.Response.Write([]byte("world"))

		})

	fmt.Println("=============================================")

	res := world.Request("GET", "/example").Do()

	if 400 != res.StatusCode {
		t.Error("Status code should be 400")
	}

	body := res.BodyString()

	body_lines := strings.Split(body, "\n")

	if "Hello" != body_lines[0] {
		t.Error("First line should be `Hello` instead of " + body_lines[0])
	}

	body_json := map[string]interface{}{}

	json.Unmarshal([]byte(body_lines[1]), &body_json)

	if !reflect.DeepEqual(float64(500), body_json["status"]) {
		t.Error("Body json status should be status:500")
	}

	if !strings.HasPrefix(body_json["description"].(string), "This is a panic!") {
		t.Error("Body json description should start by `This is a panic!` ")
	}

}
