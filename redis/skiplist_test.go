package redis

import (
	"math/rand"
	"strconv"
	"testing"
)

const maxN = 1 << 16

func TestInsertAndSearch(t *testing.T) {
	list := NewSkiplist()
	r := createStringObject("1")
	node := list.Insert(1.0, r)
	if list.tail != node {
		t.Error("node")
	}
	if list.Search(1.0) != r {
		t.Error("not search r")
	}
	if list.Search(1.1) != nil {
		t.Error("search")
	}

	for i:=0; i < maxN; i++ {
		list.Insert(float64(i), createStringObject(strconv.Itoa(i)))
	}
	for i:=0; i < maxN; i++ {
		if list.Search(float64(i)).getString() != strconv.Itoa(i) {
			t.Error(i)
		}
	}

	// Test at random positions in the list.
	list = NewSkiplist()
	rList := rand.Perm(maxN)
	for _, e := range rList {
		list.Insert(float64(e), createStringObject(strconv.Itoa(e)))
	}
	for _, e := range rList {
		if list.Search(float64(e)).getString() != strconv.Itoa(e) {
			t.Fail()
		}
	}
}

func TestDelete(t *testing.T) {
	list := NewSkiplist()

	a := createStringObject("a")
	if list.Delete(0, a) != 0 {
		t.Fail()
	}
	list.Insert(1.0, a)
	if list.Delete(1.0, a) != 1 {
		t.Fail()
	}
	if list.tail != nil {
		t.Fail()
	}
}