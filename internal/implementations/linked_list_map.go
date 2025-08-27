package implementations

import (
	"echodb/internal/interfaces"
	"echodb/internal/models"
)

// Node 链表节点
type Node struct {
	person *models.Person
	next   *Node
}

// LinkedListMap 链表实现
type LinkedListMap struct {
	head *Node
}

func NewLinkedListMap() interfaces.MyMap {
	return &LinkedListMap{head: nil}
}

func (llm *LinkedListMap) Put(id int, person *models.Person) {
	if llm.head == nil {
		llm.head = &Node{person: person}
		return
	}
	current := llm.head
	for current != nil {
		if current.person.ID == id {
			current.person = person
			return
		}
		if current.next == nil {
			break
		}
		current = current.next
	}
	current.next = &Node{person: person}
}

func (llm *LinkedListMap) Get(id int) *models.Person {
	current := llm.head
	for current != nil {
		if current.person.ID == id {
			return current.person
		}
		current = current.next
	}
	return nil
}

func (llm *LinkedListMap) Delete(id int) {
	if llm.head == nil {
		return
	}
	if llm.head.person.ID == id {
		llm.head = llm.head.next
		return
	}
	current := llm.head
	for current.next != nil {
		if current.next.person.ID == id {
			current.next = current.next.next
			return
		}
		current = current.next
	}
}

