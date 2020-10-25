package parser

import "github.com/aleferri/casmeleon/pkg/text"

//CSTNode represent a Concrete Syntaxt Tree Node
type CSTNode interface {
	Children() []CSTNode
	ID() uint32
}

//CSTLeaf is leaf of the tree (a node without children)
type CSTLeaf struct {
	list text.Symbol
	id   uint32
}

//Children return the children of the leaf that by definition is an empty list
func (c *CSTLeaf) Children() []CSTNode {
	return []CSTNode{}
}

//ID return the id of the node
func (c *CSTLeaf) ID() uint32 {
	return c.id
}

//CSTBranch is a node of the tree that ha children
type CSTBranch struct {
	list     text.Symbol
	id       uint32
	children []CSTNode
}

//Children return the children of the node
func (c *CSTBranch) Children() []CSTNode {
	return c.children
}

//ID return the id of the node
func (c *CSTBranch) ID() uint32 {
	return c.id
}
