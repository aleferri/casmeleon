package parsing

import (
	"github.com/aleferri/casmeleon/internal/text"
	"strings"
)

//MatchRule match a rule in the token stream and push the matched token to the helper result
//if not enough token are provided, pull function is called of the pHelper is called
//return tokens not parsed, keep parsing
type MatchRule func(tokens []text.Token, pState ParserState) ([]text.Token, bool)

//IncompleteError contains the token related part to throw a parsing error
type IncompleteError struct {
	errType ParsingErrorType
	errMsg  string
}

//NewIncompleteError return a new incomplete error type with type and message
func NewIncompleteError(errType ParsingErrorType, errMsg string) IncompleteError {
	return IncompleteError{errType: errType, errMsg: errMsg}
}

//MatchToken return a MatchRule that match a specified TokenType
//return false with the same MatchRule condition and if an error is found
func MatchToken(err IncompleteError, tokenType text.TokenType) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		if len(tokens) < 1 {
			pState.Pull(&tokens)
		}
		notEof := tokens[0].EnumType() != text.EOF
		match := tokens[0].EnumType() == tokenType
		if !match {
			pState.ReportError(tokens[0], err.errType, err.errMsg)
			return tokens, notEof && err.errType == NoError
		}
		pState.Push(tokens[0:1], DefaultAppend)
		return tokens[1:], notEof
	}
}

//MatchSkipToken return a MatchRule that match a specified TokenType
//it push the token with the 'IgnoreTokens' flag
//return false with the same MatchRule condition and if an error is found
func MatchSkipToken(err IncompleteError, tokenType text.TokenType) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		if len(tokens) < 1 {
			pState.Pull(&tokens)
		}
		notEof := tokens[0].EnumType() != text.EOF
		match := tokens[0].EnumType() == tokenType
		if !match {
			pState.ReportError(tokens[0], err.errType, err.errMsg)
			return tokens, notEof && err.errType == NoError
		}
		pState.Push(tokens[0:1], DefaultIgnore)
		return tokens[1:], notEof
	}
}

//MatchSkipString return a MatchRule that match a specified TokenType
//it push the token with the 'IgnoreTokens' flag
//return false with the same MatchRule condition and if an error is found
func MatchSkipString(err IncompleteError, s string) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		if len(tokens) < 1 {
			pState.Pull(&tokens)
		}
		notEof := tokens[0].EnumType() != text.EOF
		match := strings.EqualFold(tokens[0].Value(), s)
		if !match {
			pState.ReportError(tokens[0], err.errType, err.errMsg)
			return tokens, notEof && err.errType == NoError
		}
		pState.Push(tokens[0:1], DefaultIgnore)
		return tokens[1:], notEof
	}
}

//MatchNotAnyToken return a MatchRule that match every token except the specified TokenType
//return false when the token match, true otherwise
func MatchNotAnyToken(err IncompleteError, tList ...text.TokenType) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		if len(tokens) < 1 {
			pState.Pull(&tokens)
		}
		match := false
		for _, t := range tList {
			match = t == tokens[0].EnumType()
			if match {
				break
			}
		}
		if match {
			pState.ReportError(tokens[0], err.errType, err.errMsg)
			return tokens, err.errType == NoError
		}
		pState.Push(tokens[0:1], DefaultAppend)
		return tokens[1:], true
	}
}

//TryMatchToken return a MatchRule that match a specified TokenType
//skip tell the test to skip the tested token
//return false if the test is false, true otherwise
func TryMatchToken(tokenType text.TokenType, skip bool) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		if len(tokens) < 1 {
			pState.Pull(&tokens)
		}
		match := tokens[0].EnumType() == tokenType
		if skip && match {
			return tokens[1:], match
		}
		return tokens, match
	}
}

//TryNotMatchToken return a MatchRule that match a specified TokenType
//return true if the test is false, false otherwise
func TryNotMatchToken(tokenType text.TokenType, skip bool) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		if len(tokens) < 1 {
			pState.Pull(&tokens)
		}
		match := tokens[0].EnumType() != tokenType
		if skip && match {
			return tokens[1:], match
		}
		return tokens, match
	}
}

//TryMatchAnyString return a MatchRule test the first token in the slice with the arguments
//if a match is found the MatchRule return true, false otherwise
func TryMatchAnyString(sList ...string) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		if len(tokens) < 1 {
			pState.Pull(&tokens)
		}
		match := false
		for _, s := range sList {
			match = strings.EqualFold(tokens[0].Value(), s)
			if match {
				break
			}
		}
		return tokens, match
	}
}

