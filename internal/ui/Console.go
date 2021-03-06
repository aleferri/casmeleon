package ui

import "fmt"

//Console is a concrete implementation of ErrorHandler
type Console struct {
	warningAsErrors, suppressWarnings bool
	errorCount                        int
}

//NewConsole return a new Console
func NewConsole(warningAsErrors, suppressWarnings bool) *Console {
	return &Console{warningAsErrors, suppressWarnings, 0}
}

//ReportSourceError report an error to the user
func (c *Console) ReportSourceError(msg string, index uint, line Printable) {
	c.errorCount++
	fmt.Printf("Error at line %d: %v\n", line.LineNumber(), msg)
	fmt.Print("Line: ")
	line.Print(c)
	var charIndex = line.RuneIndex(c, index) + uint(len("Line: "))
	var wordSize = uint(len(line.StringAt(index)))
	for i := uint(0); i < charIndex+wordSize; i++ {
		if i >= charIndex {
			fmt.Print("^")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	fmt.Println()
}

//ReportSourceWarning report a warning to the user if suppressWarning is false
func (c *Console) ReportSourceWarning(msg string, index uint, line Printable) {
	if c.suppressWarnings {
		return
	}
	if c.warningAsErrors {
		c.ReportSourceError(msg, index, line)
		return
	}
	fmt.Printf("\nWarning at line %d: %s\n", line.LineNumber(), msg)
	fmt.Print("Line: ")
	line.Print(c)
	var charIndex = line.RuneIndex(c, index) + uint(len("Line: "))
	var wordSize = uint(len(line.StringAt(index)))
	for i := uint(0); i < charIndex+wordSize; i++ {
		if i >= charIndex {
			fmt.Print("^")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	fmt.Println()
}

//ReportError report a generic error without format
func (c *Console) ReportError(msg string, newLine bool) {
	c.errorCount++
	c.ReportMessage("Error: "+msg, newLine)
}

//ReportWarning report a warning without format, can be ignored
func (c *Console) ReportWarning(msg string, newLine bool) {
	if c.suppressWarnings {
		return
	}
	if c.warningAsErrors {
		c.ReportError(msg, newLine)
	} else {
		c.ReportMessage("Warning: "+msg, newLine)
	}
}

//ReportMessage report a message to the user
func (c *Console) ReportMessage(msg string, newLine bool) {
	fmt.Printf("%v", msg)
	if newLine {
		fmt.Print("\n")
	}
}

//GetErrorCount return the number of errors found during parsing
func (c *Console) GetErrorCount() int {
	return c.errorCount
}
