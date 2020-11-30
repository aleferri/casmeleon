package casm

//Opcode declared in the assembly language
type Opcode struct {
	name   string   //opcode name
	format []uint32 //opcode parameters format
	params []string //opcode parameters name
	types  []uint32 //param types
}
