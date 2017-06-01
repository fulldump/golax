package golax

import (
	"reflect"
	"testing"
)

func Test_SplitTail_Zero(t *testing.T) {

	r := SplitTail("", ":")

	if !reflect.DeepEqual(r, []string{""}) {
		t.Error("Result should be len 1 with an empty string")
	}

}

func Test_SplitTail_One(t *testing.T) {

	r := SplitTail("abc", ":")

	if !reflect.DeepEqual(r, []string{"abc"}) {
		t.Error("Result should be len 1 with an 'abc' string")
	}

}

func Test_SplitTail_MoreThanOne(t *testing.T) {

	r := SplitTail("a:b:c:d", ":")

	if !reflect.DeepEqual(r, []string{"a:b:c", "d"}) {
		t.Error("Result should be len 2 with: 'a:b:c' and 'd' strings")
	}

}
