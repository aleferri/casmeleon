package language

//Set of a symbol
type Set struct {
	name    string
	valueOf func(string) int32
}

//Contains the specified symbol
func (s *Set) Contains(n string) bool {
	return s.valueOf(n) > -1
}

//Value of the specified symbol
func (s *Set) Value(n string) uint32 {
	return uint32(s.valueOf(n))
}
