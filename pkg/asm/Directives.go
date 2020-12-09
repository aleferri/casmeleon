package asm

import (
	"errors"
)

type DirectiveOrg struct {
	address uint32
}

func (d *DirectiveOrg) Assemble(addr uint32, index int, ctx Context) (uint32, []uint8, error) {
	if addr > d.address {
		return 0, emptyLabelOutput, errors.New(".org directive cannot change PC backwards")
	}
	return d.address, emptyLabelOutput, nil
}

func (d *DirectiveOrg) IsAddressInvariant() bool {
	return false
}

func MakeOrg(target uint32) *DirectiveOrg {
	return &DirectiveOrg{target}
}

type DirectiveAdvance struct {
	address uint32
}

func (d *DirectiveAdvance) Assemble(addr uint32, index int, ctx Context) (uint32, []uint8, error) {
	if addr > d.address {
		return 0, emptyLabelOutput, errors.New(".advance directive cannot change PC backwards")
	}
	pad := make([]uint8, d.address-addr)
	return d.address, pad, nil
}

func (d *DirectiveAdvance) IsAddressInvariant() bool {
	return false
}

func MakeAdvance(target uint32) *DirectiveAdvance {
	return &DirectiveAdvance{target}
}

type DirectiveAlias struct {
	name  string
	value int64
}

func (d *DirectiveAlias) Address() uint32 {
	return 0
}

func (d *DirectiveAlias) Value() int64 {
	return d.value
}

func MakeAlias(name string, value int64) *DirectiveAlias {
	return &DirectiveAlias{name, value}
}

type DirectiveDeposit struct {
	binaryImage []uint8
}

func (d *DirectiveDeposit) Assemble(addr uint32, index int, ctx Context) (uint32, []uint8, error) {
	return addr + uint32(len(d.binaryImage)), d.binaryImage, nil
}

func (d *DirectiveDeposit) IsAddressInvariant() bool {
	return true
}

func MakeDeposit(values []uint8) *DirectiveDeposit {
	return &DirectiveDeposit{values}
}
