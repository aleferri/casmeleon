package exec

import "strconv"

//Executable interface around for the interpreter
type Executable interface {
	//Execute the specified operation with the current interpreter state
	Execute(i *Interpreter) error
	String() string
}

//ILoad is immediate load
type ILoad struct {
	value int64
}

func (l *ILoad) Execute(i *Interpreter) error {
	i.Push(l.value)
	return nil
}

func (l *ILoad) String() string {
	return "iconst " + strconv.FormatInt(l.value, 10)
}

func ILoadOf(value int64) *ILoad {
	return &ILoad{value}
}

//RLoad is reference load
type RLoad struct {
	ref uint32
}

func (r *RLoad) Execute(i *Interpreter) error {
	i.Push(i.current.args[r.ref])
	return nil
}

func (l *RLoad) String() string {
	return "rload " + strconv.FormatUint(uint64(l.ref), 10)
}

func RLoadOf(ref uint32) *RLoad {
	return &RLoad{ref}
}

//EmitError is the error emitted by user command and the error itself
type EmitError struct {
	ref     uint32
	message string
}

func (e *EmitError) Execute(i *Interpreter) error {
	return e
}

func (e *EmitError) Error() string {
	return e.message
}

func (e *EmitError) String() string {
	return "sigerr " + strconv.FormatUint(uint64(e.ref), 10) + ", " + e.message
}

func EmitErrorOf(ref uint32, message string) *EmitError {
	return &EmitError{ref, message}
}

//EmitWarning is the warning emitted by user command and the warning itself
type EmitWarning struct {
	ref     uint32
	message string
}

func (w *EmitWarning) Execute(i *Interpreter) error {
	return w
}

func (w *EmitWarning) Error() string {
	return w.message
}

func (w *EmitWarning) String() string {
	return "sigwrn " + strconv.FormatUint(uint64(w.ref), 10) + ", " + w.message
}

func EmitWarningOf(ref uint32, message string) *EmitWarning {
	return &EmitWarning{ref, message}
}

//BranchCode evaluate the interpreter status and optionally enqueue operations
type BranchCode struct {
	taken    []Executable
	notTaken []Executable
	cmp      int64
}

func (b *BranchCode) Execute(i *Interpreter) error {
	v := i.Pop()
	if v != b.cmp {
		i.runList = append(i.runList, b.taken...)
	} else if i.runList != nil {
		i.runList = append(i.runList, b.notTaken...)
	}
	return nil
}

func (b *BranchCode) String() string {
	return "if = " + strconv.FormatInt(b.cmp, 10)
}

func MakeBranchCode(taken []Executable, notTaken []Executable) Executable {
	return &BranchCode{taken: taken, notTaken: notTaken, cmp: 1}
}

//OutResult add output to the interpreter
type OutResult struct {
	list []Executable
}

func (o *OutResult) Execute(i *Interpreter) error {
	for _, e := range o.list {
		e.Execute(i)
		i.PushResult(i.Pop())
	}
	i.complete = true
	return nil
}

func (o *OutResult) String() string {
	return "out " + strconv.FormatInt(int64(len(o.list)), 10)
}

//MakeOutResult statement for opcodes
func MakeOutResult(list []Executable) *OutResult {
	return &OutResult{list}
}

//OutResultReverse add output to the interpreter
type OutResultReverse struct {
	list []Executable
}

func (o *OutResultReverse) Execute(i *Interpreter) error {
	vec := []int64{}
	for _, e := range o.list {
		e.Execute(i)
		vec = append(vec, i.Pop())
	}
	for k := len(vec) - 1; k >= 0; k-- {
		i.PushResult(vec[k])
	}
	i.complete = true
	return nil
}

func (o *OutResultReverse) String() string {
	return "outr " + strconv.FormatInt(int64(len(o.list)), 10)
}

//MakeOutResultReverse statement for opcodes
func MakeOutResultReverse(list []Executable) *OutResultReverse {
	return &OutResultReverse{list}
}

type RetResult struct{}

func (r *RetResult) Execute(i *Interpreter) error {
	i.PushResult(i.Pop())
	i.complete = true
	return nil
}

func (r *RetResult) String() string {
	return "ret"
}

func MakeReturn() *RetResult {
	return &RetResult{}
}

type StackExpression struct {
	list []Executable
}

func BuildStackExpression(list []Executable) *StackExpression {
	return &StackExpression{list: list}
}

func (s *StackExpression) Execute(i *Interpreter) error {
	for _, e := range s.list {
		e.Execute(i)
	}
	return nil
}

func (s *StackExpression) String() string {
	return "sexpr"
}

type Reduce struct {
	operator Operator
}

func BuildReduce(op string) *Reduce {
	ex := ExecOperator[op]
	return &Reduce{ex}
}

func (r *Reduce) Execute(i *Interpreter) error {
	r.operator(i)
	return nil
}

func (r *Reduce) String() string {
	return "binop"
}

type Negate struct{}

func BuildNegate() *Negate {
	return &Negate{}
}

func (n *Negate) Execute(i *Interpreter) error {
	i.Push(-i.Pop())
	return nil
}

func (n *Negate) String() string {
	return "neg"
}

type Complement struct{}

func BuildComplement() *Negate {
	return &Negate{}
}

func (n *Complement) Execute(i *Interpreter) error {
	i.Push(^i.Pop())
	return nil
}

func (n *Complement) String() string {
	return "inv"
}

type Not struct{}

func BuildNot() *Negate {
	return &Negate{}
}

func (n *Not) Execute(i *Interpreter) error {
	e := i.Pop()
	if e != 0 {
		i.Push(1)
	} else {
		i.Push(0)
	}
	return nil
}

func (n *Not) String() string {
	return "not"
}
