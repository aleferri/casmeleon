package text

//Error is an error reportable to an user
type Error struct {
	fileOffset uint32
	symOffset  uint32
	message    string
}

//FileName of the error
func (e *Error) FileName(index []string) string {
	return index[e.fileOffset]
}

//Position of the value
func (e *Error) Position() uint32 {
	return e.symOffset
}

//Message of the error
func (e *Error) Message() string {
	return e.message
}
