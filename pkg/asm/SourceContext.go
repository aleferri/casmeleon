package asm

type SourceContext struct {
	guards map[string]RetryQueue
}

func MakeSourceContext() *SourceContext {
	return &SourceContext{map[string]RetryQueue{}}
}

func (ctx *SourceContext) EnsureExists(name string) RetryQueue {
	retry, ok := ctx.guards[name]
	if !ok {
		retry = MakeRetryQueue()
		ctx.guards[name] = retry
	}
	return retry
}

func (ctx *SourceContext) GuardSymbol(name string, x int, addr uint32, c Compilable) {
	retry := ctx.EnsureExists(name)
	retry.Append(x, addr, c)
	ctx.guards[name] = retry
}

func (ctx *SourceContext) ClearAll() {
	for sym := range ctx.guards {
		retry := ctx.guards[sym]
		retry.changed = false
		ctx.guards[sym] = retry
	}
}

func (ctx *SourceContext) Refresh(sym Symbol) {
	retry := ctx.EnsureExists(sym.Name())
	retry.changed = true
	ctx.guards[sym.Name()] = retry
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
