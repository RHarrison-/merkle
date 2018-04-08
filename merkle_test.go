package merkle

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"
)

// cert implements the Content interface provided by merkle and represents the content to bestored in the tree.
type cert struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Test string `json:"test"`
}

func (c cert) Data() []byte {
	b, _ := json.Marshal(c)
	return b
}

//CalculateHash hashes the value of a cert
func (c cert) CalculateHash() []byte {
	hashable := c.Data()
	h := sha256.New()
	h.Write(hashable)
	return h.Sum(nil)
}

var certs = []Content{
	cert{ID: "A", Name: "Reece", Test: "1"},
	cert{ID: "B", Name: "Kara", Test: "2"},
	cert{ID: "C", Name: "Maxx", Test: "3"},
	cert{ID: "D", Name: "Conor", Test: "4"},
}

func TestNewTree(t *testing.T) {
	tree, err := newTree(certs)
	if err != nil {
		t.Error("error: unexpected error:  ", err)
	}

	tree.generateProofs()

	fmt.Println("--------------------------")
	fmt.Println(verifyProof(tree.Leafs[1].proof))
	fmt.Println("--------------------------")
}
