package api

import "net/http"

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
	my_server := Server{
		Api: this,
	}

	http.ListenAndServe("localhost:8000", my_server)
}
