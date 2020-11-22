package parser

import "github.com/aleferri/casmeleon/pkg/text"

//Parser parse the source code to get a CST structure
type Parser func(stream Stream, report ErrorOf) (CSTNode, text.Error)

//ErrorOf report an error in the source code
type ErrorOf func(sym text.Symbol, src text.Source, msg string) text.Error
