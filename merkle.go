package merkle

import (
	"crypto/sha256"
	"io"
)

// Node is the basic unit of the merkle tree
type Node struct {
	Parent, Left, Right *Node  // parent
	checksum            []byte // hashed data
	data                string // raw data (for testing)
}

// IsLeaf returns true if this is a leaf node (has no children)
func (n Node) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

// NewLeaf returns a new node from using the input data block
func NewLeaf(b []byte) (*Node, error) {
	n := new(Node)
	h := sha256.New()
	n.data = string(b)

	if _, err := io.WriteString(h, string(b)); err != nil {
		return nil, err
	}

	n.checksum = h.Sum(nil)
	return n, nil
}

// BuildTree builds a tree from leaf nodes and returns a root node
func BuildTree(theNodes ...*Node) *Node {
	var nodes []*Node

	for i := 0; i < len(theNodes); i = i + 2 {
		parentNode := new(Node)

		// only one child node left
		if i == len(theNodes)-1 {
			// parentNode.data = theNodes[i].data + theNodes[i].data
			parentNode.Left = theNodes[i]
			parentNode.Right = nil
			parentNode.Left.Parent = parentNode
			parentNode.checksum = HashMerkleBranch(parentNode.Left.checksum, parentNode.Left.checksum)

		} else {
			// parentNode.data = theNodes[i].data + theNodes[i+1].data
			parentNode.Left = theNodes[i] // bind tree
			parentNode.Right = theNodes[i+1]
			parentNode.Left.Parent = parentNode
			parentNode.Right.Parent = parentNode
			parentNode.checksum = HashMerkleBranch(parentNode.Right.checksum, parentNode.Left.checksum) // hash branch
		}

		// append to new node list
		nodes = append(nodes, parentNode)
	}

	// root node reached
	if len(nodes) == 1 {
		return nodes[0]
	}

	return BuildTree(nodes...)
}

// HashMerkleBranch ... hash two merkle branches
func HashMerkleBranch(left []byte, right []byte) []byte {
	h := sha256.New()

	if _, err := io.WriteString(h, string(left)+string(right)); err != nil {

	}

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	return h.Sum(nil)
}
