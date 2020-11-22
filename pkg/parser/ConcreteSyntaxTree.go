package parser

import "github.com/aleferri/casmeleon/pkg/text"

//CSTNode represent a Concrete Syntaxt Tree Node
type CSTNode interface {
	Children() []CSTNode
	Symbols() []text.Symbol
	ID() uint32
}

//CSTLeaf is leaf of the tree (a node without children)
type CSTLeaf struct {
	list []text.Symbol
	id   uint32
}

//BuildLeaf of the tree
func BuildLeaf(list []text.Symbol, id uint32) CSTNode {
	return &CSTLeaf{list: list, id: id}
}

//Children return the children of the leaf that by definition is an empty list
func (c *CSTLeaf) Children() []CSTNode {
	return []CSTNode{}
}

//Symbols return the symbols carried by the node
func (c *CSTLeaf) Symbols() []text.Symbol {
	return c.list
}

//ID return the id of the node
func (c *CSTLeaf) ID() uint32 {
	return c.id
}

//CSTBranch is a node of the tree that ha children
type CSTBranch struct {
	list     []text.Symbol
	id       uint32
	children []CSTNode
}

//BuildBranch of the tree
func BuildBranch(list []text.Symbol, id uint32) *CSTBranch {
	return &CSTBranch{list: list, id: id, children: []CSTNode{}}
}

//Children return the children of the node
func (c *CSTBranch) Children() []CSTNode {
	return c.children
}

//Symbols return the symbols carried by the node
func (c *CSTBranch) Symbols() []text.Symbol {
	return c.list
}

//ID return the id of the node
func (c *CSTBranch) ID() uint32 {
	return c.id
}

//InsertChild to the branch
func (c *CSTBranch) InsertChild(child CSTNode, ret bool) (CSTNode, bool) {
	if ret {
		c.children = append(c.children, child)
	}

	return c, ret
}
