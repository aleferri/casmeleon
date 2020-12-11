package parser

import "github.com/aleferri/casmeleon/pkg/text"

//ErrorReport contains the necessary info for error reporting
type ErrorReport struct {
	ref     text.MessageContext
	msg     string
	errCode int
}

//MakeErrorReport from base info
func MakeErrorReport(ref text.MessageContext, msg string, errCode int) *ErrorReport {
	return &ErrorReport{ref, msg, errCode}
}

func (e *ErrorReport) Error() string {
	return e.msg
}

func (e *ErrorReport) PrettyPrint() {

}
