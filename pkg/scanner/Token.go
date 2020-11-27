package scanner

//Token is a range of runes
type Token struct {
	slice []rune
}

func (t *Token) String() string {
	return string(t.slice)
}

//Merge ADJACENTS tokens, mind the capitalized ADJACENTS
//guaranteed memory corruption error otherwise
func (t *Token) Merge(rhs *Token) Token {
	lLen := len(t.slice)
	rLen := len(rhs.slice)
	return Token{t.slice[:lLen+rLen]}
}
