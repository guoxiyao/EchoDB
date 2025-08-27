package interfaces

import (
	"echodb/models"
	"fmt"
	"time"
)

// MyMap 接口定义
type MyMap interface {
	Put(id int, person *models.Person)
	Get(id int) *models.Person
	Delete(id int)
}

// MapType 定义Map类型
type MapType string

const (
	SliceMapType      MapType = "slice"
	LinkedListMapType MapType = "linkedlist"
	HashMapType       MapType = "hashmap"
)

// SliceMap 切片实现
type SliceMap struct {
	data []*models.Person
}

func NewSliceMap() *SliceMap {
	return &SliceMap{
		data: make([]*models.Person, 0),
	}
}

func (sm *SliceMap) Put(id int, person *models.Person) {
	for i, p := range sm.data {
		if p.ID == id {
			sm.data[i] = person
			return
		}
	}
	sm.data = append(sm.data, person)
}

func (sm *SliceMap) Get(id int) *models.Person {
	for _, p := range sm.data {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func (sm *SliceMap) Delete(id int) {
	for i, p := range sm.data {
		if p.ID == id {
			sm.data[i] = sm.data[len(sm.data)-1]
			sm.data = sm.data[:len(sm.data)-1]
			return
		}
	}
}

// Node 链表节点
type Node struct {
	person *models.Person
	next   *Node
}

// LinkedListMap 链表实现
type LinkedListMap struct {
	head *Node
}

func NewLinkedListMap() *LinkedListMap {
	return &LinkedListMap{
		head: nil,
	}
}

func (llm *LinkedListMap) Put(id int, person *models.Person) {
	// 如果链表为空，直接添加为头节点
	if llm.head == nil {
		llm.head = &Node{person: person}
		return
	}

	// 查找是否已有该ID
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

type HashMap struct {
	data map[int]*models.Person
}

func NewHashMap() *HashMap {
	return &HashMap{
		data: make(map[int]*models.Person),
	}
}

func (hm *HashMap) Put(id int, person *models.Person) {
	hm.data[id] = person
}

func (hm *HashMap) Get(id int) *models.Person {
	return hm.data[id]
}

func (hm *HashMap) Delete(id int) {
	delete(hm.data, id)
}

// NewMap 工厂函数，根据类型创建对应的Map实现
func NewMap(mapType MapType) MyMap {
	switch mapType {
	case SliceMapType:
		return NewSliceMap()
	case LinkedListMapType:
		return NewLinkedListMap()
	case HashMapType:
		return NewHashMap()
	default:
		// 默认返回HashMap
		return NewHashMap()
	}
}

// GetAvailableTypes 获取所有可用的Map类型
func GetAvailableTypes() []MapType {
	return []MapType{
		SliceMapType,
		LinkedListMapType,
		HashMapType,
	}
}

// TestCorrectness 正确性测试
func TestCorrectness() {
	fmt.Println("开始正确性测试...")

	p1 := models.NewPerson(1, "John", "Doe", 30, "Male", "john@example.com")
	p2 := models.NewPerson(2, "Jane", "Smith", 28, "Female", "jane@example.com")

	tests := []struct {
		name  string
		myMap MyMap
	}{
		{"切片实现", NewSliceMap()},
		{"链表实现", NewLinkedListMap()},
		{"哈希表实现", NewHashMap()},
	}

	for _, test := range tests {
		fmt.Printf("\n测试%s:\n", test.name)

		// 测试Put和Get
		test.myMap.Put(p1.ID, p1)
		test.myMap.Put(p2.ID, p2)

		fmt.Printf("获取ID=1: %v\n", test.myMap.Get(1))
		fmt.Printf("获取ID=2: %v\n", test.myMap.Get(2))
		fmt.Printf("获取不存在的ID=99: %v\n", test.myMap.Get(99))

		// 测试更新
		p1Updated := models.NewPerson(1, "John", "Updated", 31, "Male", "john.updated@example.com")
		test.myMap.Put(p1.ID, p1Updated) // 用更新后的内容更新ID=1
		fmt.Printf("更新后获取ID=1: %v\n", test.myMap.Get(1))

		// 测试Delete
		test.myMap.Delete(1)
		fmt.Printf("删除后获取ID=1: %v\n", test.myMap.Get(1))
	}
}

// TestPerformance 性能测试函数
func TestPerformance(m MyMap, name string, count int) {
	people := make([]*models.Person, count)
	for i := 0; i < count; i++ {
		people[i] = models.NewPerson(i, fmt.Sprintf("First%d", i), fmt.Sprintf("Last%d", i), 20+i%30, "Male", fmt.Sprintf("email%d@example.com", i))
	}

	// 测试Put性能
	start := time.Now()
	for i := 0; i < count; i++ {
		m.Put(people[i].ID, people[i])
	}
	putTime := time.Since(start)

	// 测试Get性能
	start = time.Now()
	for i := 0; i < count; i++ {
		m.Get(i)
	}
	getTime := time.Since(start)

	// 测试Delete性能
	start = time.Now()
	for i := 0; i < count; i++ {
		m.Delete(i)
	}
	deleteTime := time.Since(start)

	fmt.Printf("%s性能测试（%d条数据）:\n", name, count)
	fmt.Printf("  Put: %v (%.2f ns/op)\n", putTime, float64(putTime.Nanoseconds())/float64(count))
	fmt.Printf("  Get: %v (%.2f ns/op)\n", getTime, float64(getTime.Nanoseconds())/float64(count))
	fmt.Printf("  Delete: %v (%.2f ns/op)\n", deleteTime, float64(deleteTime.Nanoseconds())/float64(count))
	fmt.Printf("  总计: %v\n\n", putTime+getTime+deleteTime)
}
