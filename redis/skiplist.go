package redis

import (
	"math/rand"
	"time"
)

const skiplistMaxLevel = 32 /* Should be enough for 2^32 elements */
const skiplistP = 25 << 16 / 100

type skiplistLevel struct {
	forward *skiplistNode
	span    uint
}

type skiplistNode struct {
	robj     *robj
	score    float64
	backward *skiplistNode
	level    []skiplistLevel
}

type skiplist struct {
	header, tail *skiplistNode
	length       uint64
	level        int
}

func slCreateNode(level int, score float64, robj *robj) *skiplistNode {
	node := &skiplistNode{
		level: make([]skiplistLevel, level),
		score: score,
		robj:  robj,
	}
	return node
}

func slRandomLevel() int {
	level := 1
	//for rand.Float32() < skiplistP && level < skiplistMaxLevel {
	rand.Seed(time.Now().UnixNano())
	for (rand.Int31()&0xFFFF) < skiplistP && level < skiplistMaxLevel {
		level += 1
	}
	return level
}

func (sl *skiplist) Search(score float64) *robj {
	x := sl.header
	//var forward *skiplistNode
	for i := sl.level - 1; i >= 0; i-- {
		//forward = x.level[i].forward
		for x.level[i].forward != nil && x.level[i].forward.score < score {
			x = x.level[i].forward
		}
	}
	x = x.level[0].forward
	if x != nil && x.score == score {
		return x.robj
	} else {
		return nil
	}
}

func (sl *skiplist) nodeLessValue(node *skiplistNode, score float64, robj *robj) bool {
	if node == nil {
		return false
	}
	if node.score < score {
		return true
	}
	if node.score == score && compareStringObjects(node.robj, robj) < 0 {
		return true
	}
	return false
}
func (sl *skiplist) Insert(score float64, robj *robj) *skiplistNode {
	update := make([]*skiplistNode, skiplistMaxLevel)
	rank := make([]uint, skiplistMaxLevel)
	var x *skiplistNode
	var level int

	// redisAssert(!isnan(score));

	x = sl.header

	// loop invariant x.score < score
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for sl.nodeLessValue(x.level[i].forward, score, robj) {
			x = x.level[i].forward
			// 往前跳跃了span个元素
			rank[i] += x.level[i].span
		}
		update[i] = x
	}
	level = slRandomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.header
			update[i].level[i].span = uint(sl.length)
		}
		sl.level = level
	}

	x = slCreateNode(level, score, robj)
	for i := 0; i < level; i++ {
		// update new node level.forward
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		// update new node level.span
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}
	for i := level; i < sl.level; i++ {
		update[i].level[i].span++
	}

	if update[0] == sl.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		sl.tail = x
	}
	sl.length++
	return x
}

func (sl *skiplist) Delete(score float64, robj *robj) int {
	update := make([]*skiplistNode, skiplistMaxLevel)
	var x *skiplistNode

	x = sl.header
	//var forward *skiplistNode
	for i := sl.level - 1; i >= 0; i-- {
		//forward = x.level[i].forward
		for sl.nodeLessValue(x.level[i].forward, score, robj) {
			x = x.level[i].forward
		}
		update[i] = x
	}
	/* We may have multiple elements with the same score, what we need
	 * is to find the element with both the right score and object.
	 *
	 * 检查找到的元素 x ，只有在它的分值和对象都相同时，才将它删除。
	 */
	x = x.level[0].forward
	if x != nil && x.score == score && equalStringObjects(x.robj, robj) {
		sl.deleteNode(x, update)
		return 1
	} else {
		return 0
	}
}

func (sl *skiplist) deleteNode(x *skiplistNode, update []*skiplistNode) {
	for i := 0; i < sl.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].forward = x.level[i].forward
			update[i].level[i].span += x.level[i].span - 1
		} else {
			update[i].level[i].span -= 1
		}
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		sl.tail = x.backward
	}

	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

func (sl *skiplist) GetRank(x *skiplistNode, update []*skiplistNode) {

}

func NewSkiplist() *skiplist {
	sl := &skiplist{
		level:  1,
		length: 0,
		header: slCreateNode(skiplistMaxLevel, 0, nil),
	}
	for j := 0; j < skiplistMaxLevel; j++ {
		sl.header.level[j].forward = nil
		sl.header.level[j].span = 0
	}
	sl.header.backward = nil
	sl.tail = nil
	return sl
}
