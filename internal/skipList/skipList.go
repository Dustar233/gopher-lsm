package skipList

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	defaultMaxLevel = 32
	defaultMaxP     = 2
)

type node struct {
	mu       sync.RWMutex
	key      any
	value    any
	nextNode []*node
}

type skipList struct {
	level    int
	size     int
	rng      *rand.Rand
	cmp      func(a, b any) int
	headNode *node
	maxLevel int
	P        int // 1/P
}

func New(cmp func(a, b any) int) *skipList {

	return &skipList{
		headNode: &node{
			nextNode: make([]*node, defaultMaxLevel),
		},
		level:    1,
		cmp:      cmp,
		size:     0,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		maxLevel: defaultMaxLevel,
		P:        defaultMaxP,
	}

}

func NewInt() *skipList {
	return New(
		func(a, b any) int {
			return a.(int) - b.(int)
		})
}

func (sl *skipList) randomLevel() int {

	level := 1
	for level < sl.maxLevel && sl.rng.Intn(sl.P) == 0 {
		level++
	}

	return level

}

func (sl *skipList) findPredecessors(key any) []*node {

	updates := make([]*node, sl.level)

	cur := sl.headNode

	for i := sl.level - 1; i >= 0; i-- {
		for cur.nextNode[i] != nil && sl.cmp(key, cur.nextNode[i].key) > 0 {
			cur = cur.nextNode[i]
		}
		updates[i] = cur
	}

	return updates

}

func (sl *skipList) Delete(key any) (any, bool) {
	updateTmpNodes := sl.findPredecessors(key)
	updateNodeLvl_0 := updateTmpNodes[0].nextNode[0]

	if updateNodeLvl_0 == nil || sl.cmp(key, updateNodeLvl_0.key) != 0 {
		return nil, false
	}

	oldValue := updateNodeLvl_0.value

	for i := 0; i < sl.level; i++ {
		if updateTmpNodes[i].nextNode[i] == updateNodeLvl_0 {
			updateTmpNodes[i].nextNode[i] = updateNodeLvl_0.nextNode[i]
		}
	}

	for sl.level > 1 && sl.headNode.nextNode[sl.level-1] == nil {
		sl.level--
	}

	sl.size--

	return oldValue, true

}

func (sl *skipList) Get(key any) (any, bool) {

	cur := sl.headNode

	for i := sl.level - 1; i >= 0; i-- {
		for cur.nextNode[i] != nil && sl.cmp(key, cur.nextNode[i].key) > 0 {
			cur = cur.nextNode[i]
		}

		if cur.nextNode[i] != nil && sl.cmp(key, cur.nextNode[i].key) == 0 {
			log.Printf("[Server] OP:Get Key:%v Value:%v Status:%v", key, cur.nextNode[i].value, true)
			return cur.nextNode[i].value, true
		}
	}

	log.Printf("[Server] OP:Get Key:%v Value:%v Status:%v", key, nil, false)
	return nil, false
}

// return (oldValue, exsited)
func (sl *skipList) Set(key any, value any) (any, bool) {

	updateTmpNodes := sl.findPredecessors(key)
	updateNodeLvl_0 := updateTmpNodes[0].nextNode[0]

	if updateNodeLvl_0 != nil && sl.cmp(updateNodeLvl_0.key, key) == 0 {
		oldValue := updateNodeLvl_0.value
		updateNodeLvl_0.value = value
		return oldValue, true
	}

	lvl := sl.randomLevel()

	if lvl > sl.level {
		for i := sl.level; i < lvl; i++ {
			updateTmpNodes = append(updateTmpNodes, sl.headNode)
		}
		sl.level = lvl
	}

	newNode := &node{
		key:      key,
		value:    value,
		nextNode: make([]*node, lvl),
	}

	for i := lvl - 1; i >= 0; i-- {
		newNode.nextNode[i] = updateTmpNodes[i].nextNode[i]
		updateTmpNodes[i].nextNode[i] = newNode
	}

	sl.size++
	return nil, false
}
