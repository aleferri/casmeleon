package main

import "github.com/aleferri/casmeleon/pkg/asm"

type AssemblyProgram struct {
	list []asm.Compilable
}

func (a *AssemblyProgram) Add(c asm.Compilable) {
	a.list = append(a.list, c)
}

func MakeAssemblyProgram() AssemblyProgram {
	return AssemblyProgram{list: []asm.Compilable{}}
}
