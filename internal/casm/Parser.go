package casm

import (
	"errors"
	"fmt"

	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

//ParseNumberBase parse a number base directive
func ParseNumberBase(stream parser.Stream) (parser.CSTNode, error) {
	seq, err := parser.RequireSequence(stream, text.KeywordNum, text.Number, text.QuotedString, text.QuotedString)
	if err != nil {
		return nil, err
	}
	return parser.BuildLeaf(seq, NUMBER_BASE), nil
}

//ParseSet parse a set of symbols
func ParseSet(stream parser.Stream) (parser.CSTNode, error) {
	seq, err := parser.RequireSequence(stream, text.KeywordSet, text.Identifier)
	if err != nil {
		return nil, err
	}
	after, noInset := parser.AcceptInsetPattern(stream, text.CurlyOpen, text.CurlyClose, text.Identifier, text.Semicolon)
	if noInset != nil {
		return nil, noInset
	}
	list := parser.BuildLeaf(after, SYMBOL_SET)
	set := parser.BuildBranch(seq, SET_NODE)
	set.InsertChild(list, true)
	return set, nil
}

//ParseOpcode from the source stream
func ParseOpcode(stream parser.Stream) (parser.CSTNode, error) {
	seq, err := parser.RequireSequence(stream, text.KeywordOpcode, text.Identifier)
	if err != nil {
		return nil, err
	}
	args, noArgs := parser.AcceptInsetDelegate(stream, text.DoubleCurlyOpen, text.DoubleCurlyClose, ParseArgs)
	if noArgs != nil {
		return nil, noArgs
	}
	noWith := parser.Expect(stream, text.KeywordWith)
	if noWith != nil {
		return nil, noWith
	}
	withContent, invalidWith := parser.AcceptPatternWithTest(stream, text.RoundOpen, text.RoundClose, text.Comma, ParseWithArgs)
	if invalidWith != nil {
		return nil, invalidWith
	}
	noArrow := parser.Expect(stream, text.SymbolArrow)
	if noArrow != nil {
		return nil, noArrow
	}

	body, bodyErr := ParseBlock(stream)
	if bodyErr != nil {
		return nil, bodyErr
	}

	opcodeNode := parser.BuildBranch(seq, OPCODE_NODE)
	for _, arg := range args {
		opcodeNode.InsertChild(arg, true)
	}

	opcodeWith := parser.BuildBranch([]text.Symbol{}, WITH_TYPES)
	for _, withArg := range withContent {
		opcodeWith.InsertChild(withArg, true)
	}
	opcodeNode.InsertChild(opcodeWith, true)
	opcodeNode.InsertChild(body, true)

	return opcodeNode, nil
}

//ParseInline from the source stream
func ParseInline(stream parser.Stream) (parser.CSTNode, error) {
	seq, err := parser.RequireSequence(stream, text.KeywordInline, text.Identifier, text.KeywordWith)
	if err != nil {
		return nil, err
	}
	withContent, invalidWith := parser.AcceptPatternWithTest(stream, text.RoundOpen, text.RoundClose, text.Comma, ParseWithArgs)
	if invalidWith != nil {
		return nil, invalidWith
	}
	noArrow := parser.Expect(stream, text.SymbolArrow)
	if noArrow != nil {
		return nil, noArrow
	}

	body, bodyErr := ParseBlock(stream)
	if bodyErr != nil {
		return nil, bodyErr
	}

	inlineNode := parser.BuildBranch(seq, INLINE_NODE)
	inlineWith := parser.BuildBranch([]text.Symbol{}, WITH_TYPES)
	for _, withArg := range withContent {
		inlineWith.InsertChild(withArg, true)
	}
	inlineNode.InsertChild(inlineWith, true)
	inlineNode.InsertChild(body, true)

	return inlineNode, nil
}

//ParseWithArgs after the opcode or inline declaration
func ParseWithArgs(stream parser.Stream) (parser.CSTNode, error) {
	arg, err := parser.RequireSequence(stream, text.Identifier, text.Colon, text.Identifier)
	return parser.BuildLeaf(arg, OPCODE_ARGS), err
}

//ParseArgs parse the arguments of the opcode
func ParseArgs(stream parser.Stream) (parser.CSTNode, error) {
	acc := []text.Symbol{}
	for stream.Peek().ID() != text.DoubleCurlyClose {
		acc = append(acc, stream.Next())
	}
	return parser.BuildLeaf(acc, OPCODE_ARGS), nil
}

//ParseStatement inside a block
func ParseStatement(stream parser.Stream) (parser.CSTNode, error) {
	switch stream.Peek().ID() {
	case text.KeywordIF:
		return ParseBranch(stream)
	case text.KeywordError:
		return ParseError(stream)
	case text.KeywordWarning:
		return ParseWarning(stream)
	case text.KeywordOut:
		return ParseOut(stream)
	case text.KeywordReturn:
		return ParseReturn(stream)
	}
	return nil, errors.New("Unexpected symbol '" + stream.Peek().Value() + "', expecting '.if' or '.error' or '.warning' or '.out' or '.return'")
}

//ParseBlock of code
func ParseBlock(stream parser.Stream) (parser.CSTNode, error) {
	body, err := parser.AcceptInsetDelegate(stream, text.CurlyOpen, text.CurlyClose, ParseStatement)
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
	parser.Expect(stream, text.KeywordIF)
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
	parser.Expect(stream, text.KeywordOut)
	exprs, err := parser.AcceptInsetDelegate(stream, text.SquareOpen, text.SquareClose, ParseExpression)

	parser.Expect(stream, text.Semicolon)

	outNode := parser.BuildBranch([]text.Symbol{}, STMT_OUT)
	for _, e := range exprs {
		outNode.InsertChild(e, true)
	}
	return outNode, err
}

//ParseReturn parse an out
func ParseReturn(stream parser.Stream) (parser.CSTNode, error) {
	parser.Expect(stream, text.KeywordReturn)
	expr, err := ParseExpression(stream)

	p := stream.Next()
	if p.ID() != text.Semicolon {
		return nil, fmt.Errorf("Expected ';', found '%s", p.Value())
	}

	outNode := parser.BuildBranch([]text.Symbol{}, STMT_RET)
	outNode.InsertChild(expr, err != nil)
	return outNode, err
}

//ParseError parse an out
func ParseError(stream parser.Stream) (parser.CSTNode, error) {
	syms, err := parser.RequireSequence(stream, text.KeywordError, text.Identifier, text.Comma, text.QuotedString, text.Semicolon)
	if err != nil {
		return nil, err
	}
	return parser.BuildLeaf(syms, STMT_WARNING), nil
}

//ParseWarning parse an out
func ParseWarning(stream parser.Stream) (parser.CSTNode, error) {
	syms, err := parser.RequireSequence(stream, text.KeywordWarning, text.Identifier, text.Comma, text.QuotedString, text.Semicolon)
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
	if t.ID() == text.Number || t.ID() == text.Identifier {
		return append(expr, t), nil
	}

	if t.ID() == text.RoundOpen {
		symbols, errExpr := parseExpression(stream, expr)
		if errExpr != nil {
			return nil, errExpr
		}
		err := parser.Expect(stream, text.RoundClose)
		return symbols, err
	}

	if testAnyUnaryOperator(t) {
		return parseTerm(stream, append(expr, t))
	}

	if t.ID() == text.KeywordExpr {
		part1, invalidStart := parser.RequireSequence(stream, text.Identifier, text.RoundOpen)
		if invalidStart != nil {
			return expr, invalidStart
		}

		expr = append(expr, part1...)

		begin := true

		for begin || stream.Peek().ID() == text.Comma {
			arg, invalidArg := parseExpression(stream, expr)
			if invalidArg != nil {
				return expr, invalidArg
			}

			expr = append(expr, arg...)
		}

		part3, invalidEnd := parser.Require(stream, text.RoundClose)
		if invalidEnd != nil {
			return expr, invalidEnd
		}

		expr = append(expr, part3)
	}

	return expr, errors.New("Unexpected Symbol " + t.Value())
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
