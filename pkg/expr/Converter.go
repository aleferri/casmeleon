package expr

import (
	"math"

	"github.com/aleferri/casmeleon/pkg/text"
)

type Converter struct {
	q         []text.Symbol
	stack     []Atom
	flags     uint16
	nextLocal uint16
}

func MakeConverter(q []text.Symbol, nParams uint16) Converter {
	return Converter{q, []Atom{}, 0, nParams}
}

func (c *Converter) Queue() []text.Symbol {
	return c.q
}

func (c *Converter) DropFront(n uint) {
	c.q = c.q[n:]
}

func (c *Converter) IsEmptyQueue() bool {
	return len(c.q) == 0
}

func (c *Converter) Front() text.Symbol {
	return c.q[0]
}

func (c *Converter) Poll() text.Symbol {
	ret := c.q[0]
	c.q = c.q[1:]
	return ret
}

func (c *Converter) LabelAtom(atom Atom) Atom {
	if atom.local == math.MaxUint16 {
		atom.local = c.LabelLocal()
	}
	return atom
}

func (c *Converter) LabelLocal() uint16 {
	local := c.nextLocal
	c.nextLocal++
	return local
}

func (c *Converter) Pop() Atom {
	last := len(c.stack) - 1
	atom := c.stack[last]
	c.stack = c.stack[0:last]
	return atom
}

func (c *Converter) Push(atom Atom) {
	c.stack = append(c.stack, atom)
}

func (c *Converter) IsEmptyStack() bool {
	return len(c.stack) == 0
}

func (c *Converter) SetFlag(flag uint16) {
	c.flags |= flag
}

func (c *Converter) HasFlag(flag uint16) bool {
	return (c.flags & flag) != 0
}
