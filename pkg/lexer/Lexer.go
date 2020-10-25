package lexer

import (
	"github.com/aleferri/casmeleon/pkg/text"
)

//Lexer structure
type Lexer struct {
	scanners map[rune]Scanner
	backup   Scanner
}

//Of Create a new Lexer from delimiters and symbols
func Of(scanners map[rune]Scanner, backup Scanner) Lexer {
	return Lexer{scanners: scanners, backup: backup}
}

//Scan next Symbol the the line
func (lexer *Lexer) Scan(fileOffset uint32, symOffset uint32, buffer []rune, stop bool) (*Result, error) {
	symbols := []text.Symbol{}

	for len(buffer) > 0 {
		val, ok := lexer.scanners[buffer[0]]
		var scanner Scanner
		if ok {
			scanner = val
		} else {
			scanner = lexer.backup
		}

		id, left, err := scanner(buffer, stop)
		if err != nil {
			return &Result{symbols: symbols, partial: buffer}, err
		}

		if id == text.NONE {
			return &Result{symbols: symbols, partial: buffer}, nil
		}

		endOffset := len(buffer) - len(left)

		if endOffset == 0 && stop {
			panic("Program Looped")
		}

		sym := text.SymbolOf(fileOffset, symOffset, string(buffer[0:endOffset]), id)
		symOffset++
		buffer = left
		symbols = append(symbols, sym)
	}

	return &Result{symbols: symbols, partial: buffer}, nil
}
