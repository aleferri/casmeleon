package lexer

import "github.com/aleferri/casmeleon/pkg/text"

//Result represent the Lexer result
type Result struct {
	symbols []text.Symbol
	partial []rune
}

//Symbols scanned by the lexer
func (s *Result) Symbols() []text.Symbol {
	return s.symbols
}

//Unlexed runes scanned by the lexer
func (s *Result) Unlexed() []rune {
	return s.partial
}
