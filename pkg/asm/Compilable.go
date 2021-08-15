package asm

import "github.com/aleferri/casmvm/pkg/opcodes"

type Compilable interface {
	Assemble(vm opcodes.VM, addr uint32, index int, ctx Context) (uint32, []uint8, error)
	IsAddressInvariant() bool
}
