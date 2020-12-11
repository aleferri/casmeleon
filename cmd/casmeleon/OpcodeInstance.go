package main

import (
	"fmt"

	"github.com/aleferri/casmeleon/internal/casm"
	"github.com/aleferri/casmeleon/pkg/asm"
	"github.com/aleferri/casmeleon/pkg/exec"
)

type OpcodeInstance struct {
	addrInvariant bool
	parameters    []asm.Symbol
	runList       []exec.Executable
	symTable      *SymbolTable
}

func MakeOpcodeInstance(opcode casm.Opcode, format ArgumentFormat, symTable *SymbolTable) *OpcodeInstance {
	inst := OpcodeInstance{
		addrInvariant: opcode.UseAddress(), parameters: format.parameters, runList: opcode.RunList(), symTable: symTable,
	}
	return &inst
}

func (c *OpcodeInstance) Assemble(addr uint32, index int, ctx asm.Context) (uint32, []uint8, error) {

	instances := []int64{}
	for _, a := range c.parameters {
		instances = append(instances, a.Value())
		if a.IsDynamic() {
			//Mark dynamic symbols only
			fmt.Println("Guarding to", a.Name())
			ctx.GuardSymbol(a.Name(), index, addr, c)
		}
	}

	interp := exec.MakeInterpreter(exec.FrameOf(instances), c.runList)
	err := interp.Run()
	if err != nil {
		return addr, nil, err
	}

	outs := interp.PopResults()
	bin := []uint8{}
	for _, v := range outs.Content() {
		bin = append(bin, uint8(v))
	}
	return addr + uint32(len(bin)), bin, nil
}

func (c *OpcodeInstance) IsAddressInvariant() bool {
	return c.addrInvariant
}
