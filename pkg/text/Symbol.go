package text

import "fmt"

//Symbol represents a single piece of the input
type Symbol struct {
	fileOffset uint32
	symOffset  uint32
	value      string
	symID      uint32
}

//Reserved SymID
const (
	NONE = 0 + iota
	EOL
	EOF
	WHITESPACE
	RoundOpen
	RoundClose
	SquareOpen
	SquareClose
	CurlyOpen
	CurlyClose
	DoubleCurlyOpen
	DoubleCurlyClose
	Comma
	Colon
	Semicolon
	SymbolAt
	SymbolHash
	SymbolArrow
	CommentStart
	CommentEnd
	CommentLine
	QuotedString
	QuotedChar
	OperatorPlusUnary
	OperatorPlus
	OperatorMinus
	OperatorMinusUnary
	OperatorMul
	OperatorDiv
	OperatorMod
	OperatorRightShift
	OperatorLeftShift
	OperatorAnd
	OperatorLAnd
	OperatorOr
	OperatorLOr
	OperatorXor
	OperatorNot
	OperatorNeg
	OperatorLess
	OperatorLessEqual
	OperatorEqual
	OperatorGreaterEqual
	OperatorGreater
	OperatorNotEqual
	KeywordIF
	KeywordELSE
	KeywordOut
	KeywordSet
	KeywordNum
	KeywordInline
	KeywordOpcode
	KeywordWith
	KeywordExpr
	KeywordWarning
	KeywordError
	KeywordReturn
	Number
	Identifier
	ExactMatchKeyword
	LastReservedToken = ExactMatchKeyword
)

//SymbolOf create a Symbol from the provided parameters
func SymbolOf(fileOffset, symOffset uint32, value string, symID uint32) Symbol {
	return Symbol{fileOffset: fileOffset, symOffset: symOffset, value: value, symID: symID}
}

//SymbolEmpty create an empty Symbol from the file offset
func SymbolEmpty(fileOffset uint32) Symbol {
	return Symbol{fileOffset: fileOffset, symOffset: 0, value: "", symID: 0}
}

//WithID return the same symbol with different ID
func (s Symbol) WithID(symID uint32) Symbol {
	s.symID = symID
	return s
}

func (s Symbol) WithText(text string) Symbol {
	s.value = text
	return s
}

//ID is the identifier of the Symbol
func (s Symbol) ID() uint32 {
	return s.symID
}

//Value is the string value of the symbol
func (s Symbol) Value() string {
	return s.value
}

func (s *Symbol) String() string {
	return fmt.Sprint("\\ ", s.value, " \\", " at ", s.fileOffset)
}

//Equals between symbols
func (s Symbol) Equals(d Symbol) bool {
	return s.fileOffset == d.fileOffset && s.symOffset == d.symOffset
}
