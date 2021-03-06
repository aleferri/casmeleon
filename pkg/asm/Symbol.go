package asm

type Symbol interface {
	Address() uint32
	Value() int64
	Name() string
	IsDynamic() bool
}
