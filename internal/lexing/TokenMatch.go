package lexing

import (
	"bitbucket.org/mrpink95/casmeleon/internal/text"
	"strings"
)

//TokenMatchingOptions contains definition of delimiters and operators for the lexer
type TokenMatchingOptions struct {
	operators, separators string
	multiCharSymbols      string
	lineComment           string
}

func (opt TokenMatchingOptions) GetAllSeparators() string {
	return opt.separators + opt.operators + opt.lineComment + " \t "
}

func (opt TokenMatchingOptions) GetOperators() string {
	return opt.operators
}

func (opt TokenMatchingOptions) GetMultiCharSymbols() string {
	return opt.multiCharSymbols
}

func NewMatchingOptions(operators, separators, multiCharSymbols, lineCommentChar string) TokenMatchingOptions {
	return TokenMatchingOptions{operators, separators, multiCharSymbols, lineCommentChar}
}

//Convert the line to a slice of tokens, also skip whitespaces and comments
//search why fileName was used
func (opt TokenMatchingOptions) LineToTokens(l *text.SourceLine, lineNumber uint, fileName string) []text.Token {
	var tokens = make([]text.Token, 0)
	for i := uint(0); i < uint(len(l.AllStrings())); i++ {
		var value = l.StringAt(i)
		if strings.EqualFold(value, opt.lineComment) { // line comment found, discard remaining part of the line
			return tokens
		}
		if !strings.EqualFold(strings.TrimSpace(value), "") { // skip space
			var tType = FindRawTokenType(value, opt.operators, opt.separators)
			tokens = append(tokens, text.NewToken(value, tType, uint(i), l))
		}
	}
	return tokens
}

func FindRawTokenType(value, operators, separators string) text.TokenType {
	if strings.HasPrefix(value, "'") {
		return text.SingleQuotedString
	} else if strings.HasPrefix(value, "\"") {
		return text.DoubleQuotedString
	} else if strings.Contains(operators, value) {
		return text.BinaryOperator //generally assumed binary
	} else if strings.Contains(separators, value) {
		return text.GenericSeparator //generic separator
	} else if strings.IndexAny(value, "0123456789") == 0 {
		return text.Number
	}
	return text.Identifier
}
