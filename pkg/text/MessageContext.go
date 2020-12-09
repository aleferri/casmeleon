package text

//MessageContext is an error location inside the source code
type MessageContext struct {
	fileOffset uint32
	symOffset  uint32
	scopeLeft  string
	scopeRight string
}

//FileName of the error
func (e *MessageContext) FileName(index []string) string {
	return index[e.fileOffset]
}

//Position of the value
func (e *MessageContext) Position() uint32 {
	return e.symOffset
}

//MakeMessageContext create a context for the message
func MakeMessageContext(sym Symbol, scopeLeft string, scopeRight string) MessageContext {
	return MessageContext{fileOffset: sym.fileOffset, symOffset: sym.symOffset, scopeLeft: scopeLeft, scopeRight: scopeRight}
}
