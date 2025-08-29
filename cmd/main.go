package main

import (
	"echodb/internal/mymap"
	"fmt"
)

func main() {

	fmt.Println("开始MyMap测试...")
	// 测试正确性
	mymap.TestCorrectness()

	// 测试性能
	count := 10000
	fmt.Printf("\n开始MyMap性能测试（%d条数据）...\n\n", count)
	mymap.TestPerformance(mymap.NewArrayList(8), "数组实现", count)
	mymap.TestPerformance(mymap.NewLinkedList(), "链表实现", count)
	mymap.TestPerformance(mymap.NewHashTable(8), "哈希表实现", count)

	fmt.Println("开始MyTree测试...")

}
