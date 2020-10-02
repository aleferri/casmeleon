package parsing

import (
	"github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/internal/ui"
)

//ParsingErrorType is an enum representing types of errors
type ParsingErrorType int

//ParsingErrorType values
const (
	ErrorUnresolvedSymbol = 0 + iota
	ErrorDuplicatedLabel
	ErrorExpectedToken
	InvalidNumberFormat
	InvalidOpcodeArgument
	InvalidOpcodeFormat
	ErrorCyclicInclusion
	ErrorExpectedEOL
	ErrorIncompleteExpression
	ErrorUserDefined
	ErrorInternal
	NoError
)

//String implements Stringer interface for ParsingErrorType
func (eType ParsingErrorType) String() string {
	var errorTypeString = []string{
		"unresolved symbol", "duplicate label", "unexpected token",
		"invalid number format", "invalid opcode argument", "invalid format for opcode", "cyclic inclusion",
		"end of line was expected, found", "incomplete expression", "", "internal error", "no error",
	}
	if eType <= ErrorInternal {
		return errorTypeString[eType]
	}
	return ""
}

//ReportToken report in a token
func (eType ParsingErrorType) ReportToken(ui ui.UI, token text.Token, b TokenSource, arg string) {
	var line, pos, _ = token.Position()
	if line < 1 {
		ui.ReportError("Unexpected begin of file", true)
	} else {
		ui.ReportSourceError(eType.String()+" '"+token.Value()+"' "+arg+" in file "+b.SourceName(), pos, b.Lines()[line-1])
	}
}

//ParsingError contains all type of error that can occur during parsing
type ParsingError struct {
	errorType      ParsingErrorType
	arg            string
	msg            string
	position, line uint
	fileName       string
}

//NewParsingError create new parsing error
func NewParsingError(t text.Token, eType ParsingErrorType, err string) ParsingError {
	line, pos, file := t.Position()
	return ParsingError{eType, t.Value(), err, pos, line, file}
}

/*
	This section implement AssemblerError interface for the type ParsingError
*/

//Error message
func (err ParsingError) Error() string {
	return err.errorType.String() + " " + err.arg
}

//Report the error
func (err ParsingError) Report(ui ui.UI, line ui.UIPrintable) {
	ui.ReportSourceError(err.errorType.String()+" '"+err.arg+"' "+err.msg+" in file "+err.fileName, err.position, line)
}

//GetLine return the line of the error
func (err ParsingError) GetLine() uint {
	return err.line
}
