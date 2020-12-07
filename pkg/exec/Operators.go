package exec

type Operator func(i *Interpreter)

var ExecOperator = map[string]Operator{
	"&&": EvalLogicalAnd, "||": EvalLogicalOR, "==": EvalEqual, "!=": EvalNotEqual, ">": EvalGreaterThan, ">=": EvalGreaterEqualThan, "<=": EvalLessEqualThan,
	"<": EvalLessThan, "<<": EvalShiftLeft, ">>": EvalShiftRight, "+": EvalAdd, "-": EvalSub, "*": EvalMul, "/": EvalDiv, "%": EvalRem, "&": EvalAnd,
	"|": EvalOr, "^": EvalXor,
}

func EvalShiftLeft(i *Interpreter) {
	a := i.Pop()
	i.Push(i.Pop() << a)
}

func EvalShiftRight(i *Interpreter) {
	a := i.Pop()
	i.Push(i.Pop() >> a)
}

func EvalAdd(i *Interpreter) {
	i.Push(i.Pop() + i.Pop())
}

func EvalSub(i *Interpreter) {
	a := i.Pop()
	i.Push(i.Pop() - a)
}

func EvalOr(i *Interpreter) {
	i.Push(i.Pop() | i.Pop())
}

func EvalMul(i *Interpreter) {
	i.Push(i.Pop() * i.Pop())
}

func EvalDiv(i *Interpreter) {
	i.Push(i.Pop() / i.Pop())
}

func EvalRem(i *Interpreter) {
	i.Push(i.Pop() % i.Pop())
}

func EvalXor(i *Interpreter) {
	i.Push(i.Pop() ^ i.Pop())
}

func EvalAnd(i *Interpreter) {
	i.Push(i.Pop() & i.Pop())
}

func EvalLogicalOR(i *Interpreter) {
	b := i.Pop()
	a := i.Pop()
	if a != 0 || b != 0 {
		i.Push(1)
	} else {
		i.Push(0)
	}
}

func EvalLogicalAnd(i *Interpreter) {
	b := i.Pop()
	a := i.Pop()
	if a != 0 && b != 0 {
		i.Push(1)
	} else {
		i.Push(0)
	}
}

func EvalNotEqual(i *Interpreter) {
	b := i.Pop()
	a := i.Pop()
	if a != b {
		i.Push(1)
	} else {
		i.Push(0)
	}
}

func EvalEqual(i *Interpreter) {
	b := i.Pop()
	a := i.Pop()
	if a == b {
		i.Push(1)
	} else {
		i.Push(0)
	}
}

func EvalGreaterEqualThan(i *Interpreter) {
	b := i.Pop()
	a := i.Pop()
	if a >= b {
		i.Push(1)
	} else {
		i.Push(0)
	}
}

func EvalLessEqualThan(i *Interpreter) {
	b := i.Pop()
	a := i.Pop()
	if a <= b {
		i.Push(1)
	} else {
		i.Push(0)
	}
}

func EvalLessThan(i *Interpreter) {
	b := i.Pop()
	a := i.Pop()
	if a < b {
		i.Push(1)
	} else {
		i.Push(0)
	}
}

func EvalGreaterThan(i *Interpreter) {
	b := i.Pop()
	a := i.Pop()
	if a > b {
		i.Push(1)
	} else {
		i.Push(0)
	}
}
