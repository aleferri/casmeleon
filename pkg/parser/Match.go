package parser

import (
	"errors"

	"github.com/aleferri/casmeleon/pkg/text"
)

//Accept symbol of type sym
func Accept(stream Stream, sym text.SymID) (text.Symbol, bool) {
	v := stream.Peek()
	if v.ID() == sym {
		stream.Next()
		return v, true
	}
	return v, false
}

//Consume symbol of type sym
func Consume(stream Stream, sym text.SymID) bool {
	_, c := Accept(stream, sym)
	return c
}

//Require accept a symbol of type sym or return an error if the symbol is not the accepted symbol type
func Require(stream Stream, sym text.SymID) (text.Symbol, error) {
	v, c := Accept(stream, sym)
	if c {
		return v, nil
	}
	return v, errors.New("Expected symbol")
}

//Expect consume a symbol of type sym or return an error if the symbol is not the accepted symbol type
func Expect(stream Stream, sym text.SymID) error {
	_, e := Require(stream, sym)
	return e
}

//AcceptAny accept any of the proposed symbols
func AcceptAny(stream Stream, syms ...text.SymID) (text.Symbol, bool) {
	peek := stream.Peek()
	for _, sym := range syms {
		if sym == peek.ID() {
			return stream.Next(), true
		}
	}
	return peek, false
}

//ConsumeAny symbol of type syms
func ConsumeAny(stream Stream, syms ...text.SymID) bool {
	_, c := AcceptAny(stream, syms...)
	return c
}

//RequireAny of the listed symbols
func RequireAny(stream Stream, syms ...text.SymID) (text.Symbol, error) {
	v, c := AcceptAny(stream, syms...)
	if !c {
		return v, errors.New("Expected one of symbols")
	}
	return v, nil
}

//ExpectAny of the listed symbols
func ExpectAny(stream Stream, syms ...text.SymID) error {
	_, err := RequireAny(stream, syms...)
	return err
}
