package operators

var Operators = map[string]int{
	"||": 0, "&&": 0,
	"!=": 1, "==": 1, ">=": 1, "<=": 1, "<": 1, ">": 1, ".in": 1, ".get": 1,
	"<<": 2, ">>": 2, "+": 2, "-": 2, "|": 2,
	"*": 3, "/": 3, "%": 3, "^": 3, "&": 3,
	"!": 4, "~": 4, ".len": 4,
}

var OperatorEval = map[string](func(...int) int){
	"||": EvalLogicalOR, "&&": EvalLogicalAnd,
	"!=": EvalNotEqual, "==": EvalEqual, ">=": EvalGreaterEqualThan, "<=": EvalLessEqualThan, "<": EvalLessThan, ">": EvalGreaterThan,
	"<<": EvalShiftLeft, ">>": EvalShiftRight, "+": EvalAdd, "-": EvalSub, "|": EvalOr,
	"*": EvalMul, "/": EvalDiv, "%": EvalRem, "^": EvalXor, "&": EvalAnd,
	"!": EvalLogicalNot, "~": EvalNot,
}

//EvalOperation is a function that eval a single operation
type EvalOperation func(operands ...int) int

//Eval an operation with variable number of operands
func Eval(operator string, operands ...int) int {
	val, _ := OperatorEval[operator]
	return val(operands...)
}

func EvalShiftLeft(operands ...int) int {
	return operands[0] << uint(operands[1])
}

func EvalShiftRight(operands ...int) int {
	return operands[0] >> uint(operands[1])
}

func EvalAdd(operands ...int) int {
	return operands[0] + operands[1]
}

func EvalSub(operands ...int) int {
	return operands[0] - operands[1]
}

func EvalOr(operands ...int) int {
	return operands[0] | operands[1]
}

func EvalMul(operands ...int) int {
	return operands[0] * operands[1]
}

func EvalDiv(operands ...int) int {
	return operands[0] / operands[1]
}

func EvalRem(operands ...int) int {
	return operands[0] % operands[1]
}

func EvalXor(operands ...int) int {
	return operands[0] ^ operands[1]
}

func EvalAnd(operands ...int) int {
	return operands[0] & operands[1]
}

func EvalNot(operands ...int) int {
	return ^operands[1]
}

func EvalLogicalNot(operands ...int) int {
	if operands[1] == 0 {
		return 1
	}
	return 0
}

func EvalLogicalOR(operands ...int) int {
	if operands[0] > 0 || operands[1] > 0 {
		return 1
	}
	return 0
}

func EvalLogicalAnd(operands ...int) int {
	if operands[0] > 0 && operands[1] > 0 {
		return 1
	}
	return 0
}

func EvalNotEqual(operands ...int) int {
	if operands[0] != operands[1] {
		return 1
	}
	return 0
}

func EvalEqual(operands ...int) int {
	if operands[0] == operands[1] {
		return 1
	}
	return 0
}

func EvalGreaterEqualThan(operands ...int) int {
	if operands[0] >= operands[1] {
		return 1
	}
	return 0
}

func EvalLessEqualThan(operands ...int) int {
	if operands[0] <= operands[1] {
		return 1
	}
	return 0
}

func EvalLessThan(operands ...int) int {
	if operands[0] < operands[1] {
		return 1
	}
	return 0
}

func EvalGreaterThan(operands ...int) int {
	if operands[0] > operands[1] {
		return 1
	}
	return 0
}
