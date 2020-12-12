package exec

//Interpreter state
type Interpreter struct {
	current  Frame
	runList  []Executable
	complete bool
}

//Make interpreter to run code
func MakeInterpreter(frame Frame, runList []Executable) Interpreter {
	return Interpreter{current: frame, runList: runList}
}

//Run the interprter for the specified parameters
func (i *Interpreter) Run() error {
	i.complete = false
	var err error = nil
	for len(i.runList) > 0 && err == nil && !i.complete {
		q := i.runList[0]
		i.runList = i.runList[1:]
		err = q.Execute(i)
	}
	return err
}

//CallFrame replace the current interpreter frame with a new frame giving a new list of executable statements
func (i *Interpreter) CallFrame(f Frame, list []Executable) error {
	last := i.current
	queue := i.runList
	backComplete := i.complete

	i.current = f
	i.runList = list

	err := i.Run()

	last.eval.Push(i.current.ret.Pop())

	i.complete = backComplete
	i.current = last
	i.runList = queue
	return err
}

//Push value into the stack
func (i *Interpreter) Push(v int64) {
	i.current.eval.Push(v)
}

//Pop value into the stack
func (i *Interpreter) Pop() int64 {
	return i.current.eval.Pop()
}

//PushResult to the frame
func (i *Interpreter) PushResult(v int64) {
	i.current.ret.Push(v)
}

func (i *Interpreter) PopResults() *Stack {
	return i.current.ret
}
