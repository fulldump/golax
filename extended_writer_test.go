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
		})

	fmt.Println("=============================================")

	res := world.Request("GET", "/example").Do()

	if 400 != res.StatusCode {
		t.Error("Status code should be 400")
	}

	body := res.BodyString()

	bodyLines := strings.Split(body, "\n")

	if "Hello" != bodyLines[0] {
		t.Error("First line should be `Hello` instead of " + bodyLines[0])
	}

	bodyJson := map[string]interface{}{}

	json.Unmarshal([]byte(bodyLines[1]), &bodyJson)

	if !reflect.DeepEqual(float64(500), bodyJson["status"]) {
		t.Error("Body json status should be status:500")
	}

	if !strings.HasPrefix(bodyJson["description"].(string), "This is a panic!") {
		t.Error("Body json description should start by `This is a panic!` ")
	}

}
