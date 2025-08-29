package mytree

import (
	"echodb/internal/model"
	"fmt"
)

// B树节点
type BTreeNode struct {
	Keys     []int
	Persons  []*model.Person
	Children []*BTreeNode
	IsLeaf   bool
}

// B树实现
type BTree struct {
	Root   *BTreeNode
	Degree int // 最小度数
}

// 创建新的B树节点
func newBTreeNode(isLeaf bool, degree int) *BTreeNode {
	return &BTreeNode{
		Keys:     make([]int, 0, 2*degree-1),
		Persons:  make([]*model.Person, 0, 2*degree-1),
		Children: make([]*BTreeNode, 0, 2*degree),
		IsLeaf:   isLeaf,
	}
}

// 实现MyTree接口
func (t *BTree) Put(id int, person *model.Person) {
	if t.Root == nil {
		t.Root = newBTreeNode(true, t.Degree)
		t.Root.Keys = append(t.Root.Keys, id)
		t.Root.Persons = append(t.Root.Persons, person)
		return
	}

	// 如果根节点已满，需要分裂
	if len(t.Root.Keys) == 2*t.Degree-1 {
		oldRoot := t.Root
		newRoot := newBTreeNode(false, t.Degree)
		t.Root = newRoot
		newRoot.Children = append(newRoot.Children, oldRoot)
		t.splitChild(newRoot, 0)
	}
	t.insertNonFull(t.Root, id, person)
}

// 在非满节点中插入
func (t *BTree) insertNonFull(node *BTreeNode, id int, person *model.Person) {
	i := len(node.Keys) - 1

	if node.IsLeaf {
		// 叶子节点：找到插入位置并插入
		for i >= 0 && id < node.Keys[i] {
			i--
		}
		i++ // 调整到插入位置

		// 插入键值对
		node.Keys = append(node.Keys, 0)
		node.Persons = append(node.Persons, nil)
		copy(node.Keys[i+1:], node.Keys[i:])
		copy(node.Persons[i+1:], node.Persons[i:])
		node.Keys[i] = id
		node.Persons[i] = person
	} else {
		// 内部节点：找到合适的子节点
		for i >= 0 && id < node.Keys[i] {
			i--
		}
		i++ // 调整到子节点索引

		// 如果子节点已满，先分裂
		if len(node.Children[i].Keys) == 2*t.Degree-1 {
			t.splitChild(node, i)
			if id > node.Keys[i] {
				i++
			}
		}
		t.insertNonFull(node.Children[i], id, person)
	}
}

// 分裂子节点
func (t *BTree) splitChild(parent *BTreeNode, index int) {
	child := parent.Children[index]
	newChild := newBTreeNode(child.IsLeaf, t.Degree)

	// 移动后半部分键值到新节点
	mid := t.Degree - 1
	newChild.Keys = append(newChild.Keys, child.Keys[mid+1:]...)
	newChild.Persons = append(newChild.Persons, child.Persons[mid+1:]...)
	child.Keys = child.Keys[:mid]
	child.Persons = child.Persons[:mid]

	// 如果不是叶子节点，移动子节点指针
	if !child.IsLeaf {
		newChild.Children = append(newChild.Children, child.Children[mid+1:]...)
		child.Children = child.Children[:mid+1]
	}

	// 将中间键值提升到父节点
	parent.Keys = append(parent.Keys, 0)
	parent.Persons = append(parent.Persons, nil)
	copy(parent.Keys[index+1:], parent.Keys[index:])
	copy(parent.Persons[index+1:], parent.Persons[index:])
	parent.Keys[index] = child.Keys[mid]
	parent.Persons[index] = child.Persons[mid]

	// 更新子节点指针
	parent.Children = append(parent.Children, nil)
	copy(parent.Children[index+2:], parent.Children[index+1:])
	parent.Children[index+1] = newChild

	// 清理原节点的中间键值
	child.Keys = child.Keys[:mid]
	child.Persons = child.Persons[:mid]
}

// 查找键值
func (t *BTree) Get(id int) *model.Person {
	return t.search(t.Root, id)
}

