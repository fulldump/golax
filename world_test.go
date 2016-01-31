package golax

import "net/http/httptest"

type World struct {
	Api        *Api
	Server     *Server
	TestServer *httptest.Server
}

func NewWorld() *World {
	server := NewServer()
	server.Api = NewApi()

	return &World{
		Api:        server.Api,
		Server:     server,
		TestServer: httptest.NewServer(server),
	}
}

func (this *World) Destroy() {
	this.TestServer.Close()
}

func (this *World) Request(method, path string) *RequestTest {
	return NewRequestTest(method, this.TestServer.URL+path)
}
