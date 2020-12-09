package parser

import (
	"github.com/aleferri/casmeleon/pkg/text"
)

//Match a generic rule in the provided stream
type Match func(stream Stream) (CSTNode, error)

//Accept symbol of type sym
func Accept(stream Stream, sym uint32) (text.Symbol, bool) {
	v := stream.Peek()
	if v.ID() == sym {
		return stream.Next(), true
	}
	return v, false
}

//Consume symbol of type sym
func Consume(stream Stream, sym uint32) bool {
	_, c := Accept(stream, sym)
	return c
}

//Require accept a symbol of type sym or return an error if the symbol is not the accepted symbol type
func Require(stream Stream, sym uint32) (text.Symbol, error) {
	v, c := Accept(stream, sym)
	if c {
		return v, nil
	}
	return v, ExpectedSymbol(v, "Unexpected '%s', was expecting: %s", sym)
}

//Expect consume a symbol of type sym or return an error if the symbol is not the accepted symbol type
func Expect(stream Stream, sym uint32) error {
	_, e := Require(stream, sym)
	return e
}

//AcceptAny accept any of the proposed symbols
func AcceptAny(stream Stream, syms ...uint32) (text.Symbol, bool) {
	peek := stream.Peek()
	for _, sym := range syms {
		if sym == peek.ID() {
			return stream.Next(), true
		}
	}
	return peek, false
}

//ConsumeAny symbol of type syms
func ConsumeAny(stream Stream, syms ...uint32) bool {
	_, c := AcceptAny(stream, syms...)
	return c
}

//RequireAny of the listed symbols
func RequireAny(stream Stream, syms ...uint32) (text.Symbol, error) {
	v, c := AcceptAny(stream, syms...)
	if !c {
		return v, ExpectedAnyOf(v, "Unexpected '%s', was expecting: %s", syms...)
	}
	return v, nil
}

//ExpectAny of the listed symbols
func ExpectAny(stream Stream, syms ...uint32) error {
	_, err := RequireAny(stream, syms...)
	return err
}

//RequireSequence specified by caller
func RequireSequence(stream Stream, seq ...uint32) ([]text.Symbol, error) {
	acc := []text.Symbol{}
	for _, m := range seq {
		item, err := Require(stream, m)
		if err != nil {
			return acc, err
		}
		acc = append(acc, item)
	}
	return acc, nil
}

//AcceptInsetPattern in a stream
func AcceptInsetPattern(stream Stream, left uint32, right uint32, seq ...uint32) ([]text.Symbol, error) {
	acc := []text.Symbol{}
	matchLeft := Expect(stream, left)
	if matchLeft != nil {
		return acc, matchLeft
	}
	for stream.Peek().ID() != right {
		for _, m := range seq {
			item, err := Require(stream, m)
			if err != nil {
				return acc, err
			}
			acc = append(acc, item)
		}
	}
	matchRight := Expect(stream, right)
	return acc, matchRight
}

//AcceptInsetDelegate read left expected symbol, then test for the right expected symbol.
//If symbol is not the specified symbol this function will delegate the inset matching on the provided function
//In the end the right symbol is read
func AcceptInsetDelegate(stream Stream, left uint32, right uint32, dg Match) ([]CSTNode, error) {
	leafs := []CSTNode{}
	matchLeft := Expect(stream, left)
	if matchLeft != nil {
		return leafs, matchLeft
	}
	for stream.Peek().ID() != right {
		part, err := dg(stream)
		if err != nil {
			return leafs, err
		}
		leafs = append(leafs, part)
	}
	matchRight := Expect(stream, right)
	return leafs, matchRight
}

//AcceptPatternWithTest in the stream source
func AcceptPatternWithTest(stream Stream, left uint32, right uint32, test uint32, dg Match) ([]CSTNode, error) {
	leafs := []CSTNode{}
	matchLeft := Expect(stream, left)
	if matchLeft != nil {
		return leafs, matchLeft
	}

	readNext := stream.Peek().ID() != right
	for readNext {
		part, err := dg(stream)
		if err != nil {
			return leafs, err
		}
		leafs = append(leafs, part)

		if stream.Peek().ID() == test {
			stream.Next()
		} else {
			readNext = false
		}
	}
	matchRight := Expect(stream, right)
	return leafs, matchRight
}
