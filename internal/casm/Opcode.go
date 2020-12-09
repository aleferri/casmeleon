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
	useAddr bool
}

func (o Opcode) Name() string {
	return o.name
}

func (o Opcode) Format() []uint32 {
	return o.format
}

func (o Opcode) StringifyFormat(lang *Language) []string {
	desc := []string{}

	ids := 0
	for _, particle := range o.format {
		if particle == text.Identifier {
			desc = append(desc, lang.sets[o.types[ids]].name)
			ids++
		} else {
			desc = append(desc, idDescriptor[particle])
		}
	}

	return desc
}

func (o Opcode) Accept(format []uint32, types []uint32) bool {
	return false
}

func StringifyFormat(format []uint32) []string {
	desc := []string{}

	for _, particle := range format {
		desc = append(desc, idDescriptor[particle])
	}

	return desc
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

	parsedFormat := children[0].Children()
	argsFormat := []uint32{}

	if len(parsedFormat) > 0 {
		for _, f := range parsedFormat[0].Symbols() {
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
