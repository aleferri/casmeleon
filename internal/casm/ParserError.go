package casm

import (
	"fmt"

	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

//ParserError is an encapsulated MatchError
type ParserError struct {
	wrapped *parser.MatchError
	context text.MessageContext
}

func (e *ParserError) Error() string {
	return e.wrapped.Error()
}

func (e *ParserError) PrettyPrint(source *text.Source) {
	wrong := e.wrapped.Found()
	fileName, lineIndex, column := source.FindPosition(wrong)
	fmt.Printf("In file %s, error at %d, %d: ", fileName, lineIndex+1, column+1)
	fmt.Printf(e.wrapped.Error(), wrong.Value(), e.wrapped.Expected().StringFromArray(idDescriptor))
	fmt.Println()
	source.PrintContext(e.context)
}

//WrapError of underlying match
func WrapMatchError(e error, left string, right string) error {
	me, ok := e.(*parser.MatchError)
	if ok {
		return &ParserError{wrapped: me, context: text.MakeMessageContext(me.Found(), left, right)}
	}
	panic("Unexpected kind of error from the source")
}
