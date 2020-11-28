package casm

import "strings"

//NumberBase is a base accepted into the language
type NumberBase struct {
	n      uint32
	prefix string
	suffix string
}

//DIGITS of a number
const DIGITS = "0123456789ABCDEF"

//Base of the number base
func (b *NumberBase) Base() uint32 {
	return b.n
}

//Parse the number in the specified base
func (b *NumberBase) Parse(s string) (uint32, bool) {
	isBase := strings.HasSuffix(s, b.suffix) && strings.HasPrefix(s, b.prefix)
	num := uint32(0)
	if isBase {
		rm := strings.TrimSuffix(strings.TrimPrefix(s, b.prefix), b.suffix)
		list := []byte(rm)
		for i, c := range list {
			num += (uint32(len(list)-i) * b.n * uint32(c))
		}
	}
	return num, isBase
}
