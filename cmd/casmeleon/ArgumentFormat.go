package main

import (
	"github.com/aleferri/casmeleon/pkg/asm"
)

type ArgumentFormat struct {
	types      []uint32
	format     []uint32
	parameters []asm.Symbol
}

func MakeFormat() ArgumentFormat {
	return ArgumentFormat{types: []uint32{}, format: []uint32{}, parameters: []asm.Symbol{}}
}
