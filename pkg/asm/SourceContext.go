package asm

type SourceContext struct {
	guards map[Symbol]RetryQueue
}

func (ctx *SourceContext) GuardSymbol(sym Symbol, x int, addr uint32, c Compilable) {
	retry, ok := ctx.guards[sym]
	if !ok {
		retry = MakeRetryQueue()
	}
	retry.Append(x, addr, c)
	ctx.guards[sym] = retry
}

func (ctx *SourceContext) ClearAll() {
	for sym := range ctx.guards {
		retry := ctx.guards[sym]
		retry.changed = false
		ctx.guards[sym] = retry
	}
}

func (ctx *SourceContext) NotifyChange(sym Symbol) {
	retry := ctx.guards[sym]
	retry.changed = true
	ctx.guards[sym] = retry
}

func (ctx *SourceContext) RetryList() []RetryQueue {
	list := []RetryQueue{}
	for _, g := range ctx.guards {
		if g.changed {
			list = append(list, g)
		}
	}
	return list
}
