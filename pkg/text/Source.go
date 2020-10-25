package text

//Source represents a Source File that contains a list of symbols
type Source struct {
	fileName string
	symbols  []Symbol
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
