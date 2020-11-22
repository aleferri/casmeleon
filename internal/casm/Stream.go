package language

import (
	"bufio"
	"bytes"

	"github.com/aleferri/casmeleon/pkg/lexer"
	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

//SymbolsStream is the stream of the symbols
type SymbolsStream struct {
	source bufio.Reader
	lexer  lexer.Lexer
	report func(error, uint32, uint32)
	buffer []text.Symbol
	left   []rune
	repo   text.Source
}

//BuildStream for the parser
func BuildStream(source bufio.Reader, lexer lexer.Lexer, repo text.Source, report func(error, uint32, uint32)) parser.Stream {
	return &SymbolsStream{source: source, lexer: lexer, repo: repo, report: report, left: []rune{}, buffer: []text.Symbol{}}
}

//Buffer ensure that the buffer of the stream contains at least 1 element
func (s *SymbolsStream) Buffer() {
	for len(s.buffer) == 0 {
		line, ioErr := s.source.ReadBytes('\n')

		runes := s.left
		runes = append(runes, bytes.Runes(line)...)
		lexed, err := s.lexer.Scan(0, 0, runes, ioErr != nil)

		s.left = lexed.Unlexed()
		s.buffer = s.repo.AppendAll(lexed.Symbols())

		if err != nil {
			s.report(err, s.repo.FileIndex(), s.repo.Count())
			s.buffer = append(s.buffer, text.SymbolEmpty(s.repo.FileIndex()))
		}

		if ioErr != nil {
			s.buffer = append(s.buffer, text.SymbolEmpty(s.repo.FileIndex()))
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
func (s *SymbolsStream) Source() text.Source {
	return s.repo
}
