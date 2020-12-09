package text

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

//Source represents a Source File that contains a list of symbols
type Source struct {
	fileName  string
	fileIndex uint32
	symbols   []Symbol
}

//BuildSource archive for error reporting
func BuildSource(fileName string) Source {
	return Source{fileName: fileName, fileIndex: 0, symbols: []Symbol{}}
}

//Count the available symbols
func (s *Source) Count() uint32 {
	return uint32(len(s.symbols))
}

//FileIndex of the Source
func (s *Source) FileIndex() uint32 {
	return s.fileIndex
}

//Append symbol to the Source
func (s *Source) Append(sym Symbol) {
	s.symbols = append(s.symbols, sym)
}

//FindPosition of a symbol inside the source
func (s *Source) FindPosition(sym Symbol) (string, uint32, uint32) {
	line := uint32(0)
	offset := uint32(0)
	for _, t := range s.symbols {
		fmt.Printf("token '%d' --- ", t.symID)
		if t.symOffset == sym.symOffset {
			return s.fileName, line, offset
		}
		if strings.Contains(t.value, "\n") {
			line++
			offset = 0
		} else {
			offset += uint32(len(t.value))
		}
	}
	return s.fileName, line, offset
}

//SliceLine return the Line sourrounding the symbol
func (s *Source) SliceLine(sym Symbol) []Symbol {
	return s.SliceScope(sym, "\n")
}

//SliceScope return the Scope that surround the symbol delimited by the specified delimiter
func (s *Source) SliceScope(sym Symbol, delimiter string) []Symbol {
	prev := s.FindPrevDelimiter(sym, delimiter)
	next := s.FindNextDelimiter(sym, delimiter)
	return s.symbols[prev:next]
}

func (s *Source) FindDelimiter(offset uint32, direction int, delimiter string) uint32 {
	maxLen := len(s.symbols)

	index := int(offset)

	for index > 0 && index < maxLen && s.symbols[index].value != delimiter {
		index += direction
	}

	return uint32(index)
}

//FindPrevDelimiter find the previous delimiter from the symbol sym
func (s *Source) FindPrevDelimiter(sym Symbol, delimiter string) uint32 {
	offset := sym.symOffset
	offsetFound := false
	for !offsetFound && offset > 0 {
		offsetFound = s.symbols[offset].value == delimiter
		offset--
	}
	return offset
}

//FindNextDelimiter find the previous delimiter from the symbol sym
func (s *Source) FindNextDelimiter(sym Symbol, delimiter string) uint32 {
	offset := sym.symOffset
	offsetFound := false
	limit := uint32(len(s.symbols))
	for !offsetFound && offset < limit {
		offsetFound = s.symbols[offset].value == delimiter
		offset++
	}
	return offset
}

//PrintContext of a message
func (s *Source) PrintContext(context MessageContext) {
	left := s.FindDelimiter(context.symOffset, -1, context.scopeLeft)
	startLine := s.FindDelimiter(context.symOffset-1, -1, "\n")
	endLine := s.FindDelimiter(context.symOffset+1, 1, "\n")
	right := s.FindDelimiter(context.symOffset, 1, context.scopeRight)

	if right < endLine {
		right = endLine
	}

	for _, t := range s.symbols[left:endLine] {
		fmt.Print(t.Value())
	}
	fmt.Println()

	if right == endLine {
		right = endLine + 1
	}

	offendedLine := s.symbols[startLine:endLine]

	fmt.Println("Symbols: ", len(offendedLine))

	for _, t := range offendedLine {
		char := ""
		if t.symOffset == context.symOffset {
			char = "^"
		} else if t.ID() > 2 {
			char = "-"
		}
		runeCount := utf8.RuneCountInString(t.Value())
		for k := 0; k < runeCount; k++ {
			fmt.Print(char)
		}
	}

	for _, e := range s.symbols[endLine+1 : right] {
		fmt.Print(e.Value())
	}
	fmt.Println()
}
