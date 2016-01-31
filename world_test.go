package golax

import "net/http/httptest"

type World struct {
	Api    *Api
	Server *httptest.Server
}

func NewWorld() *World {
	api := NewApi()

	return &World{
		Api:    api,
		Server: httptest.NewServer(api),
	}
}

func (this *World) Destroy() {
	this.Server.Close()
}

func (this *World) Request(method, path string) *RequestTest {
	return NewRequestTest(method, this.Server.URL+path)
}
