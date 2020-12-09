package asm

type RetryQueue struct {
	list    map[int]Compilable
	addrs   map[int]uint32
	changed bool
}

func MakeRetryQueue() RetryQueue {
	return RetryQueue{list: map[int]Compilable{}, addrs: map[int]uint32{}}
}

func (r *RetryQueue) Append(j int, addr uint32, c Compilable) {
	r.list[j] = c
	r.addrs[j] = addr
}

func (r *RetryQueue) ReAssemble(ctx Context, imgs *[]BinaryImage) (int, error) {
	slots := 0
	for j, addr := range r.addrs {
		compilable := r.list[j]
		newAddr, img, err := compilable.Assemble(addr, j, ctx)
		if err != nil {
			return 0, err
		}
		(*imgs)[j] = BinaryImage{img}

		if newAddr != addr && slots == 0 {
			slots = j
		}
	}
	return slots, nil
}
