package casm

import "github.com/aleferri/casmeleon/pkg/parser"

//Set of a symbol
type Set struct {
	name    string
	index   uint32
	valueOf func(string) int32
}

//Contains the specified symbol
func (s *Set) Contains(n string) bool {
	return s.valueOf(n) > -1
}

//Value of the specified symbol
func (s *Set) Value(n string) uint32 {
	return uint32(s.valueOf(n))
}

func generateFindValue(list []string) func(string) int32 {
	return func(str string) int32 {
		return 0
	}
}

//PruneToSet reduce the Concrete Syntax Tree Branch of a Set declaration to a type
func PruneToSet(node *parser.CSTBranch, index uint32) Set {
	name := node.Symbols()[1]
	values := []string{}
	for i := 2; i < len(node.Symbols()); i += 2 {
		values = append(values, node.Symbols()[i].Value())
	}
	return Set{name: name.Value(), index: index, valueOf: generateFindValue(values)}
}
