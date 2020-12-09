package asm

type Constant struct {
	value int64
}

func (c *Constant) Value() int64 {
	return c.value
}

func (c *Constant) Address() uint32 {
	return 0
}
