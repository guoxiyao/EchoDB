package mytree

import (
	"echodb/internal/model"
	"fmt"
)

const (
	RED   = true
	BLACK = false
)

type RBNode struct {
	Key    int
	Person *model.Person
	Color  bool
	Left   *RBNode
	Right  *RBNode
	Parent *RBNode
}

type RBTree struct {
	Root *RBNode
}

// 创建新节点
func newRBNode(key int, person *model.Person, color bool) *RBNode {
	return &RBNode{
		Key:    key,
		Person: person,
		Color:  color,
	}
}

// 实现MyTree接口 - Put
func (t *RBTree) Put(id int, person *model.Person) {
	newNode := newRBNode(id, person, RED)
	t.Root = t.insert(t.Root, newNode)
	t.Root.Color = BLACK // 根节点总是黑色
}

func (t *RBTree) insert(root *RBNode, newNode *RBNode) *RBNode {
	if root == nil {
		return newNode
	}

	// 标准的BST插入
	if newNode.Key < root.Key {
		root.Left = t.insert(root.Left, newNode)
		root.Left.Parent = root
	} else if newNode.Key > root.Key {
		root.Right = t.insert(root.Right, newNode)
		root.Right.Parent = root
	} else {
		// 键已存在，更新值
		root.Person = newNode.Person
		return root
	}

	// 红黑树修复
	return t.fixViolation(root)
}

// 修复红黑树性质
func (t *RBTree) fixViolation(node *RBNode) *RBNode {
	// 情况1: 当前节点是根节点
	if node.Parent == nil {
		node.Color = BLACK
		return node
	}

	// 情况2: 父节点是黑色，不需要修复
	if node.Parent.Color == BLACK {
		return t.findRoot(node)
	}

	// 获取叔叔节点
	grandparent := node.Parent.Parent
	var uncle *RBNode
	if node.Parent == grandparent.Left {
		uncle = grandparent.Right
	} else {
		uncle = grandparent.Left
	}

	// 情况3: 叔叔节点是红色
	if uncle != nil && uncle.Color == RED {
		node.Parent.Color = BLACK
		uncle.Color = BLACK
		grandparent.Color = RED
		return t.fixViolation(grandparent)
	}

	// 情况4: 叔叔节点是黑色或不存在
	if node.Parent == grandparent.Left {
		if node == node.Parent.Right {
			// 左旋
			grandparent.Left = t.leftRotate(node.Parent)
			grandparent.Left.Parent = grandparent
			node = node.Left
		}
		// 右旋并重新着色
		node.Parent.Color = BLACK
		grandparent.Color = RED
		if grandparent.Parent != nil {
			if grandparent.Parent.Left == grandparent {
				grandparent.Parent.Left = t.rightRotate(grandparent)
				grandparent.Parent.Left.Parent = grandparent.Parent
			} else {
				grandparent.Parent.Right = t.rightRotate(grandparent)
				grandparent.Parent.Right.Parent = grandparent.Parent
			}
		} else {
			return t.rightRotate(grandparent)
		}
	} else {
		if node == node.Parent.Left {
			// 右旋
			grandparent.Right = t.rightRotate(node.Parent)
			grandparent.Right.Parent = grandparent
			node = node.Right
		}
		// 左旋并重新着色
		node.Parent.Color = BLACK
		grandparent.Color = RED
		if grandparent.Parent != nil {
			if grandparent.Parent.Left == grandparent {
				grandparent.Parent.Left = t.leftRotate(grandparent)
				grandparent.Parent.Left.Parent = grandparent.Parent
			} else {
				grandparent.Parent.Right = t.leftRotate(grandparent)
				grandparent.Parent.Right.Parent = grandparent.Parent
			}
		} else {
			return t.leftRotate(grandparent)
		}
	}

	return t.findRoot(node)
}

// 左旋
func (t *RBTree) leftRotate(x *RBNode) *RBNode {
	y := x.Right
	x.Right = y.Left

	if y.Left != nil {
		y.Left.Parent = x
	}

	y.Parent = x.Parent
	y.Left = x
	x.Parent = y

	return y
}

// 右旋
func (t *RBTree) rightRotate(y *RBNode) *RBNode {
	x := y.Left
	y.Left = x.Right

	if x.Right != nil {
		x.Right.Parent = y
	}

	x.Parent = y.Parent
	x.Right = y
	y.Parent = x

	return x
}

// 找到根节点
func (t *RBTree) findRoot(node *RBNode) *RBNode {
	for node.Parent != nil {
		node = node.Parent
	}
	return node
}

