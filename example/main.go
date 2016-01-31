package main

import (
	"encoding/json"
	"strconv"

	"github.com/fulldump/golax"
)

func main() {

	my_api := golax.NewApi()
	my_api.Prefix = "/service/v1"

	my_api.Root.Middleware(golax.MiddlewareError)

	users := my_api.Root.Node("users").
		Method("GET", get_users).
		Method("POST", post_users)

	users.Node("{user_id}").
		Middleware(middleware_user).
		Method("GET", get_user).
		Method("POST", post_user).
		Method("DELETE", delete_user)

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
 * Middleware {user_id}
 * if: `user_id` exists -> load the object and put it available in the context
 * else: raise 404
 */
var middleware_user = &golax.Middleware{
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
