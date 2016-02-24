package main

import (
	"encoding/json"
	"strconv"

	"github.com/fulldump/apidoc"
	"github.com/fulldump/golax"
)

func main() {

	my_api := golax.NewApi()
	my_api.Prefix = "/service/v1"

	my_api.Root.
		Doc(golax.Doc{
		Description: `
				_Example_ is a demonstration REST API that implements a CRUD over a collection
				of users, stored in memory.

				All API calls:
				* are returning errors with the same JSON format and
				* are logging all request to standard output.
			`,
	}).
		Interceptor(golax.InterceptorLog).
		Interceptor(golax.InterceptorError)

	users := my_api.Root.Node("users").
		Doc(golax.Doc{
		Description: `
				Resource users to list and create elements. It does not support pagination,
				sorting or filtering.
			`}).
		Method("GET", get_users, golax.Doc{
		Description: `
				Return a list with a list of user ids:

				´´´json
				[1,2,3]
				´´´
			`}).
		Method("POST", post_users, golax.Doc{
		Description: `
				Create a user:
				´´´sh
				curl http://localhost:8000/service/v1/users --data '{"name": "John"}'
				´´´
				And return the user id:
				´´´json
				{"id":4}
				´´´
		`})

	users.Node("{user_id}").
		Doc(golax.Doc{
		Description: `
				Resource user to retrieve, modify and delete. A user has this structure:

				´´´json
				{
					"name": "Menganito Menganez",
					"age": 30,
					"introduction": "Hi, I like wheels and cars"
				}
				´´´
		`}).
		Interceptor(interceptor_user).
		Method("GET", get_user, golax.Doc{
		Description: `
				Return a user in JSON format. For example:
				´´´sh
				curl http://localhost:8000/service/v1/users/4
				´´´
				Will return this:
				´´´json
				{
					"name": "John",
					"age": 0,
					"introduction": ""
				}
				´´´
			`}).
		Method("POST", post_user, golax.Doc{
		Description: `
				Modify an existing user. You do not have to send all fields, for example, to
				change only the age of the user 4:

				´´´sh
				curl http://localhost:8000/service/v1/users/4 --data '{"age": 11}'
				´´´
			`}).
		Method("DELETE", delete_user, golax.Doc{
		Description: `
				Delete an existing user:

				´´´sh
				curl -X DELETE http://localhost:8000/service/v1/users/4
				´´´
			`})

	apidoc.Build(my_api)

	my_api.Serve()
}

func get_users(c *golax.Context) {
	ids := []int{}
	for id, _ := range users {
		ids = append(ids, id)
	}

	json.NewEncoder(c.Response).Encode(ids)
}

func post_users(c *golax.Context) {
	u := &User{}

	json.NewDecoder(c.Request.Body).Decode(u)

	insert_user(u)

	c.Response.WriteHeader(201)
	json.NewEncoder(c.Response).Encode(map[string]interface{}{"id": u.id})
}

func get_user(c *golax.Context) {
	u := get_context_user(c)

	json.NewEncoder(c.Response).Encode(u)
}

func post_user(c *golax.Context) {
	u := get_context_user(c)

	json.NewDecoder(c.Request.Body).Decode(u)
}

func delete_user(c *golax.Context) {
	u := get_context_user(c)
	delete(users, u.id)
}

/**
 * Interceptor {user_id}
 * if: `user_id` exists -> load the object and put it available in the context
 * else: raise 404
 */
var interceptor_user = &golax.Interceptor{
	Documentation: golax.Doc{
		Name: "User",
		Description: `
			Extract and validate user from url. If the user does not exist, a 404 will be
			returned.
		`,
	},
	Before: func(c *golax.Context) {
		user_id, _ := strconv.Atoi(c.Parameter)
		if user, exists := users[user_id]; exists {
			c.Set("user", user)
		} else {
			c.Error(404, "user `"+c.Parameter+"` does not exist")
		}
	},
}

/**
 * Helper to get a user object from the context
 */
func get_context_user(c *golax.Context) *User {
	v, _ := c.Get("user")
	return v.(*User)
}
