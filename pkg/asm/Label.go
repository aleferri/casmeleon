package asm

var emptyLabelOutput = []uint8{}

type Label struct {
	name    string
	parent  *Label
	address uint32
}

func MakeLabel(name string, parent *Label) *Label {
	return &Label{name, parent, 0}
}

//Assemble make the pass
func (l *Label) Assemble(addr uint32, ctx Context) (uint32, []uint8, error) {
	if addr != l.address {
		ctx.NotifyChange(l)
	}
	l.address = addr
	return addr, emptyLabelOutput, nil
}

func (l *Label) IsAddressInvariant() bool {
	return false
}

func (d *Label) Address() uint32 {
	return d.address
}

func (d *Label) Value() int64 {
	return int64(d.address)
}

func (d *Label) Name() string {
	return d.name
}
