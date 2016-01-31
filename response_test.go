package golax

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ResponseTest struct {
	http.Response
}

func (this *ResponseTest) BodyBytes() []byte {
	body, err := ioutil.ReadAll(this.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func (this *ResponseTest) BodyString() string {
	return string(this.BodyBytes())
}

func (this *ResponseTest) BodyJson() interface{} {
	var body interface{}
	if err := json.Unmarshal(this.BodyBytes(), &body); err != nil {
		panic(err)
	}
	return body
}

func (this *ResponseTest) BodyJsonMap() *map[string]interface{} {
	body := map[string]interface{}{}
	if err := json.Unmarshal(this.BodyBytes(), &body); err != nil {
		panic(err)
	}
	return &body
}
