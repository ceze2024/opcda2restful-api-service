/*
 * @Date: 2022-11-11 09:49:07
 * @LastEditors: 春贰
 * @gitee: https://gitee.com/chun22222222
 * @github: https://github.com/chun222
 * @Desc:
 * @LastEditTime: 2022-11-11 10:21:57
 * @FilePath: \opcConnector\system\util\opc\tree.go
 */
package opc

import "fmt"

//Tree creates an OPC browser representation
type Tree struct {
	Name     string
	Parent   *Tree
	Branches []*Tree
	Leaves   []Leaf
}

//Leaf contains the OPC tag and forms part of the Tree struct for the  OPC browser
type Leaf struct {
	Name string
	Tag  string
	Type int16
}

//ExtractBranchByName return substree with name
func ExtractBranchByName(tree *Tree, name string) *Tree {
	if tree.Name == name {
		return tree
	}
	for _, b := range tree.Branches {
		subtree := ExtractBranchByName(b, name)
		if subtree != nil {
			return subtree
		}
	}
	return nil
}

//CollectTags traverses tree and collects all tags in string slice
func CollectTags(tree *Tree) []string {
	collection := []string{}
	for _, l := range tree.Leaves {
		collection = append(collection, l.Tag)
	}
	for _, b := range tree.Branches {
		lowerCollection := CollectTags(b)
		collection = append(collection, lowerCollection...)
	}
	return collection
}

//PrettyPrint prints tree in a nice format
func PrettyPrint(tree *Tree) {
	fmt.Println(tree.Name)
	printSubtree(tree, 1)
}

// printSubtree is a recursive helper function to traverse the tree
func printSubtree(tree *Tree, level int) {
	space := ""
	for i := 0; i < level; i++ {
		space += "  "
	}
	for _, l := range tree.Leaves {
		fmt.Println(space, "-", l.Tag)
	}
	for _, b := range tree.Branches {
		fmt.Println(space, "+", b.Name)
		printSubtree(b, level+1)
	}
}
