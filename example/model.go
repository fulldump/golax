package main

type user struct {
	id           int
	Name         string `json:"name"`
	Age          int    `json:"age"`
	Introduction string `json:"introduction"`
}

var users = map[int]*user{}
var usersLastId = 0

func insertUser(u *user) {
	usersLastId++ // NOTE: This should be thread safe in a nice server

	u.id = usersLastId
	users[u.id] = u
}

/**
 * Insert 3 sample users
 */
func init() {
	insertUser(&user{
		Name:         "Fulanito Fulanitez",
		Age:          20,
		Introduction: "Hello, I like flowers and plants",
	})

	insertUser(&user{
		Name:         "Menganito Menganez",
		Age:          30,
		Introduction: "Hi, I like wheels and cars",
	})

	insertUser(&user{
		Name:         "Zutanito Zutanez",
		Age:          40,
		Introduction: "Hey, I love cats and dogs",
	})
}
