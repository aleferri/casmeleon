package asm

//Context is the context interface for the assembler resolution
type Context interface {
	GuardSymbol(name string, x int, addr uint32, c Compilable)
	ClearAll()
	Refresh(sym Symbol)
	RetryList() []RetryQueue
	ByteSize() uint32
}
