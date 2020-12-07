package exec

//Executable interface around for the interpreter
type Executable interface {
	//Execute the specified operation with the current interpreter state
	Execute(i *Interpreter) error
}

//ILoad is immediate load
type ILoad struct {
	value int64
}

func (l *ILoad) Execute(i *Interpreter) error {
	i.Push(l.value)
	return nil
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

func MakeBranchCode(taken []Executable, notTaken []Executable) Executable {
	return &BranchCode{taken: taken, notTaken: notTaken, cmp: 0}
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

//MakeOutResult statement for opcodes
func MakeOutResult(list []Executable) *OutResult {
	return &OutResult{list}
}

type RetResult struct{}

func (r *RetResult) Execute(i *Interpreter) error {
	i.PushResult(i.Pop())
	i.complete = true
	return nil
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

type Negate struct {
	embed Executable
}

func BuildNegate(embed Executable) *Negate {
	return &Negate{embed}
}

func (n *Negate) Execute(i *Interpreter) error {
	err := n.embed.Execute(i)
	i.Push(-i.Pop())
	return err
}

type Complement struct {
	embed Executable
}

func BuildComplement(embed Executable) *Negate {
	return &Negate{embed}
}

func (n *Complement) Execute(i *Interpreter) error {
	err := n.embed.Execute(i)
	i.Push(^i.Pop())
	return err
}

type Not struct {
	embed Executable
}

func BuildNot(embed Executable) *Negate {
	return &Negate{embed}
}

func (n *Not) Execute(i *Interpreter) error {
	err := n.embed.Execute(i)
	e := i.Pop()
	if e != 0 {
		i.Push(1)
	} else {
		i.Push(0)
	}
	return err
}
