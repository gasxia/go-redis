package redis

import (
	"fmt"
	"testing"
)

func TestNewObject(t *testing.T) {
	obj := createEmbeddedStringObject("aaa")
	fmt.Printf("%+v\n", obj)
	b := createEmbeddedStringObject("aaa")
	if compareStringObjects(obj, b) != 0 {
		t.Error("is equal")
	}
}
