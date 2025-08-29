package mymap

import (
	"echodb/internal/model"
)

// HashTable 哈希表核心实现
type HashTable struct {
	buckets  []*hashNode
	count    int
	capacity int
}

type hashNode struct {
	key   int
	value interface{}
	next  *hashNode
}

const (
	loadFactor = 0.75 //判断是否需要扩容
)

func NewHashTable(capacity int) *HashTable {
	if capacity <= 0 {
		capacity = 16
	}
	return &HashTable{
		buckets:  make([]*hashNode, capacity),
		capacity: capacity,
	}
}

func (h *HashTable) hash(key int) uint {
	return uint(key) % uint(h.capacity)
}

func (h *HashTable) Put(key int, value *model.Person) {
	index := h.hash(key)

	// 检查是否已存在
	current := h.buckets[index]
	for current != nil {
		if current.key == key {
			current.value = value // 更新值
			return
		}
		current = current.next
	}

	// 添加新节点
	h.buckets[index] = &hashNode{
		key:   key,
		value: value,
		next:  h.buckets[index],
	}
	h.count++

	// 检查是否需要扩容
	if float64(h.count)/float64(h.capacity) > loadFactor {
		h.resize()
	}
}

func (h *HashTable) resize() {
	newCapacity := h.capacity * 2
	newBuckets := make([]*hashNode, newCapacity)

	// 重新哈希所有元素
	for i := 0; i < h.capacity; i++ {
		current := h.buckets[i]
		for current != nil {
			newIndex := uint(current.key) % uint(newCapacity)
			newBuckets[newIndex] = &hashNode{
				key:   current.key,
				value: current.value,
				next:  newBuckets[newIndex],
			}
			current = current.next
		}
	}

	h.buckets = newBuckets
	h.capacity = newCapacity
}

func (h *HashTable) Get(key int) *model.Person {
	index := h.hash(key)
	current := h.buckets[index]

	for current != nil {
		if current.key == key {
			return current.value.(*model.Person)
		}
		current = current.next
	}
	return nil
}

func (h *HashTable) Delete(key int) {
	index := h.hash(key)

	if h.buckets[index] == nil {
		return
	}

	// 检查头节点
	if h.buckets[index].key == key {
		h.buckets[index] = h.buckets[index].next
		h.count--
		return
	}

	// 检查其他节点
	prev := h.buckets[index]
	current := h.buckets[index].next

	for current != nil {
		if current.key == key {
			prev.next = current.next
			h.count--
			return
		}
		prev = current
		current = current.next
	}
}

func (h *HashTable) Size() int {
	return h.count
}
