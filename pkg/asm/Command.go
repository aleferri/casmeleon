package asm

import "github.com/aleferri/casmeleon/pkg/exec"

type Command struct {
	addrInvariant bool
	parameters    []Symbol
	dynamics      []Symbol
	runList       []exec.Executable
}

func (c *Command) Assemble(addr uint32, index int, ctx Context) (uint32, []uint8, error) {
	instances := []int64{}
	for _, a := range c.parameters {
		instances = append(instances, a.Value())
	}

	//Mark dynamic symbols only
	for _, d := range c.dynamics {
		ctx.GuardSymbol(d, index, addr, c)
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

func (c *Command) IsAddressInvariant() bool {
	return c.addrInvariant
}
