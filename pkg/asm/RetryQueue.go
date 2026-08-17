package asm

import (
	"sort"

	"github.com/aleferri/casmvm/pkg/opcodes"
)

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

func (r *RetryQueue) ReAssemble(ctx Context, m opcodes.VM, imgs *[]BinaryImage) (int, error) {
	//Iteration over a map is randomized in Go: without sorting, both the order of
	//reassembly and the returned slot vary between runs on the same source
	order := make([]int, 0, len(r.addrs))
	for j := range r.addrs {
		order = append(order, j)
	}
	sort.Ints(order)

	slots := 0
	for _, j := range order {
		addr := r.addrs[j]
		compilable := r.list[j]
		newAddr, img, err := compilable.Assemble(m, addr, j, ctx)
		if err != nil {
			return 0, err
		}
		(*imgs)[j] = BinaryImage{img}

		//The caller reassembles everything from here on, so it must be the first
		//slot that moved, not whichever one came up first
		if newAddr != addr && (slots == 0 || j < slots) {
			slots = j
		}
	}
	return slots, nil
}
