package langdef

import (
	"github.com/aleferri/casmeleon/internal/parsing"
	"github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/internal/ui"
)

//LangParserState is a configurable ParseState that can skip the EOL
//if needed
type LangParserState struct {
	pUI        ui.UI
	src        parsing.TokenSource
	lastParsed []text.Token
	cst        []CSTNode
	skipEOL    bool
	report     func(err *parsing.ParsingError)
}

//NewParserState return a new parser state
func NewParserState(pUI ui.UI, src parsing.TokenSource, skipEOL bool, report func(err *parsing.ParsingError)) *LangParserState {
	return &LangParserState{pUI, src, []text.Token{}, []CSTNode{}, skipEOL, report}
}

//Push the tokens in the parsed collection
func (s *LangParserState) Push(tokens []text.Token, tag int) {
	if tag == parsing.DefaultIgnore {
		return
	}
	s.lastParsed = append(s.lastParsed, tokens...)
}

//Pull the tokens from the source
func (s *LangParserState) Pull(tokens *[]text.Token) int {
	lex := s.src.NextToken(s.pUI)
	eol := lex.EnumType() == text.EOL && s.skipEOL
	for eol {
		lex = s.src.NextToken(s.pUI)
		eol = lex.EnumType() == text.EOL && s.skipEOL
	}
	lex = identifyOperators(lex.WithType(parsing.FindTokenType(lex)))
	*tokens = append(*tokens, lex)
	if lex.EnumType() == text.EOF {
		return 0
	}
	return 1
}

//Mark current state
func (s *LangParserState) Mark() (int, int) {
	return len(s.cst), len(s.lastParsed)
}

//ToCST reduce all trees and parsed tokens after the mark to a single tree
func (s *LangParserState) ToCST(markT int, markL int, tag int) {
	if tag == OnlyAppend || tag == Ignore {
		return
	}
	cst := NewRootCSTNode(tag, 0)
	cst.content = append(cst.content, s.lastParsed[markL:]...)
	s.lastParsed = s.lastParsed[:markL]
	for i := markT; i < len(s.cst); i++ {
		child := s.cst[i]
		child.parent = &cst
		cst.AddChild(child)
	}
	if len(s.cst) > markT {
		s.cst = s.cst[:markT+1]
		s.cst[markT] = cst
	} else {
		s.cst = append(s.cst, cst)
	}
}

//SetSkipEOL ask to the ParseState to filter (or not) EOL tokens. The filter (if skipEOL is true)
//will be applied to the buffered tokens too
func (s *LangParserState) SetSkipEOL(skipEOL bool, tokens *[]text.Token) {
	s.skipEOL = skipEOL
	if skipEOL {
		lastEmpty := 0
		for i, t := range *tokens {
			if t.EnumType() != text.EOL {
				if lastEmpty > -1 && i != lastEmpty {
					(*tokens)[lastEmpty] = (*tokens)[i]
				}
				lastEmpty++
			}
		}
		*tokens = (*tokens)[:lastEmpty]
	}
}

//ReportError on a token
func (s *LangParserState) ReportError(wrong text.Token, err parsing.ParsingErrorType, msg string) {
	if err == parsing.NoError {
		return
	}
	err.ReportToken(s.pUI, wrong, s.src, msg)
}
