package casm

import (
	"errors"
	"strconv"

	"github.com/aleferri/casmeleon/pkg/exec"
	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

//PruneExpressionToExecutionList rewalk through the expression to generate the execution list
func PruneExpressionToExecutionList(lang *Language, params []string, types []uint32, expr parser.CSTNode) ([]exec.Executable, error) {
	list := []exec.Executable{}

	left, err := CompileTerm(lang, params, &list, expr.Symbols())
	for len(left) > 0 && err == nil {
		left, err = CompileTerm(lang, params, &list, left)
	}

	return list, err
}

var Precedence = map[string]int{
	"||": 0, "&&": 0,
	"!=": 1, "==": 1, ">=": 1, "<=": 1, "<": 1, ">": 1, ".in": 1, ".get": 1,
	"<<": 2, ">>": 2, "+": 2, "-": 2, "|": 2,
	"*": 3, "/": 3, "%": 3, "^": 3, "&": 3,
	"!": 4, "~": 4, ".len": 4,
}

//CompileTerm compile a <Term> of the expression: either <Identifier> | <Integer> | <UnaryOp> | <ParensExpr> | <InlineCall>
func CompileTerm(lang *Language, params []string, list *[]exec.Executable, queued []text.Symbol) ([]text.Symbol, error) {
	if len(queued) == 0 {
		return queued, nil
	}

	q := queued[0]

	switch q.ID() {
	case text.Number:
		{
			base := 10
			str := q.Value()
			if len(str) > 1 {
				if str[1] == 'b' {
					base = 2
					str = str[2:]
				} else if str[1] == 'x' {
					base = 16
					str = str[2:]
				}
			}
			v, e := strconv.ParseInt(str, base, 32)
			if e != nil {
				return queued[1:], e
			}
			*list = append(*list, exec.ILoadOf(v))
			return queued[1:], nil
		}
	case text.Identifier:
		{
			for i, p := range params {
				if p == q.Value() {
					*list = append(*list, exec.RLoadOf(uint32(i)))
					if p == ".addr" {
						lang.MarkAddressUsed()
					}
					return queued[1:], nil
				}
			}
			t, f := lang.SetOf(q.Value())
			if f {
				*list = append(*list, exec.ILoadOf(int64(t.valueOf(q.Value()))))
			}
			return queued, errors.New("Parameter " + q.Value() + " not found")
		}
	case text.RoundOpen:
		{
			var err error = nil
			stack := []exec.Executable{}
			for q.ID() != text.RoundClose && err == nil {
				queued, err = CompileExpression(lang, params, &stack, queued[1:])
				q = queued[0]
			}
			*list = append(*list, exec.BuildStackExpression(stack))
			return queued[1:], nil
		}
	case text.OperatorNeg, text.OperatorNot, text.OperatorPlusUnary, text.OperatorMinusUnary:
		{
			stack := []exec.Executable{}
			left, err := CompileTerm(lang, params, &stack, queued[1:])

			var partial exec.Executable = exec.BuildStackExpression(stack)
			if q.ID() != text.OperatorPlus {
				switch q.ID() {
				case text.OperatorNeg:
					partial = exec.BuildComplement(partial)
				case text.OperatorMinus:
					partial = exec.BuildNegate(partial)
				case text.OperatorNot:
					partial = exec.BuildNot(partial)
				}
			}
			*list = append(*list, partial)
			return left, err
		}
	case text.KeywordExpr:
		{
			stack := []exec.Executable{}
			inlineName := queued[1]

			q = queued[2].WithID(text.Comma)
			queued = queued[2:]
			var err error = nil
			for q.ID() == text.Comma && err == nil {
				queued, err = CompileExpression(lang, params, &stack, queued[1:])
				q = queued[0]
			}

			inlineCode, err := lang.InlineCode(inlineName.Value())

			if err != nil {
				return queued, err
			}

			*list = append(*list, exec.BuildStackExpression(stack), inlineCode)
			return queued[1:], nil
		}
	}
	return queued, nil
}

//CompileFactor compile the left associativity part of the expression
func CompileFactor(lang *Language, params []string, list *[]exec.Executable, precedence *text.Symbol, queued []text.Symbol) ([]text.Symbol, error) {
	operator := queued[0]

	opVal := operator.Value()
	opPrec, isOp := Precedence[opVal]

	if !isOp {
		return queued, nil
	}

	left, err := CompileTerm(lang, params, list, queued[1:])
	if err != nil {
		return left, err
	}

	if len(left) == 0 {
		*list = append(*list, exec.BuildReduce(opVal))
		return left, nil
	}

	other := left[0]

	seq := other.Value()
	seqPrec, isOp := Precedence[seq]

	if isOp && seqPrec > opPrec {
		left, err = CompileFactor(lang, params, list, &other, left)
		if err != nil {
			return left, err
		}
		*list = append(*list, exec.BuildReduce(opVal))
	} else {
		*list = append(*list, exec.BuildReduce(opVal))
	}

	return CompileFactor(lang, params, list, nil, left)
}

//CompileExpression compile the whole expression
func CompileExpression(lang *Language, params []string, list *[]exec.Executable, queued []text.Symbol) ([]text.Symbol, error) {
	left, err := CompileTerm(lang, params, list, queued)
	if err != nil || len(left) == 0 {
		return left, err
	}

	return CompileFactor(lang, params, list, nil, left)
}
