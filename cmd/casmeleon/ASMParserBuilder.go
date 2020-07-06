package main

import (
	"bitbucket.org/mrpink95/casmeleon/internal/langdef"
	"bitbucket.org/mrpink95/casmeleon/internal/parsing"
	"bitbucket.org/mrpink95/casmeleon/internal/text"
	"bitbucket.org/mrpink95/casmeleon/internal/ui"
	"strings"
)

//Tag values, respect Ignore = 0 and Append = 1 default declared in parsing
const (
	IgnoreAll = 0 + iota
	OnlyAppend
	LabelTag
	IncludeTag
	OpcodeTag
)

func matchLabel() parsing.MatchRule {
	noErr := parsing.NewIncompleteError(parsing.NoError, "")
	colonErr := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected ':'")
	matchName := parsing.MatchToken(noErr, text.Identifier)
	matchColon := parsing.MatchSkipString(colonErr, ":")
	return parsing.MatchAll(LabelTag, matchName, matchColon)
}

func isASMLabel(tokens []text.Token, pState parsing.ParserState) ([]text.Token, bool) {
	for len(tokens) < 2 {
		pState.Pull(&tokens)
	}
	if strings.EqualFold(tokens[1].Value(), ":") {
		return tokens, true
	}
	return tokens, false
}

func matchInclude() parsing.MatchRule {
	noErr := parsing.NewIncompleteError(parsing.NoError, "")
	expectedQuote := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected double quote string")
	expectedEOL := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected end of line after include statement")
	matchKeyword := parsing.MatchToken(noErr, text.KeywordInclude)
	matchPath := parsing.MatchToken(expectedQuote, text.DoubleQuotedString)
	matchEOL := parsing.MatchAnyToken(expectedEOL, text.EOL, text.EOF)
	return parsing.MatchAll(IncludeTag, matchKeyword, matchPath, matchEOL)
}

func readToken(tokens *[]text.Token, pState parsing.ParserState) (text.Token, bool) {
	if len(*tokens) < 1 {
		pState.Pull(tokens)
	}
	token := (*tokens)[0]
	*tokens = (*tokens)[1:]
	return token, token.EnumType() == text.EOF
}

func skipEOL(tokens []text.Token, pState parsing.ParserState) ([]text.Token, bool) {
	for len(tokens) < 1 {
		pState.Pull(&tokens)
		if tokens[0].EnumType() == text.EOF {
			return tokens, false
		} else if tokens[0].EnumType() == text.EOL {
			tokens = tokens[1:]
		}
	}
	return tokens, true
}

func matchOpcode(lang *langdef.LangDef) parsing.MatchRule {
	return func(tokens []text.Token, pState parsing.ParserState) ([]text.Token, bool) {
		wnd := lang.NewFilterWindow()
		name, eof := readToken(&tokens, pState)
		if eof {
			pState.ReportError(name, parsing.ErrorExpectedToken, ", expected opcode name, not eof")
			return tokens, false
		}
		wnd = wnd.FilterByName(name.Value())
		args := []text.Token{}
		arg, eof := readToken(&tokens, pState)
		var keepArg bool
		for !eof && arg.EnumType() != text.EOL {
			wnd, keepArg = wnd.FilterByToken(arg)
			if wnd.IsEmpty() {
				pState.ReportError(arg, parsing.InvalidOpcodeArgument, ", invalid opcode arguments: no match found")
				return tokens, false
			}
			if keepArg {
				args = append(args, arg)
			}
			arg, eof = readToken(&tokens, pState)
		}
		_, opcodeIndex := wnd.Collect()
		pState.Push([]text.Token{name}, OnlyAppend)
		pState.Push(args, OpcodeTag)
		pState.Push([]text.Token{}, (opcodeIndex+1)<<4)
		return tokens, true
	}
}

//ASMParserState is parsed state for the ASM part
type ASMParserState struct {
	pUI            ui.UI
	src            parsing.TokenSource
	lastParsed     []text.Token
	lastOpcode     int
	customIdentify func(t text.Token) (text.Token, ui.SourceCodeError) //custom post token conversion identification
}

//NewParserState return a new ASMParserState from source
func NewParserState(pUI ui.UI, src parsing.TokenSource) *ASMParserState {
	return &ASMParserState{pUI, src, []text.Token{}, -1, nil}
}

//SetCustomIdentification set a custom TokenType identification function for every token
func (p *ASMParserState) SetCustomIdentification(identify func(t text.Token) (text.Token, ui.SourceCodeError)) {
	p.customIdentify = identify
}

//Push the tokens in the parsed collection
func (p *ASMParserState) Push(tokens []text.Token, tag int) {
	if tag == IgnoreAll {
		return
	}
	if tag >= 16 {
		p.lastOpcode = (tag - 1) >> 4
	}
	p.lastParsed = append(p.lastParsed, tokens...)
}

//Pull the tokens from the source
func (p *ASMParserState) Pull(tokens *[]text.Token) int {
	lex := p.src.NextToken(p.pUI)
	identified, err := p.customIdentify(lex.WithType(parsing.FindTokenType(lex)))
	if err != nil {
		err.Report(p.pUI, p.src.Lines()[err.GetLine()])
		*tokens = append(*tokens, lex)
		return 1
	}
	*tokens = append(*tokens, identified)
	if lex.EnumType() == text.EOF {
		return 0
	}
	return 1
}

//Mark the current parsing state
func (p *ASMParserState) Mark() (int, int) {
	return 0, 0
}

//ToCST does nothing
func (p *ASMParserState) ToCST(markT int, markL int, tag int) {
}

//ReportError on a token
func (p *ASMParserState) ReportError(wrong text.Token, err parsing.ParsingErrorType, msg string) {
	if err == parsing.NoError {
		return
	}
	err.ReportToken(p.pUI, wrong, p.src, msg)
}
