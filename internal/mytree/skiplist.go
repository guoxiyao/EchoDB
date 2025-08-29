package mytree

import (
	"echodb/internal/model"
	"math/rand"
	"time"
)

// 跳表节点
type SkipListNode struct {
	Key    int
	Person *model.Person
	Next   []*SkipListNode // 多层指针
}

// 跳表实现
type SkipList struct {
	Head     *SkipListNode
	MaxLevel int
	Level    int
	size     int
}

// 创建新节点
func newSkipListNode(key int, person *model.Person, level int) *SkipListNode {
	return &SkipListNode{
		Key:    key,
		Person: person,
		Next:   make([]*SkipListNode, level),
	}
}

// 随机生成层级（概率为1/2的level次方）
func (s *SkipList) randomLevel() int {
	level := 1
	for rand.Float32() < 0.5 && level < s.MaxLevel {
		level++
	}
	return level
}

// 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// 创建新的跳表
func NewSkipList(maxLevel int) *SkipList {
	head := newSkipListNode(0, nil, maxLevel)
	return &SkipList{
		Head:     head,
		MaxLevel: maxLevel,
		Level:    1,
		size:     0,
	}
}

// 实现MyTree接口 - Put
func (s *SkipList) Put(id int, person *model.Person) {
	update := make([]*SkipListNode, s.MaxLevel)
	current := s.Head

	// 从最高层开始查找插入位置
	for i := s.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < id {
			current = current.Next[i]
		}
		update[i] = current
	}

	// 如果键已存在，更新值
	if current.Next[0] != nil && current.Next[0].Key == id {
		current.Next[0].Person = person
		return
	}

	// 生成新节点的随机层级
	newLevel := s.randomLevel()
	newNode := newSkipListNode(id, person, newLevel)

	// 如果新节点的层级大于当前跳表层级，更新高层级的指针
	if newLevel > s.Level {
		for i := s.Level; i < newLevel; i++ {
			update[i] = s.Head
		}
		s.Level = newLevel
	}

	// 更新各层指针
	for i := 0; i < newLevel; i++ {
		newNode.Next[i] = update[i].Next[i]
		update[i].Next[i] = newNode
	}

	s.size++
}

// 实现MyTree接口 - Get
func (s *SkipList) Get(id int) *model.Person {
	current := s.Head

	// 从最高层开始查找
	for i := s.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < id {
			current = current.Next[i]
		}
	}

	// 检查最底层是否找到
	if current.Next[0] != nil && current.Next[0].Key == id {
		return current.Next[0].Person
	}

	return nil
}

// 实现MyTree接口 - RangeQuery
func (s *SkipList) RangeQuery(minId, maxId int) []*model.Person {
	result := make([]*model.Person, 0)

	// 先找到大于等于minId的第一个节点
	startNode := s.findLowerBound(minId)
	if startNode == nil {
		return result
	}

	// 从起始节点开始线性遍历，直到超过maxId
	current := startNode
	for current != nil && current.Key <= maxId {
		result = append(result, current.Person)
		current = current.Next[0]
	}

	return result
}

// 找到大于等于minId的第一个节点
func (s *SkipList) findLowerBound(minId int) *SkipListNode {
	current := s.Head

	// 从最高层开始查找
	for i := s.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < minId {
			current = current.Next[i]
		}
	}

	// 返回第一个大于等于minId的节点
	return current.Next[0]
}

// 实现MyTree接口 - Delete
func (s *SkipList) Delete(id int) {
	update := make([]*SkipListNode, s.MaxLevel)
	current := s.Head

	// 从最高层开始查找要删除的节点
	for i := s.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < id {
			current = current.Next[i]
		}
		update[i] = current
	}

	// 找到要删除的节点
	target := current.Next[0]
	if target == nil || target.Key != id {
		return // 节点不存在
	}

	// 更新各层指针
	for i := 0; i < s.Level; i++ {
		if update[i].Next[i] != target {
			break
		}
		update[i].Next[i] = target.Next[i]
	}

	// 更新跳表层级（如果最高层变空）
	for s.Level > 1 && s.Head.Next[s.Level-1] == nil {
		s.Level--
	}

	s.size--
}

// 获取跳表大小
func (s *SkipList) Size() int {
	return s.size
}

// 获取跳表高度
func (s *SkipList) Height() int {
	return s.Level
}

// 打印跳表结构（用于调试）
func (s *SkipList) Print() {
	println("SkipList Structure:")
	println("==================")

	for i := s.Level - 1; i >= 0; i-- {
		print("Level ", i, ": ")
		current := s.Head.Next[i]
		for current != nil {
			print(current.Key, " -> ")
			current = current.Next[i]
		}
		println("nil")
	}
	println("==================")
}

// 验证跳表是否有效
func (s *SkipList) Validate() bool {
	// 检查每层是否有序
	for i := 0; i < s.Level; i++ {
		current := s.Head.Next[i]
		for current != nil && current.Next[i] != nil {
			if current.Key >= current.Next[i].Key {
				return false
			}
			current = current.Next[i]
		}
	}

	// 检查底层是否包含所有元素
	level0Count := 0
	current := s.Head.Next[0]
	for current != nil {
		level0Count++
		current = current.Next[0]
	}

	return level0Count == s.size
}

// 清空跳表
func (s *SkipList) Clear() {
	for i := 0; i < s.MaxLevel; i++ {
		s.Head.Next[i] = nil
	}
	s.Level = 1
	s.size = 0
}

// 获取所有键值（按顺序）
func (s *SkipList) Keys() []int {
	keys := make([]int, 0, s.size)
	current := s.Head.Next[0]
	for current != nil {
		keys = append(keys, current.Key)
		current = current.Next[0]
	}
	return keys
}

// 范围查询的迭代器版本（更高效）
func (s *SkipList) RangeQueryIter(minId, maxId int) <-chan *model.Person {
	ch := make(chan *model.Person)
	go func() {
		defer close(ch)
		startNode := s.findLowerBound(minId)
		if startNode == nil {
			return
		}
		current := startNode
		for current != nil && current.Key <= maxId {
			ch <- current.Person
			current = current.Next[0]
		}
	}()
	return ch
}

// 统计各层节点数量（用于分析）
func (s *SkipList) LevelStats() []int {
	stats := make([]int, s.Level)
	for i := 0; i < s.Level; i++ {
		count := 0
		current := s.Head.Next[i]
		for current != nil {
			count++
			current = current.Next[i]
		}
		stats[i] = count
	}
	return stats
}