package main

import (
	"bufio"
	"bytes"
	"strings"
	"unicode"

	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/scanner"
	"github.com/aleferri/casmeleon/pkg/text"
)

// numberPrefixes holds the prefixes declared by .num in the language file.
// IdentifySymbol needs them to classify literals that do not start with a digit.
var numberPrefixes []string

// SetNumberPrefixes must be called after the language file has been parsed
func SetNumberPrefixes(prefixes []string) {
	numberPrefixes = prefixes
}

func hasDeclaredPrefix(str string) bool {
	for _, p := range numberPrefixes {
		if len(str) > len(p) && strings.HasPrefix(str, p) {
			return true
		}
	}
	return false
}

// IdentifySymbol in the assembly source. Unlike the homonym in internal/casm, this
// one must consult the grammar declared by .num in the language file: the .casm has
// already been parsed when this runs, so prefixes are known.
func IdentifySymbol(s []rune, fileOffset uint32, count uint32) text.Symbol {
	if unicode.IsDigit(s[0]) {
		return text.SymbolOf(fileOffset, count, string(s), text.Number)
	}
	if hasDeclaredPrefix(string(s)) {
		return text.SymbolOf(fileOffset, count, string(s), text.Number)
	}
	if s[0] == '\n' || s[0] == '\r' {
		return text.SymbolOf(fileOffset, count, string(s), text.EOL)
	}
	if unicode.IsSpace(s[0]) {
		return text.SymbolOf(fileOffset, count, string(s), text.WHITESPACE)
	}
	if s[0] == '"' {
		return text.SymbolOf(fileOffset, count, string(s), text.QuotedString)
	}
	if s[0] == '\'' {
		return text.SymbolOf(fileOffset, count, string(s), text.QuotedChar)
	}
	if s[0] == ';' {
		return text.SymbolOf(fileOffset, count, string(s), text.WHITESPACE)
	}
	str := string(s)
	id, ok := identifyMap[str]
	if !ok {
		return text.SymbolOf(fileOffset, count, str, text.Identifier)
	}
	return text.SymbolOf(fileOffset, count, str, id)
}

// AssemblyStream is the stream of the symbols
type AssemblyStream struct {
	parent *AssemblyStream
	source *bufio.Reader
	buffer []text.Symbol
	repo   *text.Source
}

// MakeRootStream for the parser
func MakeRootStream(source *bufio.Reader, repo *text.Source) parser.Stream {
	return &AssemblyStream{parent: nil, source: source, repo: repo, buffer: []text.Symbol{}}
}

func MakeChildStream(source *bufio.Reader, repo *text.Source, parent *AssemblyStream) parser.Stream {
	return &AssemblyStream{parent: parent, source: source, repo: repo, buffer: []text.Symbol{}}
}

// Buffer ensure that the buffer of the stream contains at least 1 element
func (s *AssemblyStream) Buffer() {
	for len(s.buffer) == 0 {
		line, ioErr := s.source.ReadBytes('\n')

		if ioErr != nil {
			s.buffer = append(s.buffer, text.SymbolEmpty(s.repo.FileIndex()).WithID(text.EOF))
			return
		}

		runes := bytes.Runes(line)

		temps, _ := scanner.FastScan(runes, true, scanFollowMap)

		scanner.ClassifyBasicASMTokens(temps)
		scanned := scanner.MergeASMLine(temps)

		for _, t := range scanned {
			sym := IdentifySymbol(t.Runes(), s.repo.FileIndex(), s.repo.Count())
			s.repo.Append(sym)
			if sym.ID() != text.WHITESPACE {
				s.buffer = append(s.buffer, sym)
			}
		}
	}
}

// Next symbol in the internal buffer
func (s *AssemblyStream) Next() text.Symbol {
	s.Buffer()

	result := s.buffer[0]
	s.buffer = s.buffer[1:]
	return result
}

// Peek the symbol from the internal buffer
func (s *AssemblyStream) Peek() text.Symbol {
	s.Buffer()

	return s.buffer[0]
}

// Source of the stream
func (s *AssemblyStream) Source() *text.Source {
	return s.repo
}
