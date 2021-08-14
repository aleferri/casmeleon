package exec

import "github.com/aleferri/casmvm/pkg/opcodes"

//Function is a function made of opcodes
type Function struct {
	name string
	body []opcodes.Opcode
}

func MakeFunction(name string) *Function {
	return &Function{name, []opcodes.Opcode{}}
}

func (f *Function) AppendOpcode(opcode opcodes.Opcode) {
	f.body = append(f.body, opcode)
}
