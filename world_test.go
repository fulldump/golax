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

func (w *World) Destroy() {
	w.Server.Close()
}

func (w *World) Request(method, path string) *RequestTest {
	return NewRequestTest(method, w.Server.URL+path)
}
