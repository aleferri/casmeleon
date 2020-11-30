package casm

import "github.com/aleferri/casmeleon/pkg/parser"

//ParserError is an encapsulated MatchError
type ParserError struct {
	internal *parser.MatchError
	context  string
}

func (e *ParserError) Error() string {
	return e.internal.Error()
}

//DecorateError of underlying match
func DecorateError(e error, s string) error {
	me, ok := e.(*parser.MatchError)
	if ok {
		return &ParserError{internal: me, context: s}
	}
	return e
}
