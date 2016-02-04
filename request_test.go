package golax

import (
	"net/http"
	"strings"
)

type RequestTest struct {
	http.Request
}

func NewRequestTest(method, path string) *RequestTest {
	request, err := http.NewRequest(method, path, strings.NewReader(""))
	if nil != err {
		panic(err)
	}

	return &RequestTest{*request}
}

func (rt *RequestTest) Do() *ResponseTest {
	response, err := http.DefaultClient.Do(&rt.Request)
	if err != nil {
		panic(err)
	}

	return &ResponseTest{*response}
}
