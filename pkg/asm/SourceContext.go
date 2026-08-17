package asm

import "sort"

type SourceContext struct {
	guards   map[string]RetryQueue
	byteSize uint32
}

func MakeSourceContext(byteSize uint32) *SourceContext {
	return &SourceContext{map[string]RetryQueue{}, byteSize}
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
	names := make([]string, 0, len(ctx.guards))
	for name := range ctx.guards {
		names = append(names, name)
	}
	sort.Strings(names)

	list := []RetryQueue{}
	for _, name := range names {
		g := ctx.guards[name]
		if g.changed {
			list = append(list, g)
		}
	}
	return list
}

func (ctx *SourceContext) ByteSize() uint32 {
	return ctx.byteSize
}
