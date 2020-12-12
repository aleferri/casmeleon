package casm

import (
	"fmt"

	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

func ParseCasm(stream parser.Stream, repo text.Source) (parser.CSTNode, error) {
	errFound := false

	id := stream.Peek().ID()

	root := parser.BuildBranch([]text.Symbol{}, ROOT_NODE)
	var err error

	for id != text.EOF && !errFound {
		var cst parser.CSTNode

		switch id {
		case text.KeywordInline:
			{
				cst, err = ParseInline(stream)
			}
		case text.KeywordOpcode:
			{
				cst, err = ParseOpcode(stream)
			}
		case text.KeywordNum:
			{
				cst, err = ParseNumberBase(stream)
			}
		case text.KeywordSet:
			{
				cst, err = ParseSet(stream)
			}
		default:
			{
				fmt.Println(len(idDescriptor))
				err = fmt.Errorf("Undefined symbol '%s'", idDescriptor[id])
			}
		}

		if err != nil {
			errFound = true
		} else {
			root.InsertChild(cst, true)
		}

		id = stream.Peek().ID()
	}
	return root, err
}

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
		return nil, WrapMatchError(bodyErr, ".opcode", "}")
	}

	opcodeNode := parser.BuildBranch(seq, OPCODE_NODE)

	opcodeFormat := parser.BuildBranch(seq, OPCODEFORMAT)
	for _, arg := range args {
		opcodeFormat.InsertChild(arg, true)
	}

	opcodeNode.InsertChild(opcodeFormat, true)

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
		return nil, WrapMatchError(bodyErr, ".inline", "}")
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
	case text.KeywordOutR:
		return ParseOut(stream)
	case text.KeywordReturn:
		return ParseReturn(stream)
	}
	return nil, parser.ExpectedAnyOf(
		stream.Peek(), "Unexpected symbol '%s', was expecting: '%s'",
		text.KeywordIF, text.KeywordOut, text.KeywordReturn, text.KeywordError, text.KeywordWarning,
	)
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

	if stream.Peek().ID() == text.KeywordELSE {
		stream.Next()
		elseBlock, errElse := ParseBlock(stream)
		if errElse != nil {
			return nil, errBlock
		}
		branch.InsertChild(elseBlock, true)
	}

	return branch, nil
}

//ParseOut parse an out
func ParseOut(stream parser.Stream) (parser.CSTNode, error) {
	outK := stream.Next()
	exprs, err := parser.AcceptPatternWithTest(stream, text.SquareOpen, text.SquareClose, text.Comma, ParseExpression)

	parser.Expect(stream, text.Semicolon)

	id := uint32(STMT_OUT)
	if outK.ID() == text.KeywordOutR {
		id = STMT_OUTR
	}

	outNode := parser.BuildBranch([]text.Symbol{}, id)
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
		return nil, parser.ExpectedSymbol(p, "Found unexpected token '%s', expected '%s' instead", text.Semicolon)
	}

	outNode := parser.BuildBranch([]text.Symbol{}, STMT_RET)
	outNode.InsertChild(expr, err == nil)
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

func testAnyBinaryOperator(stream parser.Stream) bool {
	t := stream.Peek()
	if t.ID() == text.OperatorMinusUnary || t.ID() == text.OperatorPlusUnary {
		return false
	}
	for _, p := range allBinaryOperators {
		if t.Value() == p {
			return true
		}
	}
	return false
}

func parseTerm(stream parser.Stream, expr []text.Symbol) ([]text.Symbol, error) {
	t := stream.Next()

	switch t.ID() {
	case text.Number, text.Identifier:
		{
			return append(expr, t), nil
		}
	case text.RoundOpen:
		{
			symbols, errExpr := parseFactors(stream, append(expr, t))
			if errExpr != nil {
				return nil, errExpr
			}
			cls, err := parser.Require(stream, text.RoundClose)
			symbols = append(symbols, cls)
			return symbols, err
		}
	case text.OperatorNeg, text.OperatorNot, text.OperatorPlus, text.OperatorMinus:
		{
			if t.ID() == text.OperatorPlus {
				return parseTerm(stream, append(expr, t.WithID(text.OperatorPlusUnary)))
			} else if t.ID() == text.OperatorMinus {
				return parseTerm(stream, append(expr, t.WithID(text.OperatorMinusUnary)))
			}
			return parseTerm(stream, append(expr, t))
		}
	case text.KeywordExpr:
		{
			part1, invalidStart := parser.RequireSequence(stream, text.Identifier, text.RoundOpen)
			if invalidStart != nil {
				return expr, invalidStart
			}

			expr = append(expr, t)
			expr = append(expr, part1...)

			readNextParam := stream.Peek().ID() != text.RoundClose

			for readNextParam {
				arg, invalidArg := parseFactors(stream, expr)
				if invalidArg != nil {
					return expr, invalidArg
				}

				expr = arg

				if stream.Peek().ID() == text.Comma {
					expr = append(expr, stream.Next())
				} else {
					readNextParam = false
				}
			}

			part3, invalidEnd := parser.Require(stream, text.RoundClose)
			if invalidEnd != nil {
				return expr, invalidEnd
			}

			return append(expr, part3), nil
		}
	default:
		{
			return expr, parser.ExpectedSymbol(t, "Unexpected Symbol '%s', was expecting a term like %s", text.Identifier)
		}
	}
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

//ParseExpression of any kind
func ParseExpression(stream parser.Stream) (parser.CSTNode, error) {
	expr, err := parseFactors(stream, []text.Symbol{})
	if err == nil {
		return parser.BuildLeaf(expr, EXPRESSION), nil
	}
	return nil, err
}
