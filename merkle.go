package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

//Content represents the data that is stored and verified by the tree.
//A type that implements this interface can be used as an item in the tree.
type Content interface {
	CalculateHash() []byte
	Data() []byte
}

//MerkleTree is the container for the tree. It holds a pointer to the root of the tree,.
type MerkleTree struct {
	Root       *Node   // pointer too root node
	merkleRoot []byte  //
	Leafs      []*Node // pointer too leaf nodes
}

//Node represents a node, root, or leaf in the tree.
type Node struct {
	Parent *Node   // parent pointer
	Left   *Node   // left child pointer
	Right  *Node   // right child pointer
	isLeaf bool    // ifs node a leaf
	dup    bool    // is node a duplicate
	Hash   []byte  // hash of node content
	C      Content // content of node
	Hex    string  // Hexidecimal representation of content hash
	data   []byte  // raw data before hash (useful for debug)
	proof  signature
}

// create a new merkle tree and return the root node, merkleroot and leaf nodes.
func newTree(cs []Content) (*MerkleTree, error) {
	if len(cs) == 0 {
		return nil, errors.New("error: cannot construct tree with no content")
	}

	var leafs []*Node
	for _, c := range cs {

		hash := c.CalculateHash()

		// create array of leaves
		leafs = append(leafs, &Node{
			Hash:   hash,
			Hex:    hex.EncodeToString(hash),
			C:      c,
			isLeaf: true,
			data:   c.Data(),
		})
	}

	// if un-even leaf amount, duplicate last leaf
	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash:   leafs[len(leafs)-1].Hash,
			C:      leafs[len(leafs)-1].C,
			isLeaf: true,
			dup:    true,
			Hex:    leafs[len(leafs)-1].Hex,
		}
		leafs = append(leafs, duplicate)
	}

	root := constructTree(leafs)

	t := &MerkleTree{
		Root:       root,
		merkleRoot: root.Hash,
		Leafs:      leafs,
	}

	return t, nil
}

//buildIntermediate is a helper function that for a given list of leaf nodes, constructs
//the intermediate and root levels of the tree. Returns the resulting root node of the tree.
func constructTree(nodes []*Node) *Node {
	var newNodes []*Node

	for i := 0; i < len(nodes); i += 2 {
		h := sha256.New()
		var left, right int = i, i + 1 // assign existing node to left/right child

		if i+1 == len(nodes) {
			right = i
		}

		chash := append(nodes[left].Hash, nodes[right].Hash...)
		h.Write(chash)
		hash := h.Sum(nil)

		// create new parent node and assign existing nodes as children
		n := &Node{
			Left:  nodes[left],
			Right: nodes[right],
			Hash:  hash,
			Hex:   hex.EncodeToString(hash),
			data:  []byte(string(nodes[left].data) + string(nodes[right].data)),
		}

		// assign new parent node to both children
		nodes[left].Parent = n
		nodes[right].Parent = n
		// add new parent node to new node list
		newNodes = append(newNodes, n)

		if len(nodes) == 2 {
			n.Parent = nil
			return n
		}
	}
	return constructTree(newNodes)
}

type anchor struct {
	SourceID string `json:"sourceId"`
	Type     string `json:"type"`
}

type proofpath struct {
	Right string `json:"right,omitempty"`
	Left  string `json:"left,omitempty"`
}

type signature struct {
	Context    []string    `json:"@context"`
	Type       string      `json:"type"`
	TargetHash string      `json:"targetHash"`
	MerkleRoot string      `json:"merkleRoot"`
	Anchors    []anchor    `json:"anchors"`
	Proof      []proofpath `json:"proof"`
}

func (t *MerkleTree) generateProofs() {
	for _, leaf := range t.Leafs {
		leaf.generateProof()
	}
}

func (n *Node) generateProof() {

	var signature = signature{
		Context: []string{"http://schema.org/", "https://w3id.org/security/v1"},
		Type:    "MerkleProof2017",
	}

	signature.TargetHash = hex.EncodeToString(n.Hash)

	root, path := buildPath(n, nil)

	signature.Proof = path
	signature.MerkleRoot = root

	n.proof = signature
}

// rename plz
func buildPath(n *Node, path []proofpath) (string, []proofpath) {
	var hash string
	if path == nil {
		path = []proofpath{}
	}

	if n.Parent == nil {
		hash = n.Hex
		return hash, path
	}

	previous := n
	current := n.Parent

	if current.Left == previous {
		path = append(path, proofpath{Right: current.Right.Hex})
	} else {
		path = append(path, proofpath{Left: current.Left.Hex})
	}

	return buildPath(n.Parent, path)
}

func verifyProof(s signature) bool {
	var decoded []byte
	currentHash, _ := hex.DecodeString(s.TargetHash)

	for _, lr := range s.Proof {
		h := sha256.New()
		if len(lr.Left) > 0 {
			decoded, _ = hex.DecodeString(lr.Left)
			h.Write(append(decoded, currentHash...))
		} else {
			decoded, _ = hex.DecodeString(lr.Right)
			h.Write(append(currentHash, decoded...))
		}

		currentHash = h.Sum(nil)
	}

	return hex.EncodeToString(currentHash) == s.MerkleRoot
}
