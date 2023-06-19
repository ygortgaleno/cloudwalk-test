package pkg

import (
	"fmt"
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
func (t *AVLTree[T]) Insert(id uint32, data T) {
	t.sync.Lock()
	defer t.sync.Unlock()
	t.Root = t.Root.insert(id, data)
}

// Searches node element in tree by id
func (t *AVLTree[T]) Search(id uint32) (node *Node[T]) {
	t.sync.RLock()
	defer t.sync.RUnlock()
	return t.Root.search(id)
}

// Prints tree in-order
func (t *AVLTree[T]) Print() {
	t.sync.RLock()
	defer t.sync.RUnlock()
	t.Root.print()
}

func (n *Node[T]) insert(id uint32, data T) *Node[T] {
	if n == nil {
		return &Node[T]{id, data, 1, nil, nil}
	}

	if id < n.Id {
		n.Left = n.Left.insert(id, data)
	} else if id > n.Id {
		n.Right = n.Right.insert(id, data)
	} else {
		return n
	}

	return n.balanceTree()
}

func (n *Node[T]) search(id uint32) *Node[T] {
	if n == nil {
		return nil
	}
	if id < n.Id {
		return n.Left.search(id)
	} else if id > n.Id {
		return n.Right.search(id)
	} else {
		return n
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
	if node == nil {
		return node
	}
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

func (currentNode *Node[T]) rotateLeft() *Node[T] {
	newRoot := currentNode.Right
	currentNode.Right = newRoot.Left
	newRoot.Left = currentNode

	currentNode.calcNewHeight()
	newRoot.calcNewHeight()
	return newRoot
}

func (currentNode *Node[T]) rotateRight() *Node[T] {
	newRoot := currentNode.Left
	currentNode.Left = newRoot.Right
	newRoot.Right = currentNode

	currentNode.calcNewHeight()
	newRoot.calcNewHeight()
	return newRoot
}

func (n *Node[T]) print() {
	if n == nil {
		return
	}

	n.Left.print()
	fmt.Println(n.Data, n.Id)
	n.Right.print()
}

func getMax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
