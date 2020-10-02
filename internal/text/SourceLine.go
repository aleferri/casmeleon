package text

import (
	"github.com/aleferri/casmeleon/internal/ui"
	"unicode/utf8"
)

//SourceLine keep original strings with whitespaces, line number and source name
type SourceLine struct {
	words      []string
	number     uint
	sourceName string
}

//NewSourceLine create a new source line
func NewSourceLine(words []string, lineNumber uint, sourceName string) *SourceLine {
	return &SourceLine{words, lineNumber, sourceName}
}

//LineNumber return the line number
func (line *SourceLine) LineNumber() uint {
	return line.number
}

//StringAt return the string at i index
func (line *SourceLine) StringAt(i uint) string {
	return line.words[i]
}

//SourceName return the source in which this line is found
func (line *SourceLine) SourceName() string {
	return line.sourceName
}

//AllStrings return a slice with all strings in the line
func (line *SourceLine) AllStrings() []string {
	return line.words
}

//Print print the line
func (line *SourceLine) Print(ui ui.UI) {
	for _, elem := range line.words {
		ui.ReportMessage(elem, false)
	}
	ui.ReportMessage("", true)
}

//RuneIndex calculate the index of the first rune of the indexOfWord argument
func (line *SourceLine) RuneIndex(ui ui.UI, indexOfWord uint) uint {
	var result = uint(0)
	for i := uint(0); i < indexOfWord; i++ {
		result += uint(utf8.RuneCountInString(line.words[i]))
	}
	return result
}
