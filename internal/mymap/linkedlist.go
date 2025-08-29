package mymap

import (
	"echodb/internal/model"
)

// LinkedList 链表实现
type LinkedList struct {
	head *listNode
	size int
}

type listNode struct {
	key   int
	value *model.Person
	next  *listNode
}

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

func (ll *LinkedList) Put(key int, value *model.Person) {
	// 先检查是否已存在，存在则更新
	current := ll.head
	for current != nil {
		if current.key == key {
			current.value = value
			return
		}
		current = current.next
	}

	// 不存在则添加到头部
	newNode := &listNode{
		key:   key,
		value: value,
		next:  ll.head,
	}
	ll.head = newNode
	ll.size++
}

func (ll *LinkedList) Get(key int) *model.Person {
	current := ll.head
	for current != nil {
		if current.key == key {
			return current.value
		}
		current = current.next
	}
	return nil
}

func (ll *LinkedList) Delete(key int) {
	if ll.head == nil {
		return
	}

	// 如果头节点就是要删除的节点
	if ll.head.key == key {
		ll.head = ll.head.next
		ll.size--
		return
	}

	// 查找并删除中间或尾部节点
	prev := ll.head
	current := ll.head.next

	for current != nil {
		if current.key == key {
			prev.next = current.next
			ll.size--
			return
		}
		prev = current
		current = current.next
	}
}

func (ll *LinkedList) Size() int {
	return ll.size
}
