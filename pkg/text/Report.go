package text

//Report is the full report of the message regarding the original source
type Report interface {
	Source(span []Symbol) Report
	Context(span []Symbol) Report
	Message(msg string) Report
}
