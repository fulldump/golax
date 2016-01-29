package main

import (
	"fmt"
	"net/http"

	"golax/api"
)

func main() {

	my_api := api.NewApi()

	my_api.Root.
		AddNode("this").
		AddNode("is").
		AddNode("my path").
		AddMethod("GET", func(w http.ResponseWriter, r *http.Request, c *api.Context) {
		fmt.Fprintln(w, "Hello world!! you are in This is my path! :) ")
	})

	my_api.Serve()

}
