package asm

type Constant struct {
	value int64
}

func MakeConstant(value int64) *Constant {
	return &Constant{value}
}

func (c *Constant) Value() int64 {
	return c.value
}

func (c *Constant) Address() uint32 {
	return 0
}

func (c *Constant) IsDynamic() bool {
	return false
}

func (c *Constant) Name() string {
	return "Constant"
}
