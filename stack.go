package wxpay

import (
	"container/list"
	"sync"
)

type Stack struct {
	sync.Mutex
	list *list.List
}

func NewStack() *Stack {
	return &Stack{
		list: list.New(),
	}
}

func (self *Stack) Push(v interface{}) {
	self.Lock()
	defer self.Unlock()
	self.list.PushBack(v)
}

func (self *Stack) Pop() interface{} {
	self.Lock()
	defer self.Unlock()
	top := self.list.Back()
	if top != nil {
		self.list.Remove(top)
		return top.Value
	}

	return nil
}

//type Node struct {
//	data interface{}
//	next *Node
//}
//type Stack struct {
//	header *Node
//}
//
//func (self *Stack) Push(v interface{}) {
//	n := &Node{v, self.header}
//	self.header = n
//}
//
//func (self *Stack) Pop() interface{} {
//	n := self.header
//	if n == nil {
//		return nil
//	}
//	self.header = n.next
//	return n.data
//}
