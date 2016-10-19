package main

import (
	"fmt"
	"unsafe"
)

func main() {

	uf := &UnionFind{
		Nodes: map[string]*Node{},
	}

	uf.Add("1", "2")
	uf.Add("2", "3")
	uf.Add("3", "4")
	uf.Add("4", "5")

	uf.Add("a", "b")
	uf.Add("b", "c")

	fmt.Printf("Size of numbers %v\n", uf.getSize("4"))
	fmt.Printf("Size of letters %v\n", uf.getSize("b"))

	for key, val := range uf.Nodes {
		fmt.Printf("%s has parent %s\n", key, val.Parent)
	}

}

// Node is the node in a UnionFind structure. It counts the number of items
// below it, points to its parent, and has a list of its children
type Node struct {
	// We want the node count so we can promote the largest set to parent when
	// joining two sets. This helps us keep the depth of the tree small
	Count  int32
	Parent string

	// If we wanted to we could put other graph data in here.
	// For example, A count of fraudsters in the graph
}

// ChainFind is a union-find structure that also keeps a list of the members of
// the group for quick and easy retrival
type UnionFind struct {
	Nodes map[string]*Node
}

// Add adds a connection between two items to the map
func (uf *UnionFind) Add(a, b string) int32 {

	// Find the parent nodes of these two items. If the node is new, this mints a new node.
	firstNode, secondNode := uf.getParentNodeOrNew(a), uf.getParentNodeOrNew(b) // HL

	// Join the two nodes together
	var parent *Node
	if firstNode.Parent == secondNode.Parent {
		// The parents are the same so we are done
		return firstNode.Count

	} else if firstNode.Count > secondNode.Count {
		// We pick a parent by chosing the node with the highest count
		parent = uf.setParent(firstNode, secondNode)

	} else {
		// secondNode is the parent
		parent = uf.setParent(secondNode, firstNode)
	}

	return parent.Count
}

func (uf *UnionFind) getSize(id string) int {
	n := uf.getParentNodeOrNew(id)
	return int(n.Count)
}

func (uf *UnionFind) setParent(parent, child *Node) *Node {
	// Ensure we promote the graph metadata
	parent.Count += child.Count

	// Set the parent on the child
	child.Parent = parent.Parent

	return parent
}

// indexNode finds the node for the given id.
func (uf *UnionFind) indexNode(id string) *Node {
	node, ok := uf.Nodes[id]
	if !ok {
		node = newNode(id)
		uf.Nodes[id] = node
	}
	return node
}

// getParentNodeOrNew finds the node indexed by the given id, then chains up to its parent.
// as it does so it points nodes it encounters directly to the parent. It will mint a new node
// if the node is not found.
func (uf *UnionFind) getParentNodeOrNew(id string) *Node {
	// This will grab the node from the map or create a new node and put it in the map
	node := uf.indexNode(id) // HLnew
	// ...

	for node.Parent != id { // HLpro

		// `node` is not the parent as parents always point to themselves // HLpro
		// Set our ID to that of the parent of `node` (move up the graph) // HLpro
		id = node.Parent // HLpro

		// Get the new parent after moving up the graph // HLpro
		newParent := uf.indexNode(id) // HLpro

		// Push `node` up the graph by pointing it at the parents parent (skip the middle man) // HLpro
		node.Parent = newParent.Parent // HLpro

		// Set that parent as `node` and loop. // HLpro
		// We keep doing this until we are at the top of the graph // HLpro
		node = newParent // HLpro
	} // HLpro

	return node
}

func newNode(id string) *Node {
	n := &Node{}
	n.Count = 1
	n.Parent = id

	return n
}

// Not Used in this toy example
// ByteSliceToString is used when you really want to convert a slice
// of bytes to a string without incurring overhead. It is only safe
// to use if you really know the byte slice is not going to change
// in the lifetime of the string
func ByteSliceToString(bs []byte) string {
	// This is copied from runtime. It relies on the string
	// header being a prefix of the slice header!
	return *(*string)(unsafe.Pointer(&bs))
}
