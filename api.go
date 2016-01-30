package golax

import (
	"log"
	"net/http"
)

type Api struct {
	Root   *Node
	Prefix string
}

func NewApi() *Api {
	return &Api{
		Root:   NewNode(),
		Prefix: "",
	}
}

func (this *Api) Serve() {
	my_server := NewServer()
	my_server.Api = this

	log.Println("Server listening at 0.0.0.0:8000")
	http.ListenAndServe("0.0.0.0:8000", my_server)
}
