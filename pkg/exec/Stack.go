package exec

//Stack is the standard stack structure
type Stack struct {
	content []int64
}

func (s *Stack) Content() []int64 {
	return s.content
}

//Push new content on top of the stack
func (s *Stack) Push(i int64) {
	s.content = append(s.content, i)
}

//Pop content from the top of the stack
func (s *Stack) Pop() int64 {
	l := len(s.content)
	val := s.content[l-1]
	s.content = s.content[:l-1]
	return val
}

//EmptyStack create an empty stack
func EmptyStack() Stack {
	return Stack{content: []int64{}}
}
