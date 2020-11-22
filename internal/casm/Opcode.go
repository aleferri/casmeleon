package language

import "github.com/aleferri/casmeleon/pkg/text"

//Opcode declared in the assembly language
type Opcode struct {
	name   string       //opcode name
	format []text.SymID //opcode parameters format
	params []string     //opcode parameters name
}
