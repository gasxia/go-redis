package structure

import (
	"fmt"
	"math/rand"
	"testing"
)

const maxN = 1 << 16

func TestInsertAndSearch(t *testing.T) {
	list := NewSkiplist()
	node := list.Insert(1.0, 1)
	if list.tail != node {
		t.Fail()
	}
	if list.Search(1.0) != 1 {
		t.Fail()
	}
	if list.Search(1.1) != nil {
		t.Fail()
	}

	for i:=0; i < maxN; i++ {
		list.Insert(float64(i), i)
	}
	for i:=0; i < maxN; i++ {
		if list.Search(float64(i)) != i {
			t.Fail()
		}
	}

	// Test at random positions in the list.
	list = NewSkiplist()
	rList := rand.Perm(maxN)
	fmt.Printf("%+v", rList)
	for _, e := range rList {
		list.Insert(float64(e), e)
	}
	for _, e := range rList {
		if list.Search(float64(e)) != e {
			t.Fail()
		}
	}
}

func TestDelete(t *testing.T) {
	list := NewSkiplist()

	if list.Delete(0, 0) != 0 {
		t.Fail()
	}
	list.Insert(1.0, 1)
	if list.Delete(1.0, 1) != 1 {
		t.Fail()
	}
	if list.tail != nil {
		t.Fail()
	}
}