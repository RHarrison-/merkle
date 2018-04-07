package merkle

import "fmt"

// helper functions for working with and debugging merkle

func printTree(n *Node, detailed bool) {
	printed := false
	fmt.Println("-----------------------------")

	if n.isLeaf {
		// fmt.Println("Leaf")
		if !printed {
			printed = true

			fmt.Println(string(n.data))
			if detailed {
				// fmt.Println(string(n.Hash))
				fmt.Println(string(n.Hex))
			}
		}
		return
	}

	if n.Left != nil {
		if !printed {
			printed = true

			fmt.Println(string(n.data))
			if detailed {
				// fmt.Println(string(n.Hash))
				fmt.Println(string(n.Hex))
			}
		}
		printTree(n.Left, detailed)
	}

	if n.Right != nil {
		if !printed {
			printed = true
			fmt.Println(string(n.data))

			if detailed {
				// fmt.Println(string(n.Hash))
				fmt.Println(string(n.Hex))
			}

		}
		printTree(n.Right, detailed)
	}
}
