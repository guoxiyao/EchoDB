package mymap

import (
	"echodb/internal/model"
	"fmt"
	"time"
)

// MyMap 接口定义
type MyMap interface {
	Put(id int, person *model.Person)
	Get(id int) *model.Person
	Delete(id int)
}

// TestCorrectness 正确性测试
func TestCorrectness() {
	fmt.Println("开始正确性测试...")

	// 测试所有三种实现
	testImplementation("数组实现", NewArrayList(8))
	testImplementation("链表实现", NewLinkedList())
	testImplementation("哈希实现", NewHashTable(8))

	fmt.Println("所有实现测试完成！")
}

// testImplementation 测试单个实现
func testImplementation(name string, m MyMap) {
	fmt.Printf("测试%s...\n", name)

	// 测试数据
	p1 := model.NewPerson(1, "Alice", "Smith", 25, "Female", "alice@example.com")
	p2 := model.NewPerson(2, "Bob", "Johnson", 30, "Male", "bob@example.com")

	// 测试Put和Get
	m.Put(1, p1)
	m.Put(2, p2)
	if m.Get(1) != p1 {
		fmt.Printf("%s: Put/Get测试失败\n", name)
		return
	}
	if m.Get(2) != p2 {
		fmt.Printf("%s: Put/Get测试失败\n", name)
		return
	}

	// 测试Delete
	m.Delete(1)
	if m.Get(1) != nil {
		fmt.Printf("%s: Delete测试失败\n", name)
		return
	}
	if m.Get(2) == nil {
		fmt.Printf(" %s: Delete测试失败\n", name)
		return
	}

	// 测试覆盖写入
	p1Updated := model.NewPerson(1, "Alice", "Brown", 26, "Female", "new_alice@example.com")
	m.Put(1, p1Updated)
	if m.Get(1) != p1Updated {
		fmt.Printf("%s: 覆盖写入测试失败\n", name)
		return
	}

	fmt.Printf("%s: 所有测试通过\n", name)
}

// TestPerformance 性能测试函数
func TestPerformance(m MyMap, name string, count int) {
	people := make([]*model.Person, count)
	for i := 0; i < count; i++ {
		people[i] = model.NewPerson(i, fmt.Sprintf("First%d", i), fmt.Sprintf("Last%d", i), 20+i%30, "Male", fmt.Sprintf("email%d@example.com", i))
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
