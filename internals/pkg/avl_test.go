package pkg_test

import (
	"testing"

	"github.com/ygortgaleno/cloudwalk-test/internals/pkg"
)

type avlInsertTest struct {
	nodes                                     []int
	expectedRoot, expectedLeft, expectedRight *pkg.Node[int]
}

var avlInTests = []avlInsertTest{
	{[]int{1, 2, 3}, &pkg.Node[int]{Id: uint32(2)}, &pkg.Node[int]{Id: uint32(1), Data: 2}, &pkg.Node[int]{Id: uint32(3)}},
	{[]int{1, 3, 2}, &pkg.Node[int]{Id: uint32(2)}, &pkg.Node[int]{Id: uint32(1), Data: 2}, &pkg.Node[int]{Id: uint32(3)}},
	{[]int{3, 2, 1}, &pkg.Node[int]{Id: uint32(2)}, &pkg.Node[int]{Id: uint32(1), Data: 2}, &pkg.Node[int]{Id: uint32(3)}},
	{[]int{3, 1, 2}, &pkg.Node[int]{Id: uint32(2)}, &pkg.Node[int]{Id: uint32(1), Data: 2}, &pkg.Node[int]{Id: uint32(3)}},
	{[]int{1, 2, 1}, &pkg.Node[int]{Id: uint32(1)}, nil, &pkg.Node[int]{Id: uint32(2)}},
	{[]int{2, 1, 1}, &pkg.Node[int]{Id: uint32(2)}, &pkg.Node[int]{Id: uint32(1)}, nil},
}

type avlSearchTest struct {
	nodes    []int
	searchId uint32
	expected *pkg.Node[int]
}

var avlSearcTests = []avlSearchTest{
	{[]int{1, 2, 3}, 1, &pkg.Node[int]{Id: uint32(1)}},
	{[]int{1, 2, 3}, 2, &pkg.Node[int]{Id: uint32(2)}},
	{[]int{1, 2, 3}, 3, &pkg.Node[int]{Id: uint32(3)}},
	{[]int{1, 2, 3}, 4, nil},
}

func TestInsert(t *testing.T) {
	for _, test := range avlInTests {
		aTree := pkg.AVLTree[int]{}

		for _, val := range test.nodes {
			aTree.Insert(uint32(val), val)
		}

		if test.expectedRoot.Id != aTree.Root.Id {
			t.Errorf("expected %d, received %d", test.expectedRoot.Id, aTree.Root.Id)
		}

		if test.expectedRight == nil && (test.expectedRight != aTree.Root.Right) {
			t.Errorf("expected %v, received %v", test.expectedRight, aTree.Root.Right)
		} else if test.expectedRight != nil && (test.expectedRight.Id != aTree.Root.Right.Id) {
			t.Errorf("expected %v, received %v", test.expectedRight, aTree.Root.Right)
		}

		if test.expectedLeft == nil && (test.expectedLeft != aTree.Root.Left) {
			t.Errorf("expected %v, received %v", test.expectedLeft, aTree.Root.Left)
		} else if test.expectedLeft != nil && (test.expectedLeft.Id != aTree.Root.Left.Id) {
			t.Errorf("expected %v, received %v", test.expectedLeft, aTree.Root.Left)
		}
	}
}

func TestSearch(t *testing.T) {
	for _, test := range avlSearcTests {
		aTree := pkg.AVLTree[int]{}

		for _, val := range test.nodes {
			aTree.Insert(uint32(val), val)
		}

		expected := test.expected
		received := aTree.Search(test.searchId)

		if expected == nil && (expected != received) {
			t.Errorf("expected %v, received %v", expected, received)
		} else if expected != nil && (expected.Id != received.Id) {
			t.Errorf("expected %v, received %v", expected, received)
		}
	}
}
