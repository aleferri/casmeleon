package asm

import (
	"fmt"

	"github.com/aleferri/casmvm/pkg/opcodes"
)

var emptyLabelOutput = []uint8{}

type Label struct {
	name     string
	parent   *Label
	address  uint32
	byteSize uint32
}

func MakeLabel(name string, parent *Label, byteSize uint32) *Label {
	return &Label{name, parent, 0, byteSize}
}

//Assemble make the pass
func (l *Label) Assemble(m opcodes.VM, addr uint32, index int, ctx Context) (uint32, []uint8, error) {
	if addr != l.address {
		ctx.Refresh(l)
	}
	l.address = addr
	return addr, emptyLabelOutput, nil
}

func (l *Label) IsAddressInvariant() bool {
	return false
}

func (d *Label) Address() uint32 {
	return d.address / (d.byteSize / 8)
}

func (d *Label) Value() int64 {
	return int64(d.address / (d.byteSize / 8))
}

func (d *Label) Name() string {
	return d.name
}

func (d *Label) IsDynamic() bool {
	return true
}

func (d *Label) String() string {
	return d.name + "@" + fmt.Sprint(d.address)
}
