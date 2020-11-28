package text

import "fmt"

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

//SliceLine return the Line sourrounding the symbol
func (s *Source) SliceLine(sym Symbol) []Symbol {
	return s.SliceScope(sym, "\n")
}

//Println surrounding the element
func (s *Source) Println(sym Symbol) {
	syms := s.SliceLine(sym)
	for _, p := range syms {
		fmt.Print(p.Value())
	}
	fmt.Println()
}

//SliceScope return the Scope that surround the symbol delimited by the specified delimiter
func (s *Source) SliceScope(sym Symbol, delimiter string) []Symbol {
	prev := s.FindPrevDelimiter(sym, delimiter)
	next := s.FindNextDelimiter(sym, delimiter)
	return s.symbols[prev:next]
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
