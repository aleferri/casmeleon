package main

import (
	"github.com/aleferri/casmeleon/internal/casm"
	"github.com/aleferri/casmeleon/pkg/asm"
	"github.com/aleferri/casmvm/pkg/opcodes"
)

type OpcodeInstance struct {
	addrInvariant bool
	name          string
	parameters    []asm.Symbol
	symTable      *SymbolTable
	atom          uint32
	index         int32
}

func MakeOpcodeInstance(opcode casm.Opcode, format ArgumentFormat, symTable *SymbolTable, atom uint32) *OpcodeInstance {
	inst := OpcodeInstance{
		addrInvariant: opcode.UseAddress(), name: opcode.Name(), parameters: format.parameters,
		symTable: symTable, atom: atom, index: opcode.Frame(),
	}
	return &inst
}

func (c *OpcodeInstance) Assemble(m opcodes.VM, addr uint32, index int, ctx asm.Context) (uint32, []uint8, error) {
	k := uint16(0)
	params := []uint16{}
	for _, a := range c.parameters {
		m.Frame().Values().Put(k, a.Value())
		if a.IsDynamic() {
			//Mark dynamic symbols only
			ctx.GuardSymbol(a.Name(), index, addr, c)
		}
		params = append(params, k)
		k++
	}

	m.Frame().Values().Put(k, int64(addr))

	_, err := m.Enter(c.index, params...)

	if err != nil {
		return addr, nil, err
	}

	outs := m.Frame().Returns()
	bin := []uint8{}
	size := uint16(outs.Size())
	for i := uint16(0); i < size; i++ {
		//TODO switch for atom size
		bin = append(bin, uint8(outs.Peek(i)))
	}
	return addr + uint32(size)/c.atom, bin, nil
}

func (c *OpcodeInstance) IsAddressInvariant() bool {
	return c.addrInvariant
}
