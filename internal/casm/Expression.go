package casm

import (
	"errors"
	"strconv"

	"github.com/aleferri/casmeleon/pkg/expr"
	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/operators"
)

const USE_THIS_ADDR = 1

// WalkCSTExpression walk the concrete syntax tree of an expression to convert it into SSA form
func WalkCSTExpression(lang *Language, params []string, types []uint32, node parser.CSTNode) (expr.Converter, error) {
	listing := []opcodes.Opcode{}
	status := expr.MakeConverter(node.Symbols(), uint16(len(params)))

	err := CompileTerm(lang, params, &listing, &status)
	for len(status.Queue()) > 0 && err == nil {
		err = CompileTerm(lang, params, &listing, &status)
	}
	return status, err
}

var Precedence = map[string]int{
	"||": 0, "&&": 0,
	"!=": 1, "==": 1, ">=": 1, "<=": 1, "<": 1, ">": 1, ".in": 1, ".get": 1,
	"<<": 2, ">>": 2, "+": 2, "-": 2, "|": 2,
	"*": 3, "/": 3, "%": 3, "^": 3, "&": 3,
	"!": 4, "~": 4, ".len": 4,
}

// CompileTerm compile a <Term> of the expression: either <Identifier> | <Integer> | <UnaryOp> | <ParensExpr> | <InlineCall>
func CompileTerm(lang *Language, params []string, listing *[]opcodes.Opcode, status *expr.Converter) error {
	if status.IsEmptyQueue() {
		return nil
	}

	q := status.Poll()

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
				return e
			}
			atom := status.LabelAtom(expr.MakeLiteral(str, v))
			*listing = append(*listing, opcodes.MakeIConst(atom.Local(), v))
			status.Push(atom)
			return nil
		}
	case text.Identifier:
		{
			for i, p := range params {
				if p == q.Value() {
					status.Push(expr.MakeParameter(p, int64(i), uint16(i)))
					if p == ".addr" {
						status.SetFlag(USE_THIS_ADDR)
					}
					return nil
				}
			}
			t, f := lang.SetOf(q.Value())
			if f {
				refId := t.valueOf(q.Value())
				atom := status.LabelAtom(expr.MakeMember(q.Value(), int64(refId)))
				*listing = append(*listing, opcodes.MakeIConst(atom.Local(), atom.Value()))
				status.Push(atom)
				return nil
			}
			return errors.New("Parameter " + q.Value() + " not found")
		}
	case text.RoundOpen:
		{
			var err error = nil
			for q.ID() != text.RoundClose && err == nil {
				err = CompileExpression(lang, params, listing, status)
				q = status.Front()
			}
			status.Poll()
			return err
		}
	case text.OperatorNeg, text.OperatorNot, text.OperatorPlusUnary, text.OperatorMinusUnary:
		{
			err := CompileTerm(lang, params, listing, status)

			if q.ID() == text.OperatorPlusUnary || err != nil {
				return err
			}

			a := status.Pop()
			resultLocal := status.LabelLocal()
			status.Push(expr.MakeLocal("local", 0, resultLocal))

			switch q.ID() {
			case text.OperatorNeg:
				*listing = append(*listing, opcodes.MakeUnaryOp(resultLocal, "com", opcodes.IntShape, a.Local(), operators.UnaryOperatorsSymbols["~"]))
			case text.OperatorMinusUnary:
				*listing = append(*listing, opcodes.MakeUnaryOp(resultLocal, "neg", opcodes.IntShape, a.Local(), operators.UnaryOperatorsSymbols["-"]))
			case text.OperatorNot:
				*listing = append(*listing, opcodes.MakeUnaryOp(resultLocal, "not", opcodes.IntShape, a.Local(), operators.UnaryOperatorsSymbols["!"]))
			}

			return nil
		}
	case text.KeywordExpr:
		{
			funcName := status.Poll()

			stack := expr.MakeConverter(status.Queue(), status.LabelLocal())

			q = stack.Poll().WithID(text.Comma) // open paren
			var err error = nil
			for q.ID() == text.Comma && err == nil {
				err = CompileExpression(lang, params, listing, &stack)
				q = stack.Poll() //comma or close paren
			}

			diff := len(status.Queue()) - len(stack.Queue())
			status.DropFront(uint(diff))

			if err != nil {
				return err
			}

			refs := []uint16{}
			for !stack.IsEmptyStack() {
				refs = append(refs, stack.Pop().Local())
			}

			addr, found := lang.FindAddressOf(funcName.Value())
			if !found {
				return errors.New("Cannot find function " + funcName.Value())
			}

			retLabel := status.LabelLocal()
			call := opcodes.MakeEnter([]uint16{retLabel}, addr, refs)

			*listing = append(*listing, call)
			status.Push(expr.MakeLocal("ret", 0, retLabel))
			return nil
		}
	}
	return nil
}

func ReduceBinaryExpression(op string, listing *[]opcodes.Opcode, status *expr.Converter) expr.Atom {
	resultLocal := status.LabelLocal()
	operation := operators.BinaryOperatorsSymbols[op]
	b := status.Pop()
	a := status.Pop()
	res := expr.MakeLocal("local", 0, resultLocal)
	*listing = append(*listing, opcodes.MakeBinaryOp(resultLocal, op, opcodes.IntShape, a.Local(), b.Local(), operation))
	status.Push(res)
	return res
}

// CompileFactor compile the left associativity part of the expression
func CompileFactor(lang *Language, params []string, listing *[]opcodes.Opcode, status *expr.Converter) error {
	if status.IsEmptyQueue() {
		return nil
	}
	operator := status.Front()

	opVal := operator.Value()
	opPrec, isOp := Precedence[opVal]

	if !isOp {
		return nil
	}

	status.Poll()
	err := CompileTerm(lang, params, listing, status)
	if err != nil {
		return err
	}

	if status.IsEmptyQueue() {
		ReduceBinaryExpression(opVal, listing, status)
		return nil
	}

	other := status.Front()

	seq := other.Value()
	seqPrec, isOp := Precedence[seq]

	if isOp && seqPrec > opPrec {
		err = CompileFactor(lang, params, listing, status)
		if err != nil {
			return err
		}
	}
	ReduceBinaryExpression(opVal, listing, status)

	return CompileFactor(lang, params, listing, status)
}

func CompileExpression(lang *Language, params []string, listing *[]opcodes.Opcode, status *expr.Converter) error {
	err := CompileTerm(lang, params, listing, status)
	if err != nil || status.IsEmptyQueue() {
		return err
	}

	return CompileFactor(lang, params, listing, status)
}
