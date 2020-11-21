package ui

//UI is an interface to report errors, warnings and messages to user
type UI interface {
	ReportSourceError(msg string, index uint, line Printable)
	ReportSourceWarning(msg string, index uint, line Printable)
	ReportError(msg string, newLine bool)
	ReportWarning(msg string, newLine bool)
	ReportMessage(msg string, newLine bool)
	GetErrorCount() int
}
