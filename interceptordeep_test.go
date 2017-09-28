package golax

import (
	"testing"
)

func Test_InterceptorDeep_OK(t *testing.T) {
	world := NewWorld()
	defer world.Destroy()

	Print := func(c *Context) {
		value, exist := c.Get("list")
		if exist {
			list := value.(string)
			c.Response.Write([]byte(list))
		}
	}

	fail_s := ""
	Append := func(s string) *Interceptor {
		return &Interceptor{
			Before: func(c *Context) {

				if s == fail_s {
					c.Error(999, "Fail cause: "+fail_s)
					Print(c)
					return
				}

				list := ""
				if value, exist := c.Get("list"); exist {
					list = value.(string)
				}
				c.Set("list", list+s)
			},
		}
	}

	root := world.Api.Root

	root.Interceptor(Append("[root"))
	root.InterceptorDeep(Append("]"))

	a := root.Node("a")
	a.Interceptor(Append(",a_i1"))
	a.Interceptor(Append(",a_i2"))
	a.InterceptorDeep(Append(",a_d1"))
	a.InterceptorDeep(Append(",a_d2"))
	{
		b := a.Node("b")
		b.Interceptor(Append(",b_i1"))
		b.Interceptor(Append(",b_i2"))
		b.InterceptorDeep(Append(",b_d1"))
		b.InterceptorDeep(Append(",b_d2"))
		b.Method("GET", Print)
	}

	grandma := root.Node("grandma")
	grandma.Interceptor(Append(",grandma"))
	{
		mary := grandma.Node("mary")
		mary.Interceptor(Append(",Mary"))
		mary.Method("GET", Print)
	}

	// Test1
	res1 := world.Request("GET", "/grandma/mary").Do()
	body1 := res1.BodyString()
	//fmt.Println(body1)
	if "[root,grandma,Mary]" != body1 {
		t.Error("DeepInterceptor should be executed at the end")
	}

	// Test2
	res2 := world.Request("GET", "/a/b").Do()
	body2 := res2.BodyString()
	//fmt.Println(body2)
	if "[root,a_i1,a_i2,b_i1,b_i2,b_d2,b_d1,a_d2,a_d1]" != body2 {
		t.Error("DeepInterceptor are executed in reverse order")
	}

	// Test3
	cases := map[string]string{
		"[root": "",
		",a_i1": "[root",
		",a_i2": "[root,a_i1",
		",b_i1": "[root,a_i1,a_i2",
		",b_i2": "[root,a_i1,a_i2,b_i1",
		",b_d2": "[root,a_i1,a_i2,b_i1,b_i2",
		",b_d1": "[root,a_i1,a_i2,b_i1,b_i2,b_d2",
		",a_d2": "[root,a_i1,a_i2,b_i1,b_i2,b_d2,b_d1",
		",a_d1": "[root,a_i1,a_i2,b_i1,b_i2,b_d2,b_d1,a_d2",
		"]":     "[root,a_i1,a_i2,b_i1,b_i2,b_d2,b_d1,a_d2,a_d1",
		"":      "[root,a_i1,a_i2,b_i1,b_i2,b_d2,b_d1,a_d2,a_d1]",
	}

	for s, expected := range cases {
		fail_s = s
		res := world.Request("GET", "/a/b").Do()
		body := res.BodyString()
		if expected != body {
			t.Error("s:", s, "Expected:", expected, "Obtained:", body)
		}

	}
}
