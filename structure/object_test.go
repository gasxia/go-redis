package structure

import (
	"fmt"
	"testing"
)

func TestNewObject(t *testing.T) {
	obj := newEmbeddedStringObject("aaa")
	fmt.Printf("%+v\n", obj)
	b := newEmbeddedStringObject("aaa")
	if !obj.compareStringObjects(b) {
		t.Error("is equal")
	}
}
