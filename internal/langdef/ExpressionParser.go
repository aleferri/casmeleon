package langdef

import (
	"bitbucket.org/mrpink95/casmeleon/internal/operators"
	"bitbucket.org/mrpink95/casmeleon/internal/parsing"
	"bitbucket.org/mrpink95/casmeleon/internal/text"
)

func matchOperand(tokens []text.Token, pState parsing.ParserState) ([]text.Token, bool) {
	if len(tokens) < 1 {
		pState.Pull(&tokens)
	}
	isOperand := tokens[0].EnumType() == text.Number || tokens[0].EnumType() == text.Identifier
	if isOperand {
		pState.Push(tokens[0:1], parsing.DefaultAppend)
		return tokens[1:], true
	}
	if tokens[0].EnumType() != text.LParen {
		pState.ReportError(tokens[0], parsing.ErrorExpectedToken, ", expected number, identifer or '('")
		return tokens, false
	}
	tokensLeft, hasExpression := matchExpression(tokens[1:], pState)
	if hasExpression {
		matchRPar := parsing.MatchSkipToken(parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected ')'"), text.RParen)
		return matchRPar(tokensLeft, pState)
	}
	pState.ReportError(tokens[0], parsing.ErrorIncompleteExpression, ", no <expression> and ')' after '('")
	return tokens, false
}

var testBinaryOperator = parsing.TryMatchAnyString("+", "-", "*", "/", "%", "<<", ">>", "^", "&", "|", "&&", "||", "==", "!=", ">", "<", ">=", "<=", ".in", ".get")

func matchBinaryOperation(tokens []text.Token, pState parsing.ParserState, oldPrecedence int) ([]text.Token, bool) {
	tokens, exist := testBinaryOperator(tokens, pState)
	if !exist {
		return tokens, true
	}
	op := tokens[0]
	precedence := operators.Operators[op.Value()]
	if oldPrecedence >= precedence {
		return tokens, true
	}
	tokensLeft, isUnary := matchUnaryOperation(tokens[1:], pState)
	if !isUnary {
		return tokensLeft, false
	}
	var isBinary bool
	tokensLeft, isBinary = matchBinaryOperation(tokensLeft, pState, precedence)
	if !isBinary {
		return tokensLeft, false
	}
	pState.Push(tokens[0:1], parsing.DefaultAppend)
	return matchBinaryOperation(tokensLeft, pState, oldPrecedence)
}

var matchUnaryOperator = parsing.TryMatchAnyString("+", "-", "~", "!", ".len")

func matchUnaryOperation(tokens []text.Token, pState parsing.ParserState) ([]text.Token, bool) {
	if len(tokens) < 1 {
		pState.Pull(&tokens)
	}
	tokensLeft, isUnary := matchUnaryOperator(tokens, pState)
	if isUnary {
		var noError bool
		tokens[0] = tokens[0].WithType(text.UnaryOperator)
		tokensLeft, noError = matchOperand(tokensLeft[1:], pState)
		pState.Push(tokens[0:1], parsing.DefaultAppend)
		return tokensLeft, noError
	}
	return matchOperand(tokensLeft, pState)
}

func matchExpression(tokens []text.Token, pState parsing.ParserState) ([]text.Token, bool) {
	tokensLeft, noError := matchUnaryOperation(tokens, pState)
	if noError {
		tokensLeft, noError = matchBinaryOperation(tokensLeft, pState, -1)
	}
	return tokensLeft, noError
}

//MatchExpression match a full expression in the parse stream
func MatchExpression(tokens []text.Token, pState parsing.ParserState) ([]text.Token, bool) {
	markT, markL := pState.Mark()
	tokensLeft, noError := matchExpression(tokens, pState)
	if noError {
		pState.ToCST(markT, markL, TagExpression)
	}
	return tokensLeft, noError
}
