package casm

import (
	"strconv"
	"strings"

	"github.com/aleferri/casmeleon/pkg/parser"
)

//NumberBase is a base accepted into the language
type NumberBase struct {
	n      uint32
	prefix string
	suffix string
}

//PruneToNumBase create a NumberBase from a CSTNode
func PruneToNumBase(cst parser.CSTNode) (NumberBase, error) {
	tokens := cst.Symbols()
	n, err := strconv.ParseUint(tokens[1].Value(), 10, 64)
	countPrefix := len(tokens[2].Value()) - 1
	countSuffix := len(tokens[3].Value()) - 1
	return NumberBase{n: uint32(n), prefix: tokens[2].Value()[1:countPrefix], suffix: tokens[3].Value()[1:countSuffix]}, err
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
