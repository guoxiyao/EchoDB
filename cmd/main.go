package main

import (
	"echodb/interfaces"
	"fmt"
)

func main() {
	// 测试正确性
	interfaces.TestCorrectness()

	// 测试性能
	count := 10000
	fmt.Printf("开始性能测试（%d条数据）...\n\n", count)
	interfaces.TestPerformance(interfaces.NewArrayList(8), "数组实现", count)
	interfaces.TestPerformance(interfaces.NewLinkedList(), "链表实现", count)
	interfaces.TestPerformance(interfaces.NewHashTable(8), "哈希表实现", count)
}
