package casm

import (
	"errors"

	"github.com/aleferri/casmeleon/pkg/exec"
	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

//Opcode declared in the assembly language
type Opcode struct {
	name    string            //opcode name
	format  []uint32          //opcode parameters format
	params  []string          //opcode parameters name
	types   []uint32          //param types
	runList []exec.Executable //executable operations
}

//Param Types
const (
	NUMBER = 0
	LABEL  = 1
)

//PruneToOpcode remove the header from the opcode CST and return Opcode and Body CST
func PruneToOpcode(lang *Language, op parser.CSTNode) (Opcode, parser.CSTNode, error) {
	toks := op.Symbols()
	name := toks[1]
	children := op.Children()

	argsLUT, err := extractTypes(lang, children[1].Children())
	if err != nil {
		return Opcode{}, nil, err
	}

	parsedFormat := children[0]
	argsFormat := []uint32{}

	for _, f := range parsedFormat.Symbols() {
		if f.ID() == text.Identifier {
			id, ok := argsLUT[f.Value()]
			if !ok {
				labels, _ := lang.SetByName("_FormatLabels")
				id = labels.index
			}
			argsFormat = append(argsFormat, id)
		} else {
			argsFormat = append(argsFormat, f.ID())
		}
	}

	params := []string{}
	types := []uint32{}
	for k, v := range argsLUT {
		params = append(params, k)
		types = append(types, v)
	}

	params = append(params, ".addr")
	types = append(types, 1)

	body := children[2]
	return Opcode{name: name.Value(), format: argsFormat, params: params, types: types}, body, nil
}

func extractTypes(lang *Language, args []parser.CSTNode) (map[string]uint32, error) {
	lut := map[string]uint32{}
	for _, a := range args {
		tokens := a.Symbols()
		name := tokens[0].Value()
		setName := tokens[2].Value()
		set, found := lang.SetByName(setName)
		if !found {
			return nil, errors.New("Set " + setName + " do not exists")
		}
		lut[name] = set.index
	}
	return lut, nil
}
