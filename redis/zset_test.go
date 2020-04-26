package redis

import (
	"fmt"
	"testing"
)

func TestZaddCommand(t *testing.T) {
	//score := 1.0
	member := createStringObject("a")
	zsl := NewSkiplist()
	//zsl.Insert(score, member)
	zs := &zset{
		dict: &dict{
			ht: map[interface{}]interface{}{
			},
		},
		zsl: zsl,
	}
	zobj := createZsetObject(zs)

	dt := make(map[*robj]*robj)
	key := createStringObject("abc")
	dt[key] = zobj
	db := &Db{
		id: 1,
		dict: dt,
	}
	c := &client{
		db: db,
		argv: []*robj{nil, key, member},
	}

	zaddCommand(c)
	rs := replyBuf
	if rs != "$1\r\n1\r\n" {
		fmt.Printf("reply: %q\n", rs)
		t.Error("not equal")
	}
}
func TestZscoreCommand(t *testing.T) {
	score := 1.0
	//sobj := createStringObjectFromFloat64(score)
	member := createStringObject("a")
	zsl := NewSkiplist()
	zsl.Insert(score, member)
	zs := &zset{
		dict: &dict{
			ht: map[interface{}]interface{}{
				member: 1.0,
			},
		},
		zsl: zsl,
	}
	zobj := createZsetObject(zs)
	dt := make(map[*robj]*robj)
	key := createStringObject("abc")
	dt[key] = zobj
	db := &Db{
		id: 1,
		dict: dt,
	}
	c := &client{
		db: db,
		argv: []*robj{nil, key, member},
	}

	zscoreCommand(c)
	rs := replyBuf
	if rs != "$1\r\n1\r\n" {
		fmt.Printf("reply: %q\n", rs)
		t.Error("not equal")
	}
}
