package main

import (
	"github.com/aleferri/casmeleon/internal/langdef"
	"github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/internal/ui"
	"strings"
)

var nilASMLabel = ASMLabel{"", nil, nil, true, -1}

//ASMLabel is a label in the asm source code
type ASMLabel struct {
	name      string    //name of the label, i.e. loop, count, skip_loop
	reference *ASMLabel //global label reference, reference.name + name = full path, func.loop is full path for .loop local label
	section   ASMSource //section of the code with the label
	global    bool      //true if the label is global
	index     int       //index is the precedent opcode index
}

//NewASMLabel create a new ASMLabel with a name and the current status
func NewASMLabel(name string, index int, a *ASMReader) ASMLabel {
	local := strings.HasPrefix(name, ".")
	reference := &nilASMLabel
	if local {
		reference = a.lastGlobalLabel
	}
	return ASMLabel{name, reference, a.program, !local, index}
}

//FullName return full name of the label
func (p ASMLabel) FullName() string {
	return p.reference.name + p.name
}

//OpcodeInfo contains information to assemble the opcode
type OpcodeInfo struct {
	opcodeIndex int              //index of the opcode in the lang array
	argsName    []text.Token     //argsName are name of the arguments
	argsValue   []int            //argsValue are values of the arguments
	argsLabel   []bool           //argsLabel tells if the correspondent argsName is a label
	lang        *langdef.LangDef //lang is the language of the program
}

//NewOpcodeInfo return a new OpcodeAssemble
func NewOpcodeInfo(opcodeIndex int, argsName []text.Token, argsValue []int, lang *langdef.LangDef) OpcodeInfo {
	return OpcodeInfo{opcodeIndex, argsName, argsValue, nil, lang}
}

//Assemble the opcode
func (p OpcodeInfo) Assemble(flags uint) ([]uint, ui.SourceCodeError) {
	return p.lang.AssembleOpcode(p.opcodeIndex, p.argsName, p.argsValue, flags)
}

//String return opcode name with all his parameters
func (p OpcodeInfo) String() string {
	return p.lang.StringOpcode(p.opcodeIndex, p.argsName)
}
