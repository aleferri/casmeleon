package casm

import (
	"github.com/aleferri/casmeleon/pkg/parser"
)

// Inline is a temporary structure
type Inline struct {
	name   string   //name
	params []string //parameters name
	types  []uint32 //param types
}

// PruneToInline remove the header from the inline CST and return Inline and Body CST
func PruneToInline(lang *Language, op parser.CSTNode) (Inline, parser.CSTNode, error) {
	toks := op.Symbols()
	name := toks[1]
	children := op.Children()

	argsLUT, err := extractTypes(lang, children[0].Children())
	if err != nil {
		return Inline{}, nil, err
	}

	//Order of the parameters is the order of declaration: iterating argsLUT would
	//take it from a Go map, and the randomized iteration order would bind the
	//parameters differently at every run
	params := []string{}
	types := []uint32{}
	for _, arg := range children[0].Children() {
		name := arg.Symbols()[0].Value()
		params = append(params, name)
		types = append(types, argsLUT[name])
	}

	body := children[1]
	return Inline{name: name.Value(), params: params, types: types}, body, nil
}
