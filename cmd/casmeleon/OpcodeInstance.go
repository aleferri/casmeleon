package main

import (
	"fmt"

	"github.com/aleferri/casmeleon/internal/casm"
	"github.com/aleferri/casmeleon/pkg/asm"
	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/vmex"
)

type OpcodeInstance struct {
	name          string
	parameters    []asm.Symbol
	symTable      *SymbolTable
	atom          uint32
	bigEndian     bool
	invokeTarget  int32
	line          uint32
	addrInvariant bool
}

func MakeOpcodeInstance(opcode casm.Opcode, format ArgumentFormat, symTable *SymbolTable, atom uint32, bigEndian bool) *OpcodeInstance {
	if atom == 0 {
		atom = 1
	}
	//an opcode that reads .addr is exactly the one that is NOT invariant to it
	inst := OpcodeInstance{
		addrInvariant: !opcode.UseAddress(), name: opcode.Name(), parameters: format.parameters,
		symTable: symTable, atom: atom, bigEndian: bigEndian, invokeTarget: opcode.InvokeTarget(),
	}
	return &inst
}

func (c *OpcodeInstance) Assemble(m opcodes.VM, addr uint32, index int, ctx asm.Context) (uint32, []uint8, error) {
	k := uint16(0)

	frame := vmex.MakeVMFrame()

	for _, a := range c.parameters {
		frame.Values().Put(k, a.Value())
		k++
		if a.IsDynamic() {
			//Mark dynamic symbols only
			ctx.GuardSymbol(a.Name(), index, addr, c)
		}
	}

	frame.Values().Put(k, int64(addr/(ctx.ByteSize()/8)))

	err := m.Start(c.invokeTarget, &frame)

	if err != nil {
		return addr, nil, err
	}

	outs := frame.Returns()
	bin := []uint8{}
	size := uint16(outs.Size())
	for i := uint16(0); i < size; i++ {
		//Every value returned by .out is one atom, spread over atom bytes
		v := uint64(outs.Peek(i))
		if c.bigEndian {
			for b := c.atom; b > 0; b-- {
				bin = append(bin, uint8(v>>(8*(b-1))))
			}
		} else {
			for b := uint32(0); b < c.atom; b++ {
				bin = append(bin, uint8(v>>(8*b)))
			}
		}
	}
	return addr + uint32(size)*c.atom, bin, nil
}

func (c *OpcodeInstance) IsAddressInvariant() bool {
	return c.addrInvariant
}

func (c *OpcodeInstance) String() string {
	return fmt.Sprint(c.line) + ": " + c.name
}
