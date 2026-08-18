package asm

import (
	"fmt"
	"strconv"

	"github.com/aleferri/casmvm/pkg/opcodes"
)

//DirectiveDepositSymbols is a .db or .dw whose values are not known at parse
//time, because some of them are labels. The values are resolved at assembly
//time and the item registers itself on every dynamic symbol it uses, so that a
//pass reassembles it when one of those symbols moves.
//
//The length never changes, one entry is always size bytes, so the item is
//invariant to its own address: what makes it dirty is a symbol, not a shift.
type DirectiveDepositSymbols struct {
	values    []Symbol
	size      uint32
	bigEndian bool
}

func MakeDepositSymbols(values []Symbol, size uint32, bigEndian bool) *DirectiveDepositSymbols {
	if size == 0 {
		size = 1
	}
	return &DirectiveDepositSymbols{values, size, bigEndian}
}

func (d *DirectiveDepositSymbols) Assemble(m opcodes.VM, addr uint32, index int, ctx Context) (uint32, []uint8, error) {
	bin := make([]uint8, 0, uint32(len(d.values))*d.size)
	for _, s := range d.values {
		if s.IsDynamic() {
			ctx.GuardSymbol(s.Name(), index, addr, d)
		}
		v := uint64(s.Value())
		if d.bigEndian {
			for b := d.size; b > 0; b-- {
				bin = append(bin, uint8(v>>(8*(b-1))))
			}
		} else {
			for b := uint32(0); b < d.size; b++ {
				bin = append(bin, uint8(v>>(8*b)))
			}
		}
	}
	return addr + uint32(len(bin)), bin, nil
}

func (d *DirectiveDepositSymbols) IsAddressInvariant() bool {
	return true
}

func (d *DirectiveDepositSymbols) String() string {
	if d.size == 1 {
		return ".db[" + strconv.Itoa(len(d.values)) + "]"
	}
	return ".d" + fmt.Sprint(d.size) + "[" + strconv.Itoa(len(d.values)) + "]"
}