//MatchAnyString return a MatchRule that match any string in a list
//return false with the same MatchRule condition and if an error is found
func MatchAnyString(err IncompleteError, sList ...string) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		if len(tokens) < 1 {
			pState.Pull(&tokens)
		}
		match := false
		for _, s := range sList {
			match = strings.EqualFold(tokens[0].Value(), s)
			if match {
				break
			}
		}
		if !match {
			pState.ReportError(tokens[0], err.errType, err.errMsg)
			return tokens, err.errType == NoError
		}
		pState.Push(tokens[0:1], DefaultAppend)
		return tokens[1:], true
	}
}

//MatchAnyToken return a MatchRule that match any token in a list
//return false with the same MatchRule condition and if an error is found
func MatchAnyToken(err IncompleteError, tList ...text.TokenType) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		if len(tokens) < 1 {
			pState.Pull(&tokens)
		}
		match := false
		for _, t := range tList {
			match = t == tokens[0].EnumType()
			if match {
				break
			}
		}
		if !match {
			pState.ReportError(tokens[0], err.errType, err.errMsg)
			return tokens, err.errType == NoError
		}
		pState.Push(tokens[0:1], DefaultAppend)
		return tokens[1:], true
	}
}

//MatchAll combine multiple MatchRule together and return a MatchRule
//that match all of the MatchRule in sequence
func MatchAll(tag int, matchList ...MatchRule) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		tokensLeft := tokens
		keepMatch := true
		markT, markL := pState.Mark()
		for _, match := range matchList {
			tokensLeft, keepMatch = match(tokensLeft, pState)
			if !keepMatch {
				return tokensLeft, false
			}
		}
		pState.Push([]text.Token{}, tag)
		pState.ToCST(markT, markL, tag)
		return tokensLeft, true
	}
}

//TryMatch combine a test and a match and return a MatchRule
//that parse the token stream only if the test match
func TryMatch(tag int, test MatchRule, match MatchRule) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		tokensLeft, keepMatch := test(tokens, pState)
		markT, markL := pState.Mark()
		if !keepMatch {
			return tokensLeft, true
		}
		tokensLeft, keepMatch = match(tokensLeft, pState)
		if !keepMatch {
			return tokensLeft, false
		}
		pState.Push([]text.Token{}, tag)
		pState.ToCST(markT, markL, tag)
		return tokensLeft, true
	}
}

//TryMatchRepeat combine a test and a match and return a MatchRule
//that parse the token stream while the test match
func TryMatchRepeat(tag int, test MatchRule, match MatchRule) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		tMatch, testSuccess := test(tokens, pState)
		ruleMatch := true
		for testSuccess && ruleMatch {
			markT, markL := pState.Mark()
			tMatch, ruleMatch = match(tMatch, pState)
			if !ruleMatch {
				return tMatch, false
			}
			pState.Push([]text.Token{}, tag)
			pState.ToCST(markT, markL, tag)
			tMatch, testSuccess = test(tMatch, pState)
		}
		return tMatch, true
	}
}

//MultiMatch is a collection of matches with their initial token
type MultiMatch struct {
	matches map[text.TokenType]MatchRule
}

//NewMultiMatch return a MultiMatch object
func NewMultiMatch() MultiMatch {
	return MultiMatch{map[text.TokenType]MatchRule{}}
}

//AddMatch add a match to the match map
func (mMatch *MultiMatch) AddMatch(tType text.TokenType, match MatchRule) {
	mMatch.matches[tType] = match
}

//MatchWithMap return a MatchRule that select a match using the first token found
//that MatchRule return false if EOF is reached or if the token is not known
func (mMatch *MultiMatch) MatchWithMap(err IncompleteError) MatchRule {
	return func(tokens []text.Token, pState ParserState) ([]text.Token, bool) {
		if len(tokens) < 1 {
			pState.Pull(&tokens)
		}
		tType := tokens[0].EnumType()
		match, exist := mMatch.matches[tType]
		if !exist {
			println(tType.String())
			pState.ReportError(tokens[0], err.errType, err.errMsg)
			return tokens, err.errType == NoError
		}
		tResult, keepParsing := match(tokens, pState)
		return tResult, keepParsing
	}
}