func (t *BTree) search(node *BTreeNode, id int) *model.Person {
	if node == nil {
		return nil
	}

	// 在当前节点中查找
	i := 0
	for i < len(node.Keys) && id > node.Keys[i] {
		i++
	}

	if i < len(node.Keys) && id == node.Keys[i] {
		return node.Persons[i] // 找到
	}

	if node.IsLeaf {
		return nil // 叶子节点但没找到
	}

	// 递归在子节点中查找
	return t.search(node.Children[i], id)
}

// 范围查询实现
func (t *BTree) RangeQuery(minId, maxId int) []*model.Person {
	result := make([]*model.Person, 0)
	t.rangeQuery(t.Root, minId, maxId, &result)
	return result
}

func (t *BTree) rangeQuery(node *BTreeNode, minId, maxId int, result *[]*model.Person) {
	if node == nil {
		return
	}

	i := 0
	// 找到第一个可能包含[minId, maxId]的子节点
	for i < len(node.Keys) && minId > node.Keys[i] {
		i++
	}

	// 如果不是叶子节点，递归查询子节点
	if !node.IsLeaf {
		t.rangeQuery(node.Children[i], minId, maxId, result)
	}

	// 检查当前节点的键值
	for i < len(node.Keys) && node.Keys[i] <= maxId {
		if node.Keys[i] >= minId {
			*result = append(*result, node.Persons[i])
		}
		// 如果不是叶子节点，查询右子节点
		if !node.IsLeaf {
			t.rangeQuery(node.Children[i+1], minId, maxId, result)
		}
		i++
	}
}

// 创建新的B树
func NewBTree(degree int) *BTree {
	return &BTree{
		Degree: degree,
		Root:   nil,
	}
}

// 辅助方法：打印树结构（用于调试）
func (t *BTree) Print() {
	t.print(t.Root, 0)
}

func (t *BTree) print(node *BTreeNode, level int) {
	if node == nil {
		return
	}

	fmt.Printf("Level %d: ", level)
	for i, key := range node.Keys {
		fmt.Printf("%d ", key)
		if i < len(node.Keys)-1 {
			fmt.Printf("| ")
		}
	}
	fmt.Println()

	if !node.IsLeaf {
		for _, child := range node.Children {
			t.print(child, level+1)
		}
	}
}

// 获取树的高度
func (t *BTree) Height() int {
	return t.getHeight(t.Root)
}

func (t *BTree) getHeight(node *BTreeNode) int {
	if node == nil {
		return 0
	}
	if node.IsLeaf {
		return 1
	}
	return 1 + t.getHeight(node.Children[0])
}

func (t *BTree) delete(node *BTreeNode, id int) {
	// 步骤1：找到键所在位置
	i := 0
	for i < len(node.Keys) && id > node.Keys[i] {
		i++
	}

	// 情况1：键在当前节点中
	if i < len(node.Keys) && id == node.Keys[i] {
		if node.IsLeaf {
			// 情况1A：叶子节点直接删除
			node.Keys = append(node.Keys[:i], node.Keys[i+1:]...)
			node.Persons = append(node.Persons[:i], node.Persons[i+1:]...)
		} else {
			// 情况1B：内部节点需要替换前驱/后继
			t.deleteInternal(node, i)
		}
		return
	}

	// 情况2：键不在当前节点且是叶子节点（无需操作）
	if node.IsLeaf {
		return
	}

	// 情况3：递归删除前确保子节点有足够键
	if len(node.Children[i].Keys) < t.Degree {
		t.fillChild(node, i)
	}

	// 递归删除（调整i可能因合并而变化）
	if i > len(node.Keys) {
		i = len(node.Keys)
	}
	t.delete(node.Children[i], id)
}

