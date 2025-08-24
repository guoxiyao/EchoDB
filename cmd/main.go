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
	interfaces.TestPerformance(interfaces.NewSliceMap(), "切片实现", count)
	interfaces.TestPerformance(interfaces.NewLinkedListMap(), "链表实现", count)
	interfaces.TestPerformance(interfaces.NewHashMap(), "哈希表实现", count)
}
