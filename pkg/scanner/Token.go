package scanner

//Token is a range of runes
type Token struct {
	slice   []rune
	basicID int32
}

func (t *Token) String() string {
	return string(t.slice)
}

//Merge ADJACENTS tokens, mind the capitalized ADJACENTS
//guaranteed memory corruption error otherwise
func (t Token) Merge(rhs Token) Token {
	lLen := len(t.slice)
	if lLen == 0 {
		return rhs
	}
	rLen := len(rhs.slice)
	if rLen == 0 {
		return t
	}
	slice := append([]rune{}, t.slice...)
	return Token{append(slice, rhs.slice...), 0}
}