// 处理内部节点删除
func (t *BTree) deleteInternal(node *BTreeNode, index int) {
	// 方法1：用前驱（左子树最大键）替换
	child := node.Children[index]
	if len(child.Keys) >= t.Degree {
		predecessor := t.getPredecessor(child)
		node.Keys[index] = predecessor.Keys[0]
		node.Persons[index] = predecessor.Persons[0]
		t.delete(child, predecessor.Keys[0])
		return
	}

	// 方法2：用后继（右子树最小键）替换
	child = node.Children[index+1]
	if len(child.Keys) >= t.Degree {
		successor := t.getSuccessor(child)
		node.Keys[index] = successor.Keys[0]
		node.Persons[index] = successor.Persons[0]
		t.delete(child, successor.Keys[0])
		return
	}

	// 方法3：合并子节点
	t.mergeChildren(node, index)
	t.delete(node.Children[index], node.Keys[index])
}

// 获取前驱节点（最大键）
func (t *BTree) getPredecessor(node *BTreeNode) *BTreeNode {
	for !node.IsLeaf {
		node = node.Children[len(node.Children)-1]
	}
	return &BTreeNode{
		Keys:    []int{node.Keys[len(node.Keys)-1]},
		Persons: []*model.Person{node.Persons[len(node.Persons)-1]},
	}
}

// 获取后继节点（最小键）
func (t *BTree) getSuccessor(node *BTreeNode) *BTreeNode {
	for !node.IsLeaf {
		node = node.Children[0]
	}
	return &BTreeNode{
		Keys:    []int{node.Keys[0]},
		Persons: []*model.Person{node.Persons[0]},
	}
}

// 填充键不足的子节点
func (t *BTree) fillChild(parent *BTreeNode, index int) {
	// 方法1：从左兄弟借键
	if index > 0 && len(parent.Children[index-1].Keys) >= t.Degree {
		t.borrowFromLeft(parent, index)
		return
	}

	// 方法2：从右兄弟借键
	if index < len(parent.Children)-1 && len(parent.Children[index+1].Keys) >= t.Degree {
		t.borrowFromRight(parent, index)
		return
	}

	// 方法3：合并兄弟节点
	if index != len(parent.Children)-1 {
		t.mergeChildren(parent, index)
	} else {
		t.mergeChildren(parent, index-1)
	}
}

// 从左兄弟借键
func (t *BTree) borrowFromLeft(parent *BTreeNode, index int) {
	child := parent.Children[index]
	leftSibling := parent.Children[index-1]

	// 父节点键下移
	child.Keys = append([]int{parent.Keys[index-1]}, child.Keys...)
	child.Persons = append([]*model.Person{parent.Persons[index-1]}, child.Persons...)

	// 左兄弟最大键上移
	parent.Keys[index-1] = leftSibling.Keys[len(leftSibling.Keys)-1]
	parent.Persons[index-1] = leftSibling.Persons[len(leftSibling.Persons)-1]

	// 移动子节点指针（若非叶子节点）
	if !leftSibling.IsLeaf {
		child.Children = append([]*BTreeNode{leftSibling.Children[len(leftSibling.Children)-1]}, child.Children...)
		leftSibling.Children = leftSibling.Children[:len(leftSibling.Children)-1]
	}

	// 移除左兄弟的键
	leftSibling.Keys = leftSibling.Keys[:len(leftSibling.Keys)-1]
	leftSibling.Persons = leftSibling.Persons[:len(leftSibling.Persons)-1]
}

// 从右兄弟借键（对称逻辑）
func (t *BTree) borrowFromRight(parent *BTreeNode, index int) {
	// 类似borrowFromLeft，方向相反
}

// 合并子节点
func (t *BTree) mergeChildren(parent *BTreeNode, index int) {
	left := parent.Children[index]
	right := parent.Children[index+1]

	// 父节点键下移
	left.Keys = append(left.Keys, parent.Keys[index])
	left.Persons = append(left.Persons, parent.Persons[index])

	// 合并右兄弟键
	left.Keys = append(left.Keys, right.Keys...)
	left.Persons = append(left.Persons, right.Persons...)

	// 合并子节点（若非叶子节点）
	if !left.IsLeaf {
		left.Children = append(left.Children, right.Children...)
	}

	// 移除父节点的键和指针
	parent.Keys = append(parent.Keys[:index], parent.Keys[index+1:]...)
	parent.Persons = append(parent.Persons[:index], parent.Persons[index+1:]...)
	parent.Children = append(parent.Children[:index+1], parent.Children[index+2:]...)
}
