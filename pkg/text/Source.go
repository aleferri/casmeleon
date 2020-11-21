package text

//Source represents a Source File that contains a list of symbols
type Source struct {
	fileName  string
	fileIndex uint32
	symbols   []Symbol
}

//BuildSource archive for error reporting
func BuildSource(fileName string) Source {
	return Source{fileName: fileName, symbols: []Symbol{}}
}

//Count the available symbols
func (s *Source) Count() uint32 {
	return uint32(len(s.symbols))
}

//FileIndex of the Source
func (s *Source) FileIndex() uint32 {
	return s.fileIndex
}

//AppendAll symbols to the Source
func (s *Source) AppendAll(syms []Symbol) []Symbol {
	s.symbols = append(s.symbols, syms...)
	notEmpty := []Symbol{}
	for _, sym := range syms {
		if sym.symID != WHITESPACE {
			notEmpty = append(notEmpty, sym)
		}
	}
	return notEmpty
}

//Append symbol to the Source
func (s *Source) Append(sym Symbol) {
	s.symbols = append(s.symbols, sym)
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

//FindPrevDelimiter find the previous delimiter from the symbol sym
func (s *Source) FindPrevDelimiter(sym Symbol, delimiter string) uint32 {
	offset := sym.symOffset
	offsetFound := false
	for !offsetFound {
		if s.symbols[offset].value == delimiter {
			offsetFound = true
		} else {
			offset--
			offsetFound = offset != 0
		}
	}
	return offset
}

//FindNextDelimiter find the previous delimiter from the symbol sym
func (s *Source) FindNextDelimiter(sym Symbol, delimiter string) uint32 {
	offset := sym.symOffset
	offsetFound := false
	for !offsetFound {
		if s.symbols[offset].value == delimiter {
			offsetFound = true
		} else {
			offset++
			offsetFound = (offset == uint32(len(s.symbols)))
		}
	}
	return offset
}
