package pkg

import (
	"sync"
)

type AVLTree[T any] struct {
	Root *Node[T]
	sync sync.RWMutex
}

type Node[T any] struct {
	Id          uint32
	Data        T
	height      int
	Left, Right *Node[T]
}

// Insert new element into AVL Tree
func (tree *AVLTree[T]) Insert(id uint32, data T) {
	tree.sync.Lock()
	defer tree.sync.Unlock()
	tree.Root = tree.Root.insert(id, data)
}

// Searches node element in tree by id
func (tree *AVLTree[T]) Search(id uint32) (node *Node[T]) {
	tree.sync.RLock()
	defer tree.sync.RUnlock()
	return tree.Root.search(id)
}

func (node *Node[T]) insert(id uint32, data T) *Node[T] {
	if node == nil {
		return &Node[T]{id, data, 1, nil, nil}
	}

	if id < node.Id {
		node.Left = node.Left.insert(id, data)
	} else if id > node.Id {
		node.Right = node.Right.insert(id, data)
	} else {
		return node
	}

	return node.balanceTree()
}

func (node *Node[T]) search(id uint32) *Node[T] {
	if node == nil {
		return nil
	}
	if id < node.Id {
		return node.Left.search(id)
	} else if id > node.Id {
		return node.Right.search(id)
	} else {
		return node
	}
}

func (node *Node[T]) getHeight() int {
	if node == nil {
		return 0
	}
	return node.height
}

func (node *Node[T]) calcNewHeight() {
	node.height = 1 + getMax(node.Left.getHeight(), node.Right.getHeight())
}

func (node *Node[T]) balanceTree() *Node[T] {
	node.calcNewHeight()

	balanceFactor := node.Left.getHeight() - node.Right.getHeight()

	if balanceFactor < -1 {
		if node.Right.Left.getHeight() > node.Right.Right.getHeight() {
			node.Right = node.Right.rotateRight()
		}
		return node.rotateLeft()
	} else if balanceFactor > 1 {
		if node.Left.Right.getHeight() > node.Left.Left.getHeight() {
			node.Left = node.Left.rotateLeft()
		}
		return node.rotateRight()
	}

	return node
}

func (node *Node[T]) rotateLeft() *Node[T] {
	newRoot := node.Right
	node.Right = newRoot.Left
	newRoot.Left = node

	node.calcNewHeight()
	newRoot.calcNewHeight()
	return newRoot
}

func (node *Node[T]) rotateRight() *Node[T] {
	newRoot := node.Left
	node.Left = newRoot.Right
	newRoot.Right = node

	node.calcNewHeight()
	newRoot.calcNewHeight()
	return newRoot
}

func getMax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
