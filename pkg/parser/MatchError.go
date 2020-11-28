package parser

import "github.com/aleferri/casmeleon/pkg/text"

//ERROR CLASSES
const (
	NOMATCH = 0 + iota
)

//MatchError is the error found during match
type MatchError struct {
	expected uint32
	found    text.Symbol
	message  string
	errCode  uint32
}

func (m *MatchError) Error() string {
	return m.message
}

//Expected symbol
func (m *MatchError) Expected() uint32 {
	return m.expected
}

//Found symbol
func (m *MatchError) Found() text.Symbol {
	return m.found
}

//UnexpectedSymbol found
func UnexpectedSymbol(expected uint32, found text.Symbol, message string) error {
	return &MatchError{expected: expected, found: found, message: message, errCode: NOMATCH}
}

//MatchPatternError is the error found during pattern matching
type MatchPatternError struct {
	pattern []uint32
	found   []text.Symbol
	message string
	errCode uint32
}
