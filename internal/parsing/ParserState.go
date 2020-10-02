package parsing

import "github.com/aleferri/casmeleon/internal/text"

//Defaults for tags
const (
	DefaultIgnore = 0
	DefaultAppend = 1
)

//ParserState supply function to interact with user and to pull tokens from stream
//and push tokens to parsed stream
//Pull put new tokens on the parser stream and return the number of token read,
//if EOF is reached a TokenEOF must be added to the slice and the result must be 0
type ParserState interface {
	ReportError(wrong text.Token, err ParsingErrorType, msg string) //ReportError report an error on a token
	Pull(tokens *[]text.Token) int                                  //Pull see the above
	Push(tokens []text.Token, tag int)                              //Push parsed tokens on the internal collection
	Mark() (int, int)                                               //Mark return a key needed to assemble the CST
	ToCST(markT int, markL, tag int)                                //Build a CST from the parsed tokens from mark and tag it with tag
}
