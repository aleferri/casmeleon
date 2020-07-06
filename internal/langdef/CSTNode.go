package langdef

import "bitbucket.org/mrpink95/casmeleon/internal/text"

//CSTNode is a Concrete Syntax Tree node built by the parser
type CSTNode struct {
	tag      int //tag is the parse name for the
	parent   *CSTNode
	children []CSTNode
	content  []text.Token
	flags    int //flags is a generic container for any flags needed
}

//NewCSTNode create a new CSTNode with a tag, a parent, a content and some flags
func NewCSTNode(tag int, parent *CSTNode, content []text.Token, flags int) CSTNode {
	return CSTNode{tag, parent, []CSTNode{}, content, flags}
}

//NewEmptyCSTNode same as NewCSTNode but without the content, assumed to be empty
func NewEmptyCSTNode(tag int, parent *CSTNode, flags int) CSTNode {
	return CSTNode{tag, parent, []CSTNode{}, []text.Token{}, flags}
}

//NewRootCSTNode return a new CSTNode without parent
func NewRootCSTNode(tag int, flags int) CSTNode {
	return CSTNode{tag, nil, []CSTNode{}, []text.Token{}, flags}
}

//AddChild add a child to the children of the CSTNode
func (node *CSTNode) AddChild(child CSTNode) {
	node.children = append(node.children, child)
}

//Root return the current root of the tree
func (node *CSTNode) Root() CSTNode {
	if node.parent == nil {
		return *node
	}
	return node.parent.Root()
}
