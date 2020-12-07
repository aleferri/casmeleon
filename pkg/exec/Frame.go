package exec

//Frame for function call, i tried to avoid it, but it get messy without it
type Frame struct {
	args []int64
	eval Stack
	ret  Stack
}

//FrameOf
func FrameOf(args []int64) Frame {
	return Frame{args: args, eval: EmptyStack(), ret: EmptyStack()}
}
