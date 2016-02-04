package golax

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ResponseTest struct {
	http.Response
}

func (rt *ResponseTest) BodyBytes() []byte {
	body, err := ioutil.ReadAll(rt.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func (rt *ResponseTest) BodyString() string {
	return string(rt.BodyBytes())
}

func (rt *ResponseTest) BodyJson() interface{} {
	var body interface{}
	if err := json.Unmarshal(rt.BodyBytes(), &body); err != nil {
		panic(err)
	}
	return body
}

func (rt *ResponseTest) BodyJsonMap() *map[string]interface{} {
	body := map[string]interface{}{}
	if err := json.Unmarshal(rt.BodyBytes(), &body); err != nil {
		panic(err)
	}
	return &body
}
