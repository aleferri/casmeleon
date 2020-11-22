package language

import (
	"errors"

	reserved "github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

//ParseNumberBase parse a number base directive
func ParseNumberBase(stream parser.Stream, report parser.ErrorOf) (parser.CSTNode, bool) {
	seq, err := parser.RequireSequence(stream, reserved.KEYWORD_NUM, text.NUMBER, text.QUOTED_STRING, text.QUOTED_STRING)
	if err != nil {
		report(stream.Peek(), stream.Source(), err.Error())
		return nil, false
	}
	return parser.BuildLeaf(seq, NUMBER_BASE), true
}

//ParseSet parse a set of symbols
func ParseSet(stream parser.Stream, report parser.ErrorOf) (parser.CSTNode, bool) {
	seq, err := parser.RequireSequence(stream, reserved.KEYWORD_SET, text.IDENTIFIER)
	if err != nil {
		report(stream.Peek(), stream.Source(), err.Error())
		return nil, false
	}
	after, noInset := parser.AcceptInsetPattern(stream, text.CURLY_OPEN, text.CURLY_CLOSE, text.IDENTIFIER, text.SEMICOLON)
	if noInset != nil {
		report(stream.Peek(), stream.Source(), noInset.Error())
		return nil, false
	}
	list := parser.BuildLeaf(after, SYMBOL_SET)
	set := parser.BuildBranch(seq, SET_NODE)
	set.InsertChild(list, true)
	return set, true
}

//ParseOpcode from the source stream
func ParseOpcode(stream parser.Stream, report parser.ErrorOf) (parser.CSTNode, bool) {
	seq, err := parser.RequireSequence(stream, reserved.KEYWORD_OPCODE, text.IDENTIFIER)
	if err != nil {
		report(stream.Peek(), stream.Source(), err.Error())
		return nil, false
	}
	args, noArgs := parser.AcceptInsetDelegate(stream, reserved.DOUBLE_CURLY_OPEN, reserved.DOUBLE_CURLY_CLOSE, ParseArgs)
	if noArgs != nil {
		report(stream.Peek(), stream.Source(), noArgs.Error())
		return nil, false
	}
	noWith := parser.Expect(stream, reserved.KEYWORD_WITH)
	if noWith != nil {
		report(stream.Peek(), stream.Source(), noWith.Error())
		return nil, false
	}
	withContent, noParens := parser.AcceptInsetPattern(stream, text.ROUND_OPEN, text.ROUND_CLOSE, text.IDENTIFIER, text.COLON, text.IDENTIFIER)
	if noParens != nil {
		report(stream.Peek(), stream.Source(), noParens.Error())
		return nil, false
	}
	noArrow := parser.Expect(stream, reserved.ARROW)
	if noArrow != nil {
		report(stream.Peek(), stream.Source(), noArrow.Error())
		return nil, false
	}

	body, bodyErr := ParseBlock(stream)
	if bodyErr != nil {
		report(stream.Peek(), stream.Source(), bodyErr.Error())
		return nil, false
	}

	opcodeNode := parser.BuildBranch(seq, OPCODE_NODE)
	for _, arg := range args {
		opcodeNode.InsertChild(arg, true)
	}
	opcodeWith := parser.BuildLeaf(withContent, WITH_TYPES)
	opcodeNode.InsertChild(opcodeWith, true)
	opcodeNode.InsertChild(body, true)

	return opcodeNode, false
}

//ParseArgs parse the arguments of the opcode
func ParseArgs(stream parser.Stream) (parser.CSTNode, error) {
	acc := []text.Symbol{}
	for stream.Peek().ID() != reserved.DOUBLE_CURLY_CLOSE {
		acc = append(acc, stream.Next())
	}
	return parser.BuildLeaf(acc, OPCODE_ARGS), nil
}

//ParseStatement inside a block
func ParseStatement(stream parser.Stream) (parser.CSTNode, error) {
	switch stream.Peek().ID() {
	case reserved.KEYWORD_IF:
		return ParseBranch(stream)
	case reserved.KEYWORD_ERROR:
		return ParseError(stream)
	case reserved.KEYWORD_WARNING:
		return ParseWarning(stream)
	case reserved.KEYWORD_OUT:
		return ParseOut(stream)
	}
	return nil, errors.New("Unexpected symbol '" + stream.Peek().Value() + "'")
}

//ParseBlock of code
func ParseBlock(stream parser.Stream) (parser.CSTNode, error) {
	body, err := parser.AcceptInsetDelegate(stream, text.CURLY_OPEN, text.CURLY_CLOSE, ParseStatement)
	if err != nil {
		return nil, err
	}

	block := parser.BuildBranch([]text.Symbol{}, STMT_BLOCK)
	for _, arg := range body {
		block.InsertChild(arg, true)
	}
	return block, nil
}

//ParseBranch parse a branch
func ParseBranch(stream parser.Stream) (parser.CSTNode, error) {
	parser.Expect(stream, reserved.KEYWORD_IF)
	expr, err := ParseExpression(stream)
	if err != nil {
		return nil, err
	}

	block, errBlock := ParseBlock(stream)
	if errBlock != nil {
		return nil, errBlock
	}

	branch := parser.BuildBranch([]text.Symbol{}, STMT_BRANCH)
	branch.InsertChild(expr, true)
	branch.InsertChild(block, true)
	return branch, nil
}

//ParseOut parse an out
func ParseOut(stream parser.Stream) (parser.CSTNode, error) {
	parser.Expect(stream, reserved.KEYWORD_OUT)
	exprs, err := parser.AcceptInsetDelegate(stream, text.SQUARE_OPEN, text.SQUARE_CLOSE, ParseExpression)

	outNode := parser.BuildBranch([]text.Symbol{}, STMT_OUT)
	for _, e := range exprs {
		outNode.InsertChild(e, true)
	}
	return outNode, err
}

//ParseError parse an out
func ParseError(stream parser.Stream) (parser.CSTNode, error) {
	syms, err := parser.RequireSequence(stream, reserved.KEYWORD_ERROR, text.IDENTIFIER, text.COMMA, text.QUOTED_STRING, text.SEMICOLON)
	if err != nil {
		return nil, err
	}
	return parser.BuildLeaf(syms, STMT_WARNING), nil
}

//ParseWarning parse an out
func ParseWarning(stream parser.Stream) (parser.CSTNode, error) {
	syms, err := parser.RequireSequence(stream, reserved.KEYWORD_WARNING, text.IDENTIFIER, text.COMMA, text.QUOTED_STRING, text.SEMICOLON)
	if err != nil {
		return nil, err
	}
	return parser.BuildLeaf(syms, STMT_WARNING), nil
}

var allBinaryOperators = []string{"+", "-", "*", "/", "%", "<<", ">>", "^", "&", "|", "&&", "||", "==", "!=", ">", "<", ">=", "<=", ".in", ".get"}
var allUnaryOperators = []string{"+", "-", "!", "~"}

func testAnyBinaryOperator(stream parser.Stream) bool {
	t := stream.Peek()
	for _, p := range allBinaryOperators {
		if t.Value() == p {
			return true
		}
	}
	return false
}

func testAnyUnaryOperator(t text.Symbol) bool {
	for _, p := range allBinaryOperators {
		if t.Value() == p {
			return true
		}
	}
	return false
}

func parseTerm(stream parser.Stream, expr []text.Symbol) ([]text.Symbol, error) {
	t := stream.Next()
	if t.ID() == text.NUMBER || t.ID() == text.IDENTIFIER {
		return append(expr, t), nil
	}

	if t.ID() == text.ROUND_OPEN {
		symbols, errExpr := parseExpression(stream, expr)
		if errExpr != nil {
			return nil, errExpr
		}
		err := parser.Expect(stream, text.ROUND_CLOSE)
		return symbols, err
	}

	if testAnyUnaryOperator(t) {
		return parseTerm(stream, append(expr, t))
	}

	if t.ID() == reserved.KEYWORD_EXPR {
		part1, invalidStart := parser.RequireSequence(stream, text.IDENTIFIER, text.ROUND_OPEN)
		if invalidStart != nil {
			return expr, invalidStart
		}

		expr = append(expr, part1...)

		begin := true

		for begin || stream.Peek().ID() == text.COMMA {
			arg, invalidArg := parseExpression(stream, expr)
			if invalidArg != nil {
				return expr, invalidArg
			}

			expr = append(expr, arg...)
		}

		part3, invalidEnd := parser.Require(stream, text.ROUND_CLOSE)
		if invalidEnd != nil {
			return expr, invalidEnd
		}

		expr = append(expr, part3)
	}

	return expr, errors.New("Unexpected Token " + t.Value())
}

func parseFactors(stream parser.Stream, expr []text.Symbol) ([]text.Symbol, error) {
	partial, errLeft := parseTerm(stream, expr)
	if errLeft != nil {
		return expr, errLeft
	}

	if testAnyBinaryOperator(stream) {
		return parseFactors(stream, append(partial, stream.Next()))
	}

	return partial, nil
}

func parseExpression(stream parser.Stream, expr []text.Symbol) ([]text.Symbol, error) {
	return parseFactors(stream, expr)
}

//ParseExpression of any kind
func ParseExpression(stream parser.Stream) (parser.CSTNode, error) {
	expr, err := parseExpression(stream, []text.Symbol{})
	if err != nil {
		return parser.BuildLeaf(expr, EXPRESSION), nil
	}
	return nil, err
}
