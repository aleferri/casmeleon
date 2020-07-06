package lexing

import (
	"bitbucket.org/mrpink95/casmeleon/internal/text"
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"strings"
)

//SourceReader read the source file/string and split line into tokens (raw)
type SourceReader struct {
	reader     *bufio.Reader
	closeOp    func() bool
	separators string
	sourceName string
	endReached bool
	lineNumber uint
}

//SourceName return the source name
func (source *SourceReader) SourceName() string {
	return source.sourceName
}

//HasNext return true if the end of file is reached
func (source *SourceReader) HasNext() bool {
	return !source.endReached
}

//NextLine read the line and then split in tokens (raw),
//return a list
func (source *SourceReader) NextLine() *list.List {
	var line, err = source.reader.ReadString('\n')
	source.lineNumber++
	if err != nil {
		source.endReached = true
		if err == io.EOF {
			source.closeOp()
			return SplitKeepSeparators(line, source.separators)
		}
		return list.New()
	}
	trimLine := strings.TrimSuffix(line, "\n")
	return SplitKeepSeparators(strings.TrimSuffix(trimLine, "\r"), source.separators)
}

//LineNumber return, yes your guess is right: the line number, but go lint complain if i don't write it
func (source *SourceReader) LineNumber() uint {
	return source.lineNumber
}

//NormalizeLine process the raw tokens in line regrouping quoted strings and multi-char operators
func (source *SourceReader) NormalizeLine(line *list.List, lineNumber uint, symbols string) *text.SourceLine {
	return text.NewSourceLine(RegroupSymbols(JoinQuote(line), symbols), lineNumber, source.sourceName)
}

//FilterLine remove single line comments and filter whitespace, return a slice of tokens
func (source *SourceReader) FilterLine(line *text.SourceLine, options TokenMatchingOptions) []text.Token {
	var tokens = options.LineToTokens(line, source.lineNumber, source.sourceName)
	return tokens
}

//NewSourceReaderFromFile create a SourceReader from a file and a separators string
func NewSourceReaderFromFile(filePath, separators string) *SourceReader {
	var file, err = os.Open(filePath)
	if err != nil {
		fmt.Printf("Error during opening of file %v\n", filePath)
		return NewSourceReader(strings.NewReader(""), separators, "Error during file opening")
	}
	closeOp := func() bool {
		err := file.Close()
		return err != nil
	}
	return &SourceReader{bufio.NewReader(file), closeOp, separators, filePath, false, 0}
}

// NewSourceReader create a SourceReader from a reader and a separator string
func NewSourceReader(reader io.Reader, separators, sourceName string) *SourceReader {
	return &SourceReader{bufio.NewReader(reader), func() bool { return true }, separators, sourceName, false, 0}
}
