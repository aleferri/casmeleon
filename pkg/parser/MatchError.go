package parser

import "github.com/aleferri/casmeleon/pkg/text"

//ERROR CLASSES
const (
	NOMATCH = 0 + iota
)

//MatchError is the error found during match
type MatchError struct {
	expected ExpectedOptions
	found    text.Symbol
	message  string
	errCode  uint32
}

func (m *MatchError) Error() string {
	return m.message
}

//Expected symbol
func (m *MatchError) Expected() ExpectedOptions {
	return m.expected
}

//Found symbol
func (m *MatchError) Found() text.Symbol {
	return m.found
}

//ExpectedAnyOf following the list of allowed symbols
func ExpectedAnyOf(wrong text.Symbol, message string, list ...uint32) *MatchError {
	return &MatchError{expected: MakeExpectedAny(list...), found: wrong, message: message, errCode: NOMATCH}
}

//ExpectedSymbol return an error for an unmatched expected symbol
func ExpectedSymbol(wrong text.Symbol, message string, kind uint32) *MatchError {
	return &MatchError{expected: MakeExpectedKind(kind), found: wrong, message: message, errCode: NOMATCH}
}
