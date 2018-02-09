package merkle

import (
	"fmt"
	"testing"
)

// TestMerkleTree ... help
func TestMerkleTree(t *testing.T) {
	var (
		dataBlock1 = "A"
		dataBlock2 = "B"
		dataBlock3 = "C"
		dataBlock4 = "D"
		dataBlock5 = "E"
		dataBlock6 = "F"
		// dataBlock7 = "G"
		// dataBlock8 = "H"
	)

	// Create leaf nodes
	leafNode1, err := NewLeaf([]byte(dataBlock1))
	if err != nil {
		t.Error(err)
	}

	leafNode2, err := NewLeaf([]byte(dataBlock2))
	leafNode3, err := NewLeaf([]byte(dataBlock3))
	leafNode4, err := NewLeaf([]byte(dataBlock4))
	leafNode5, err := NewLeaf([]byte(dataBlock5))
	leafNode6, err := NewLeaf([]byte(dataBlock6))
	// leafNode7, err := NewLeaf([]byte(dataBlock7))
	// leafNode8, err := NewLeaf([]byte(dataBlock8))

	// Build the merkle tree
	rootNode := BuildTree(leafNode1, leafNode2, leafNode3, leafNode4, leafNode5, leafNode6)

	fmt.Println("----------------")
	printTree(rootNode)

}

// printTree ... depth first traversal print of merkle tree
func printTree(n *Node) string {
	if n.IsLeaf() {
		fmt.Println(n.data)
		return n.data
	}

	if n.Left != nil {
		fmt.Println(n.data)
		printTree(n.Left)
	}

	if n.Right != nil {
		fmt.Println(n.data)
		printTree(n.Right)
	}
	return "test"
}
