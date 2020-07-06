package ui

import "fmt"

//UIPrintable is an interface of a user printable object
type UIPrintable interface {
	Print(ui UI)                            //Print the object using ui
	RuneIndex(ui UI, indexOfWord uint) uint //RuneIndex calculate the index of the first rune of the specified word
	LineNumber() uint                       //LineNumber return the lineNumber
	StringAt(i uint) string                 //StringAt return the i-th string
}

//UI is an interface to report errors, warnings and messages to user
type UI interface {
	ReportSourceError(msg string, index uint, line UIPrintable)
	ReportSourceWarning(msg string, index uint, line UIPrintable)
	ReportError(msg string, newLine bool)
	ReportWarning(msg string, newLine bool)
	ReportMessage(msg string, newLine bool)
	GetErrorCount() int
}

//ConsoleUI is a concrete implementation of ErrorHandler
type ConsoleUI struct {
	warningAsErrors, suppressWarnings bool
	errorCount                        int
}

//NewConsoleUI return a new ConsoleUI
func NewConsoleUI(warningAsErrors, suppressWarnings bool) *ConsoleUI {
	return &ConsoleUI{warningAsErrors, suppressWarnings, 0}
}

//ReportSourceError report an error to the user
func (c *ConsoleUI) ReportSourceError(msg string, index uint, line UIPrintable) {
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
func (c *ConsoleUI) ReportSourceWarning(msg string, index uint, line UIPrintable) {
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
func (c *ConsoleUI) ReportError(msg string, newLine bool) {
	c.errorCount++
	c.ReportMessage("Error: "+msg, newLine)
}

//ReportWarning report a warning without format, can be ignored
func (c *ConsoleUI) ReportWarning(msg string, newLine bool) {
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
func (c *ConsoleUI) ReportMessage(msg string, newLine bool) {
	fmt.Printf("%v", msg)
	if newLine {
		fmt.Print("\n")
	}
}

//GetErrorCount return the number of errors found during parsing
func (c *ConsoleUI) GetErrorCount() int {
	return c.errorCount
}
