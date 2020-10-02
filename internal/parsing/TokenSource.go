package parsing

import (
	"github.com/aleferri/casmeleon/internal/lexing"
	"github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/internal/ui"
	"io"
)

//FindTokenType to a text.Token
func FindTokenType(t text.Token) text.TokenType {
	improvedTokenType, found := text.TokenTypeMap[t.Value()]
	if found {
		return improvedTokenType
	}
	return t.EnumType()
}

//TokenSource is the input of the parser
type TokenSource interface {
	SourceName() string                //SourceName return the name of the source (may be a location)
	Lines() []*text.SourceLine         //Lines return all saved lines
	NextToken(ui ui.UI) text.Token     //NextToken return the next token in the source
	SkipEndOfLine(ui ui.UI) text.Token //SkipEndOfLine skip the next empty lines
}

//TokenBuffer keep a one line buffer of the tokens
type TokenBuffer struct {
	source    *lexing.SourceReader
	options   lexing.TokenMatchingOptions
	lines     []*text.SourceLine
	syncLines func([]*text.SourceLine)
	lastLine  []text.Token
	index     int
}

//NewTokenBufferFromFile return a new TokenBuffer from a file
func NewTokenBufferFromFile(file string, options lexing.TokenMatchingOptions) *TokenBuffer {
	var buffer = TokenBuffer{}
	buffer.options = options
	buffer.source = lexing.NewSourceReaderFromFile(file, options.GetAllSeparators())
	buffer.index = 0
	return &buffer
}

//NewTokenBuffer return a new TokenBuffer from any reader
func NewTokenBuffer(name string, text io.Reader, options lexing.TokenMatchingOptions) *TokenBuffer {
	var buffer = TokenBuffer{}
	buffer.options = options
	buffer.source = lexing.NewSourceReader(text, options.GetAllSeparators(), name)
	buffer.index = 0
	return &buffer
}

//SyncLines return
func (buffer *TokenBuffer) SyncLines(syncLines func([]*text.SourceLine)) {
	buffer.syncLines = syncLines
}

//SourceName return the source name of the tokens
func (buffer *TokenBuffer) SourceName() string {
	return buffer.source.SourceName()
}

//Lines return all saved lines
func (buffer *TokenBuffer) Lines() []*text.SourceLine {
	return buffer.lines
}

func (buffer *TokenBuffer) saveNotNil(line *text.SourceLine) {
	buffer.lines = append(buffer.lines, line)
	if buffer.syncLines != nil {
		buffer.syncLines(buffer.lines)
	}
}

//NextToken return the next token in the source file
//empty lines are skipped
func (buffer *TokenBuffer) NextToken(ui ui.UI) text.Token {
	isInit := buffer.lastLine == nil
	if buffer.index < len(buffer.lastLine) {
		buffer.index++
		return buffer.lastLine[buffer.index-1]
	}
	var lastToken text.Token
	if buffer.index > 0 {
		lastToken = buffer.lastLine[buffer.index-1]
	} else {
		lastToken = text.NewToken("EOF", text.EOF, 0, text.NewSourceLine([]string{""}, uint(len(buffer.Lines())), buffer.source.SourceName()))
	}
	for buffer.index >= len(buffer.lastLine) {
		if !buffer.source.HasNext() {
			return lastToken.WithType(text.EOF).WithValue("EOF")
		}
		buffer.index = 0
		var lineWords = buffer.source.NextLine()
		var lineNumber = buffer.source.LineNumber()
		var line = buffer.source.NormalizeLine(lineWords, lineNumber, buffer.options.GetMultiCharSymbols())
		buffer.saveNotNil(line)
		var tokens = buffer.source.FilterLine(line, buffer.options)
		buffer.lastLine = tokens
	}
	if isInit {
		return buffer.NextToken(ui)
	}
	return lastToken.WithType(text.EOL).WithValue("EOL")
}

//SkipEndOfLine return the first token that is not EOL
func (buffer *TokenBuffer) SkipEndOfLine(ui ui.UI) text.Token {
	var token = buffer.NextToken(ui)
	for token.EnumType() == text.EOL {
		token = buffer.NextToken(ui)
	}
	return token
}
