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
func (s *Set) Value(n string) (uint32, bool) {
	v := s.valueOf(n)
	return uint32(v), v > -1
}

func generateFindValue(list []string) func(string) int32 {
	return func(str string) int32 {
		for i, s := range list {
			if s == str {
				return int32(i)
			}
		}
		return -1
	}
}

//PruneToSet reduce the Concrete Syntax Tree Branch of a Set declaration to a type
func PruneToSet(node parser.CSTNode, index uint32) Set {
	name := node.Symbols()[1]
	values := []string{}
	for i := 2; i < len(node.Symbols()); i += 2 {
		values = append(values, node.Symbols()[i].Value())
	}
	return Set{name: name.Value(), index: index, valueOf: generateFindValue(values)}
}
