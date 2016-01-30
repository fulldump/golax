package main

import (
	"fmt"

	"github.com/fulldump/golax"
)

func main() {

	my_api := golax.NewApi()

	service := my_api.Root.AddNode("service")
	service.AddMethod("GET", func(c *golax.Context) {
		fmt.Fprintln(c.Response, "I am the service")
	})

	v1 := service.AddNode("v1")
	v1.AddMethod("GET", func(c *golax.Context) {
		fmt.Fprintln(c.Response, "I am the version")
	})

	users := v1.AddNode("users").
		AddMethod("GET", get_users).
		AddMethod("POST", post_users)

	users.AddNode("{user_id}").
		AddMethod("GET", get_user).
		AddMethod("POST", post_user).
		AddMethod("DELETE", delete_user)

	my_api.Serve()
}

func get_users(c *golax.Context) {
	fmt.Fprintln(c.Response, "Fulanito, Menganito, Zutanito")
}

func post_users(c *golax.Context) {
	fmt.Fprintln(c.Response, "Creating a new user...")
}

func get_user(c *golax.Context) {
	fmt.Fprintln(c.Response, "Reading the user "+c.Parameter)
}

func post_user(c *golax.Context) {
	fmt.Fprintln(c.Response, "updating the user "+c.Parameter)
}

func delete_user(c *golax.Context) {
	fmt.Fprintln(c.Response, "Deleting the user "+c.Parameter)
}