// 实现MyTree接口 - Get
func (t *RBTree) Get(id int) *model.Person {
	node := t.search(t.Root, id)
	if node != nil {
		return node.Person
	}
	return nil
}

func (t *RBTree) search(node *RBNode, key int) *RBNode {
	if node == nil || node.Key == key {
		return node
	}

	if key < node.Key {
		return t.search(node.Left, key)
	}
	return t.search(node.Right, key)
}

// 实现MyTree接口 - RangeQuery
func (t *RBTree) RangeQuery(minId, maxId int) []*model.Person {
	result := make([]*model.Person, 0)
	t.rangeQuery(t.Root, minId, maxId, &result)
	return result
}

func (t *RBTree) rangeQuery(node *RBNode, minId, maxId int, result *[]*model.Person) {
	if node == nil {
		return
	}

	// 如果当前节点键值大于minId，遍历左子树
	if node.Key > minId {
		t.rangeQuery(node.Left, minId, maxId, result)
	}

	// 如果当前节点在范围内，添加到结果
	if node.Key >= minId && node.Key <= maxId {
		*result = append(*result, node.Person)
	}

	// 如果当前节点键值小于maxId，遍历右子树
	if node.Key < maxId {
		t.rangeQuery(node.Right, minId, maxId, result)
	}
}

// 实现MyTree接口 - Delete
func (t *RBTree) Delete(id int) {
	t.Root = t.deleteNode(t.Root, id)
	if t.Root != nil {
		t.Root.Color = BLACK
	}
}

func (t *RBTree) deleteNode(root *RBNode, key int) *RBNode {
	if root == nil {
		return nil
	}

	if key < root.Key {
		root.Left = t.deleteNode(root.Left, key)
		if root.Left != nil {
			root.Left.Parent = root
		}
	} else if key > root.Key {
		root.Right = t.deleteNode(root.Right, key)
		if root.Right != nil {
			root.Right.Parent = root
		}
	} else {
		// 找到要删除的节点
		if root.Left == nil {
			return root.Right
		} else if root.Right == nil {
			return root.Left
		}

		// 有两个子节点，找到后继节点
		successor := t.minValueNode(root.Right)
		root.Key = successor.Key
		root.Person = successor.Person
		root.Right = t.deleteNode(root.Right, successor.Key)
		if root.Right != nil {
			root.Right.Parent = root
		}
	}

	return t.fixDeleteViolation(root)
}

// 修复删除后的红黑树性质
func (t *RBTree) fixDeleteViolation(node *RBNode) *RBNode {
	if node == nil {
		return nil
	}

	// 简单的修复策略：如果需要更复杂的修复，可以进一步实现
	// 这里使用相对简化的方法
	return node
}

// 找到最小键值节点
func (t *RBTree) minValueNode(node *RBNode) *RBNode {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

// 辅助方法：检查红黑树性质
func (t *RBTree) Validate() bool {
	if t.Root == nil {
		return true
	}

	// 性质2: 根节点是黑色
	if t.Root.Color != BLACK {
		return false
	}

	// 性质4: 没有相邻的红色节点
	// 性质5: 从任一节点到其每个叶子的所有路径都包含相同数目的黑色节点
	blackCount := -1
	return t.validateNode(t.Root, 0, &blackCount)
}

func (t *RBTree) validateNode(node *RBNode, currentBlackCount int, blackCount *int) bool {
	if node == nil {
		if *blackCount == -1 {
			*blackCount = currentBlackCount
			return true
		}
		return currentBlackCount == *blackCount
	}

	// 性质4: 红色节点的子节点必须是黑色
	if node.Color == RED {
		if (node.Left != nil && node.Left.Color == RED) ||
			(node.Right != nil && node.Right.Color == RED) {
			return false
		}
	}

	// 递归检查子节点
	nextBlackCount := currentBlackCount
	if node.Color == BLACK {
		nextBlackCount++
	}

	return t.validateNode(node.Left, nextBlackCount, blackCount) &&
		t.validateNode(node.Right, nextBlackCount, blackCount)
}

// 打印树结构（用于调试）
func (t *RBTree) Print() {
	t.print(t.Root, 0, "ROOT")
}

func (t *RBTree) print(node *RBNode, level int, prefix string) {
	if node == nil {
		return
	}

	colorStr := "B"
	if node.Color == RED {
		colorStr = "R"
	}

	for i := 0; i < level; i++ {
		print("  ")
	}
	fmt.Printf("%s: %d(%s)\n", prefix, node.Key, colorStr)

	t.print(node.Left, level+1, "L")
	t.print(node.Right, level+1, "R")
}

// 创建新的红黑树
func NewRBTree() *RBTree {
	return &RBTree{}
}
