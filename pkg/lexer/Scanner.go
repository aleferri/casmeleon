package lexer

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/aleferri/casmeleon/pkg/text"
)

//Scanner scan the runes and return the partial result and if the rune is accepted into the symbol
//The function should an empty result until it has received the last available symbol
//Initial state is zero and it should return the initial state when completed
type Scanner func(buffer []rune, stop bool) (id text.SymID, left []rune, err error)

//NUMBER_LIMIT_MAP map the second digit index to the char limit
var numberLimitMap = [...]int32{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 2, 4, 8, 16, 10}
var followSetA = map[rune]rune{'&': '&', '|': '|', '<': '<', '>': '>', '=': '!'}

//ScanBase10Number scan a number in base10
func ScanBase10Number(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	requireDigit := true
	for i := 0; i < len(buffer); i++ {
		x := buffer[i] - '0'
		isDigit := x >= 0 && x < 10

		if !isDigit {
			if requireDigit {
				return text.NONE, []rune{}, fmt.Errorf("Expected a digit in base 10, got '%c'", buffer[i])
			}
			if buffer[i] != '_' {
				return text.NUMBER, buffer[i:], nil
			}
		}
		requireDigit = !isDigit
	}

	if stop {
		if requireDigit {
			return text.NONE, []rune{}, errors.New("Unexpected end of input, expected a digit")
		}
		return text.NUMBER, []rune{}, nil
	}
	return text.NONE, buffer, nil
}

//ScanBase16Number scan a number in base 16
func ScanBase16Number(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	requireDigit := true
	for i := 2; i < len(buffer); i++ {
		x := buffer[i] - '0'
		a := buffer[i] - 'A'
		l := buffer[i] - 'a'
		isDigit := (x >= 0 && x < 10) || (a >= 0 && a < 6) || (l >= 0 && l < 6)

		if !isDigit {
			if requireDigit {
				return text.NONE, []rune{}, fmt.Errorf("Expected a digit in base 16, got '%c'", buffer[i])
			}
			if buffer[i] != '_' {
				return text.NUMBER, buffer[i:], nil
			}
		}
		requireDigit = !isDigit
	}

	if stop {
		if requireDigit {
			return text.NONE, []rune{}, errors.New("Unexpected end of input, expected a digit")
		}
		return text.NUMBER, []rune{}, nil
	}
	return text.NONE, buffer, nil
}

//ScanBase2Number scan a number in base 16
func ScanBase2Number(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	requireDigit := true
	for i := 2; i < len(buffer); i++ {
		x := buffer[i] - '0'
		isDigit := x < 2

		if !isDigit {
			if requireDigit {
				return text.NONE, []rune{}, fmt.Errorf("Expected a digit in base 2, got '%c'", buffer[i])
			}
			if buffer[i] != '_' {
				return text.NUMBER, buffer[i:], nil
			}
		}
		requireDigit = !isDigit
	}

	if stop {
		if requireDigit {
			return text.NONE, []rune{}, errors.New("Unexpected end of input, expected a digit")
		}
		return text.NUMBER, []rune{}, nil
	}
	return text.NONE, buffer, nil
}

//ScanNumber scan a number
func ScanNumber(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	if buffer[0] == '0' {
		if len(buffer) < 2 {
			return ScanBase10Number(buffer, stop)
		}
		if buffer[1] == 'b' || buffer[1] == 'B' {
			return ScanBase2Number(buffer, stop)
		}
		if buffer[1] == 'x' || buffer[1] == 'X' {
			return ScanBase16Number(buffer, stop)
		}
	}
	return ScanBase10Number(buffer, stop)
}

//ScanIdentifier scan an identifier
func ScanIdentifier(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	for i, c := range buffer {
		isUnder := (c == '_' || c == '$')
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) && !isUnder {
			return text.IDENTIFIER, buffer[i:], nil
		}
	}

	if stop {
		return text.IDENTIFIER, []rune{}, nil
	}
	return text.NONE, buffer, nil
}

//ScanSeparator scan a separator
func ScanSeparator(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	switch buffer[0] {
	case '(':
		return text.ROUND_OPEN, buffer[1:], nil
	case ')':
		return text.ROUND_CLOSE, buffer[1:], nil
	case '[':
		return text.SQUARE_OPEN, buffer[1:], nil
	case ']':
		return text.SQUARE_CLOSE, buffer[1:], nil
	case '{':
		return text.CURLY_OPEN, buffer[1:], nil
	case '}':
		return text.CURLY_CLOSE, buffer[1:], nil
	case ',':
		return text.COMMA, buffer[1:], nil
	case ';':
		return text.SEMICOLON, buffer[1:], nil
	case ':':
		return text.COLON, buffer[1:], nil
	case '@':
		return text.AT, buffer[1:], nil
	case '#':
		return text.HASH, buffer[1:], nil
	default:
		return text.NONE, buffer, errors.New("Not a separator")
	}
}

//ScanQuotedString scan a quoted string
func ScanQuotedString(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	for i := 1; i < len(buffer); i++ {
		if buffer[i] == '"' {
			return text.QUOTED_STRING, buffer[i+1:], nil
		}
	}

	if stop {
		return text.NONE, []rune{}, errors.New("Unexpected end of input")
	}
	return text.NONE, buffer, nil
}

//ScanOperatorSingle scan a single char operator
func ScanOperatorSingle(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	return text.OPERATOR, buffer[1:], nil
}

//ScanOperator scan an operator from the input
func ScanOperator(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	c := buffer[0]

	//COMPOSABLE
	if len(buffer) == 1 {
		if stop {
			return text.OPERATOR, buffer[1:], nil
		}
		return text.NONE, buffer, nil
	}
	s := buffer[1]
	if s == '=' {
		return text.OPERATOR, buffer[2:], nil
	}
	m, found := followSetA[s]
	if !found || m != c {
		return text.OPERATOR, buffer[1:], nil
	}

	if c == '>' && s == '>' && !stop {
		if len(buffer) == 2 {
			return text.NONE, buffer, nil
		}
		if buffer[2] == '>' {
			return text.OPERATOR, buffer[3:], nil
		}
	}

	return text.OPERATOR, buffer[2:], nil
}

//ScanWhite scan for a whitespace
func ScanWhite(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	return text.WHITESPACE, buffer[1:], nil
}

//ScanCommentLine recognize the line comment
func ScanCommentLine(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	for k, c := range buffer {
		if c == '\n' {
			return text.WHITESPACE, buffer[k:], nil
		}
	}
	if stop {
		return text.WHITESPACE, []rune{}, nil
	}
	return text.NONE, buffer, nil
}

//ScanCommentBlock recognize a comment block started by /* and terminated by */
func ScanCommentBlock(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	star := false
	for i, b := range buffer {
		if b == '/' && star {
			return text.WHITESPACE, buffer[i:], nil
		}
		star = b == '*'
	}
	return text.NONE, buffer, nil
}

//ScanSlash scan maybe a line comment, maybe a block comment, maybe a division operator
func ScanSlash(buffer []rune, stop bool) (id text.SymID, left []rune, err error) {
	if len(buffer) > 1 {
		if buffer[1] == '*' {
			return ScanCommentBlock(buffer, stop)
		} else if buffer[1] == '/' {
			return ScanCommentLine(buffer, stop)
		}
	}
	return ScanOperatorSingle(buffer, stop)
}
