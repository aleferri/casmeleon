package exec

//ExpressionList contains a list of Expression
type ExpressionList struct {
	list []Expression
}

//Resolve the opcode output
func (elist *ExpressionList) Resolve(params []uint32) []uint8 {
	result := []uint8{}
	for _, e := range elist.list {
		result = append(result, uint8(e.Resolve(params)))
	}
	return result
}
