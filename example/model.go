package main

type User struct {
	id           int    `json:"-"`
	Name         string `json:"name"`
	Age          int    `json:"age"`
	Introduction string `json:"introduction"`
}

var users = map[int]*User{}
var users_last_id = 0

func insert_user(u *User) {
	users_last_id++ // NOTE: This should be thread safe in a nice server

	u.id = users_last_id
	users[u.id] = u
}

/**
 * Insert 3 sample users
 */
func init() {
	insert_user(&User{
		Name:         "Fulanito Fulanitez",
		Age:          20,
		Introduction: "Hello, I like flowers and plants",
	})

	insert_user(&User{
		Name:         "Menganito Menganez",
		Age:          30,
		Introduction: "Hi, I like wheels and cars",
	})

	insert_user(&User{
		Name:         "Zutanito Zutanez",
		Age:          40,
		Introduction: "Hey, I love cats and dogs",
	})
}
