package interfaces

import (
	"echodb/internal/models"
	"fmt"
)

// MyMap 接口定义
type MyMap interface {
	Put(id int, person *models.Person)
	Get(id int) *models.Person
	Delete(id int)
}

// TestCorrectness 测试正确性
func TestCorrectness() {
	fmt.Println("开始正确性测试...")
	// 这里可以添加测试逻辑
}

// TestPerformance 测试性能
func TestPerformance(m MyMap, name string, count int) {
	fmt.Printf("%s性能测试（%d条数据）:\n", name, count)
	// 这里可以添加性能测试逻辑
}
