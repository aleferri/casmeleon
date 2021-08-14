package main

import (
	"github.com/aleferri/casmeleon/internal/casm"
	"github.com/aleferri/casmeleon/pkg/asm"
	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/vm"
	"github.com/aleferri/casmvm/pkg/vmio"
)

type OpcodeInstance struct {
	addrInvariant bool
	name          string
	parameters    []asm.Symbol
	runList       []opcodes.Opcode
	symTable      *SymbolTable
	atom          uint32
}

func MakeOpcodeInstance(opcode casm.Opcode, format ArgumentFormat, symTable *SymbolTable, atom uint32) *OpcodeInstance {
	inst := OpcodeInstance{
		addrInvariant: opcode.UseAddress(), name: opcode.Name(), parameters: format.parameters, runList: opcode.RunList(), symTable: symTable, atom: atom,
	}
	return &inst
}

func (c *OpcodeInstance) Assemble(addr uint32, index int, ctx asm.Context) (uint32, []uint8, error) {
	instances := []int64{}
	for _, a := range c.parameters {
		instances = append(instances, a.Value())
		if a.IsDynamic() {
			//Mark dynamic symbols only
			ctx.GuardSymbol(a.Name(), index, addr, c)
		}
	}

	instances = append(instances, int64(addr))

	callable := vm.MakeCallable(c.runList)
	log := vmio.MakeVMLoggerConsole(vmio.ALL)
	vm := vm.MakeVerboseNaiveVM([]vm.Callable{}, log, vm.MakeVMFrame())
	err := vm.Run(callable, true)

	if err != nil {
		return addr, nil, err
	}

	outs := vm.Frame().Returns()
	bin := []uint8{}
	for i := uint16(0); i < uint16(outs.Size()); i++ {
		bin = append(bin, uint8(outs.Peek(i)))
	}
	return addr + uint32(len(bin))/c.atom, bin, nil
}

func (c *OpcodeInstance) IsAddressInvariant() bool {
	return c.addrInvariant
}
