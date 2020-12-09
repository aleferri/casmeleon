package parser

//Parser parse the source code to get a CST structure
type Parser func(stream Stream) (CSTNode, error)
