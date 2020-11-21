package exec

//Expression evaluable in stack (postifx)
type Expression struct {
	postfix []StackApplicable
}

//Resolve the expression
func (e *Expression) Resolve(params []uint32) uint32 {
	stack := BuildStack()
	for _, apply := range e.postfix {
		apply(&stack, params)
	}
	return stack.Pop()
}
