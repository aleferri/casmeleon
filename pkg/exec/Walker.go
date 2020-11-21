package exec

//Walker walk the tree of possibility to interpret the opcode
type Walker interface {
	Accept(p string) bool
	Walk(p string) (Walker, Opcode, bool)
}

//RootWalker is the top container for every opcode
type RootWalker struct {
	children []Walker
}

//Accept string p as continuation for the branch
func (w *RootWalker) Accept(current string) bool {
	return false
}

//Walk inside one path
func (w *RootWalker) Walk(next string) (Walker, Opcode, bool) {
	for _, c := range w.children {
		isPath := c.Accept(next)
		if isPath {
			return c.Walk("")
		}
	}
	return w, nil, false
}

//BranchWalker is the walker of branch of opcode
type BranchWalker struct {
	resolver func(a string) bool
	opcode   Opcode
	children []Walker
}

//Accept string p as continuation for the branch
func (w *BranchWalker) Accept(current string) bool {
	return w.resolver(current)
}

//Walk inside one path
func (w *BranchWalker) Walk(next string) (Walker, Opcode, bool) {
	opcode := w.opcode
	for _, c := range w.children {
		isPath := c.Accept(next)
		if isPath {
			return c.Walk("")
		}
	}
	return w, opcode, next == ""
}

//LeafWalker is the leaf of the walker tree
type LeafWalker struct {
	resolver func(a string) bool
	opcode   Opcode
}

//Accept string p as continuation for the branch
func (w *LeafWalker) Accept(current string) bool {
	return w.resolver(current)
}

//Walk return false unless the string is empty
func (w *LeafWalker) Walk(next string) (Walker, Opcode, bool) {
	return w, w.opcode, next == ""
}
