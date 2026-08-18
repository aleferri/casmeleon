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

//ReAssemble is no longer used by AssembleSource: reassembling a subset of the
//slots at the addresses recorded by a previous pass is what made stale label
//addresses survive. Kept for API compatibility, safe to drop.
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

//Slots returns the indexes registered on this queue, sorted. AssembleSource
//uses them to know which items depend on a symbol that moved, so that a pass
//can reassemble those and reuse the bytes of everything else.
func (r *RetryQueue) Slots() []int {
	order := make([]int, 0, len(r.list))
	for j := range r.list {
		order = append(order, j)
	}
	sort.Ints(order)
	return order
}
