package text

import "fmt"

//Symbol represents a single piece of the input
type Symbol struct {
	fileOffset uint32
	symOffset  uint32
	value      string
	symID      SymID
}

//SymID is the identifier for the symbol type
type SymID = uint32

//END_OF_THE_LINE is default SymID for line endings
const (
	NONE = 0 + iota
	END_OF_THE_LINE
	WHITESPACE
	ROUND_OPEN
	ROUND_CLOSE
	SQUARE_OPEN
	SQUARE_CLOSE
	CURLY_OPEN
	CURLY_CLOSE
	COMMA
	COLON
	SEMICOLON
	AT
	HASH
	OPERATOR
	NUMBER
	IDENTIFIER
	QUOTED_STRING
)

//SymbolOf create a Symbol from the provided parameters
func SymbolOf(fileOffset, symOffset uint32, value string, symID SymID) Symbol {
	return Symbol{fileOffset: fileOffset, symOffset: symOffset, value: value, symID: symID}
}

//SymbolEmpty create an empty Symbol from the file offset
func SymbolEmpty(fileOffset uint32) Symbol {
	return Symbol{fileOffset: fileOffset, symOffset: 0, value: "", symID: 0}
}

//ID is the identifier of the Symbol
func (s Symbol) ID() SymID {
	return s.symID
}

//Value is the string value of the symbol
func (s Symbol) Value() string {
	return s.value
}

func (s *Symbol) String() string {
	return fmt.Sprint("\\ ", s.value, " \\", " at ", s.fileOffset)
}
