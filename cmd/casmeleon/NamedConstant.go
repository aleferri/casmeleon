package main

//NamedConstant is a symbol defined by .alias: a value bound to a name at parse
//time. Unlike a label it has no address and never changes, so it is invariant
//with respect to the assembler fixed point.
type NamedConstant struct {
	name  string
	value int64
}

func MakeNamedConstant(name string, value int64) *NamedConstant {
	return &NamedConstant{name: name, value: value}
}

func (c *NamedConstant) Name() string {
	return c.name
}

func (c *NamedConstant) Value() int64 {
	return c.value
}

func (c *NamedConstant) Address() uint32 {
	return 0
}

func (c *NamedConstant) IsDynamic() bool {
	return false
}
