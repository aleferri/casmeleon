package parser

import "github.com/aleferri/casmeleon/pkg/text"

//Stream of the tokens
type Stream interface {
	Next() text.Symbol
	Peek() text.Symbol
}
