package mymap

import (
	"echodb/internal/model"
)

//仿 slice

// ArrayEntry 数组条目
type ArrayEntry struct {
	key   int
	value *model.Person
}

// ArrayList 数组线性表实现的MyMap
type ArrayList struct {
	array  *[]ArrayEntry // 指向数组的指针
	len    int           // 当前元素数量
	cap    int           // 数组容量
	sorted bool          //标记是否已排序
}

func NewArrayList(capacity int) *ArrayList {
	if capacity <= 0 {
		capacity = 8
	}
	arr := make([]ArrayEntry, capacity) // 先创建slice
	return &ArrayList{
		array:  &arr, // 再取地址
		cap:    capacity,
		sorted: false,
	}
}

// Put 插入或更新元素（自动标记为未排序）
func (al *ArrayList) Put(key int, value *model.Person) {
	// 顺序查找是否已存在
	array := *al.array
	for i := 0; i < al.len; i++ {
		if array[i].key == key {
			array[i].value = value
			al.sorted = false // 更新后可能破坏排序
			return
		}
	}

	// 扩容检查
	if al.len >= al.cap {
		al.resize()
		// 重新获取数组引用，因为resize可能改变了地址
		array = *al.array
	}

	// 追加新元素
	array[al.len] = ArrayEntry{key, value}
	al.len++
	al.sorted = false // 新增元素后未排序
}

// resize 先扩容再复制到新数组
func (al *ArrayList) resize() {
	// 扩容
	newCap := al.cap * 2
	if newCap == 0 {
		newCap = 16
	}

	newArray := make([]ArrayEntry, newCap)
	oldArray := *al.array

	// 复制旧数据到新数组
	for i := 0; i < al.len; i++ {
		newArray[i] = oldArray[i]
	}

	al.array = &newArray
	al.cap = newCap
}

// Get 支持顺序查找和二分查找
func (al *ArrayList) Get(key int) *model.Person {
	array := *al.array
	if al.sorted {
		// 已排序时使用二分查找
		if entry := al.binarySearch(key); entry != nil {
			return entry.value
		}
	} else {
		// 未排序时顺序查找
		for i := 0; i < al.len; i++ {
			if array[i].key == key {
				return array[i].value
			}
		}
	}
	return nil
}

// 二分查找
func (al *ArrayList) binarySearch(key int) *ArrayEntry {
	array := *al.array
	left, right := 0, al.len-1
	for left <= right {
		mid := left + (right-left)/2
		if array[mid].key == key {
			return &array[mid]
		} else if array[mid].key < key {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return nil
}

func (al *ArrayList) Delete(key int) {
	array := *al.array
	for i := 0; i < al.len; i++ {
		if array[i].key == key {
			for j := i; j < al.len-1; j++ {
				array[j] = array[j+1]
			}
			// 清空最后一个元素
			array[al.len-1] = ArrayEntry{}
			al.len--
			return
		}
	}
}

func (al *ArrayList) Size() int {
	return al.len
}

func (al *ArrayList) Capacity() int {
	return al.cap
}
