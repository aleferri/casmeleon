package casm

import (
	"bufio"
	"bytes"
	"unicode"

	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/scanner"
	"github.com/aleferri/casmeleon/pkg/text"
)

//IdentifySymbol in the casm language
func IdentifySymbol(s []rune, fileOffset uint32, count uint32) text.Symbol {
	if unicode.IsDigit(s[0]) {
		return text.SymbolOf(fileOffset, count, string(s), text.Number)
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
	if len(s) > 1 && s[0] == '/' && (s[1] == '/' || s[1] == '*') {
		return text.SymbolOf(fileOffset, count, string(s), text.WHITESPACE)
	}
	str := string(s)
	id, ok := identifyMap[str]
	if !ok {
		return text.SymbolOf(fileOffset, count, str, text.Identifier)
	}
	return text.SymbolOf(fileOffset, count, str, id)
}

//SymbolsStream is the stream of the symbols
type SymbolsStream struct {
	source     *bufio.Reader
	buffer     []text.Symbol
	repo       *text.Source
	left       []rune
	lastMerged *scanner.Token
	lastID     int32
}

//BuildStream for the parser
func BuildStream(source *bufio.Reader, repo *text.Source) parser.Stream {
	return &SymbolsStream{source: source, repo: repo, left: []rune{}, buffer: []text.Symbol{}, lastMerged: nil, lastID: -1}
}

//Buffer ensure that the buffer of the stream contains at least 1 element
func (s *SymbolsStream) Buffer() {
	for len(s.buffer) == 0 {
		line, ioErr := s.source.ReadBytes('\n')

		runes := s.left
		if len(runes) == 0 {
			runes = bytes.Runes(line)
		} else {
			runes = append(runes, bytes.Runes(line)...)
		}

		temps, left := scanner.FastScan(runes, ioErr != nil, scanFollowMap)
		s.left = left

		scanner.ClassifyMergeableTokens(temps)
		scanned, nextMerged, nextID := scanner.Merge(temporaryTokenMarks, temps, s.lastMerged, s.lastID)
		s.lastMerged = nextMerged
		s.lastID = nextID

		for _, t := range scanned {
			sym := IdentifySymbol(t.Runes(), s.repo.FileIndex(), s.repo.Count())
			s.repo.Append(sym)
			if sym.ID() != text.WHITESPACE {
				s.buffer = append(s.buffer, sym)
			}
		}

		if ioErr != nil {
			s.buffer = append(s.buffer, text.SymbolEmpty(s.repo.FileIndex()).WithID(text.EOF))
		}
	}
}

//Next symbol in the internal buffer
func (s *SymbolsStream) Next() text.Symbol {
	s.Buffer()

	result := s.buffer[0]
	s.buffer = s.buffer[1:]
	return result
}

//Peek the symbol from the internal buffer
func (s *SymbolsStream) Peek() text.Symbol {
	s.Buffer()

	return s.buffer[0]
}

//Source of the stream
func (s *SymbolsStream) Source() *text.Source {
	return s.repo
}
