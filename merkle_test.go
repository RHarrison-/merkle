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

	var signatureJson = []byte(`{"@context":["http://schema.org/","https://w3id.org/security/v1"],"type":"MerkleProof2017","targetHash":"f272489e6ecf9e3ea7d575160325c630c4515dde753c1a1f84a33544ec43cacd","merkleRoot":"7389f388b1b9468c7ff9d27e33c73daf6c6a18ba25a2c33c1ea990a2127c5e82","anchors":null,"proof":[{"right":"fab8dfe93c886b0214d15525b072792722cbb66cfea4d04a663b25dad5190706"},{"right":"348a770d83cbe5c8d4ad11152a6aefa701f7b86c9392897b3b0a4a51241358a5"}]}`)

	var dat signature
	err = json.Unmarshal(signatureJson, &dat)

	tree.generateProofs()
	fmt.Println("--------------------------")
	fmt.Println(verifyProof(dat))
	fmt.Println("--------------------------")

}
