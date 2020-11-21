package exec

//Stack structure for evaluation purpose
type Stack struct {
	content []uint32
}

//Push value into the stack
func (s *Stack) Push(i uint32) {
	s.content = append(s.content, i)
}

//Pop value into the stack
func (s *Stack) Pop() uint32 {
	len := len(s.content)
	val := s.content[len-1]
	s.content = s.content[:len-1]
	return val
}

//BuildStack build the Stack structure
func BuildStack() Stack {
	return Stack{content: []uint32{}}
}

//StackApplicable is function that can be applied in a stack
type StackApplicable func(s *Stack, values []uint32)
