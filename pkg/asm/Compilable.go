package asm

type Compilable interface {
	Assemble(addr uint32, index int, ctx Context) (uint32, []uint8, error)
	IsAddressInvariant() bool
}
