package exec

import "github.com/aleferri/casmeleon/pkg/text"

//Opcode is the representation of a defined opcode
type Opcode interface {
	Params() []string
	Assemble(params []uint32) ([]uint8, Guard, text.Error)
}
