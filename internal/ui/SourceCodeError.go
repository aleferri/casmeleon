package ui

//SourceCodeError is an enrichment of the Error interface in go
type SourceCodeError interface {
	Error() string
	Report(ui UI, line Printable)
	GetLine() uint
}
