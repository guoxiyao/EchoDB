package mytree

import "echodb/internal/model"


type AVLNode struct {
	Key    int
	Person *model.Person
	Height int
	Left   *AVLNode
	Right  *AVLNode
}

type AVLTree struct {
	Root *AVLNode
}

func (t *AVLTree) Put(id int, person *model.Person) {
	t.Root = t.insert(t.Root, id, person)
}

func (t *AVLTree) insert(node *AVLNode, key int, person *model.Person) *AVLNode {
	if node == nil {
		return &AVLNode{
			Key:    key,
			Person: person,
			Height: 1,
		}
	}

	if key < node.Key {
		node.Left = t.insert(node.Left, key, person)
	} else if key > node.Key {
		node.Right = t.insert(node.Right, key, person)
	} else {
		// 键已存在，更新值
		node.Person = person
		return node
	}

	// 更新高度
	node.Height = 1 + max(t.height(node.Left), t.height(node.Right))

	// 获取平衡因子并调整
	balance := t.getBalance(node)

	// 左左情况
	if balance > 1 && key < node.Left.Key {
		return t.rightRotate(node)
	}

	// 右右情况
	if balance < -1 && key > node.Right.Key {
		return t.leftRotate(node)
	}

	// 左右情况
	if balance > 1 && key > node.Left.Key {
		node.Left = t.leftRotate(node.Left)
		return t.rightRotate(node)
	}

	// 右左情况
	if balance < -1 && key < node.Right.Key {
		node.Right = t.rightRotate(node.Right)
		return t.leftRotate(node)
	}

	return node
}

func (t *AVLTree) Get(id int) *model.Person {
	node := t.search(t.Root, id)
	if node != nil {
		return node.Person
	}
	return nil
}

func (t *AVLTree) Delete(id int) {
	t.Root = t.deleteNode(t.Root, id)
}

func (t *AVLTree) deleteNode(node *AVLNode, key int) *AVLNode {
	if node == nil {
		return nil
	}

	if key < node.Key {
		node.Left = t.deleteNode(node.Left, key)
	} else if key > node.Key {
		node.Right = t.deleteNode(node.Right, key)
	} else {
		// 找到要删除的节点
		if node.Left == nil || node.Right == nil {
			// 有一个或没有子节点
			var temp *AVLNode
			if node.Left != nil {
				temp = node.Left
			} else {
				temp = node.Right
			}

			if temp == nil {
				// 没有子节点
				return nil
			} else {
				// 有一个子节点
				node = temp
			}
		} else {
			// 有两个子节点：找到右子树的最小节点
			temp := t.minValueNode(node.Right)
			node.Key = temp.Key
			node.Person = temp.Person
			node.Right = t.deleteNode(node.Right, temp.Key)
		}
	}

	if node == nil {
		return nil
	}

	// 更新高度和平衡
	node.Height = 1 + max(t.height(node.Left), t.height(node.Right))
	balance := t.getBalance(node)

	// 重新平衡
	if balance > 1 && t.getBalance(node.Left) >= 0 {
		return t.rightRotate(node)
	}
	if balance > 1 && t.getBalance(node.Left) < 0 {
		node.Left = t.leftRotate(node.Left)
		return t.rightRotate(node)
	}
	if balance < -1 && t.getBalance(node.Right) <= 0 {
		return t.leftRotate(node)
	}
	if balance < -1 && t.getBalance(node.Right) > 0 {
		node.Right = t.rightRotate(node.Right)
		return t.leftRotate(node)
	}

	return node
}

func (t *AVLTree) minValueNode(node *AVLNode) *AVLNode {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

func (t *AVLTree) search(node *AVLNode, key int) *AVLNode {
	if node == nil || node.Key == key {
		return node
	}

	if key < node.Key {
		return t.search(node.Left, key)
	}
	return t.search(node.Right, key)
}

func (t *AVLTree) RangeQuery(minId, maxId int) []*model.Person {
	result := []*model.Person{}
	t.rangeQuery(t.Root, minId, maxId, &result)
	return result
}

func (t *AVLTree) rangeQuery(node *AVLNode, minId, maxId int, result *[]*model.Person) {
	if node == nil {
		return
	}

	// 只在需要时遍历子树
	if node.Key > minId {
		t.rangeQuery(node.Left, minId, maxId, result)
	}

	if node.Key >= minId && node.Key <= maxId {
		*result = append(*result, node.Person)
	}

	if node.Key < maxId {
		t.rangeQuery(node.Right, minId, maxId, result)
	}
}

// 辅助函数
func (t *AVLTree) height(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func (t *AVLTree) getBalance(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return t.height(node.Left) - t.height(node.Right)
}

func (t *AVLTree) leftRotate(x *AVLNode) *AVLNode {
	y := x.Right
	T2 := y.Left

	y.Left = x
	x.Right = T2

	x.Height = max(t.height(x.Left), t.height(x.Right)) + 1
	y.Height = max(t.height(y.Left), t.height(y.Right)) + 1

	return y
}

func (t *AVLTree) rightRotate(y *AVLNode) *AVLNode {
	x := y.Left
	T2 := x.Right

	x.Right = y
	y.Left = T2

	y.Height = max(t.height(y.Left), t.height(y.Right)) + 1
	x.Height = max(t.height(x.Left), t.height(x.Right)) + 1

	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
