package asm

import (
	"errors"
	"fmt"

	"github.com/aleferri/casmvm/pkg/opcodes"
)

type BinaryImage struct {
	content []uint8
}

//maxPasses caps the fixed point search: with automatic short/long form
//selection a pathological source can oscillate instead of converging
const maxPasses = 10

//AssembleSource assembles the list until the addresses stop moving.
//
//An item is reassembled only when it can actually produce different bytes,
//that is when it depends on a symbol that moved during the previous pass, or
//when it depends on its own address and that address changed. Everything else
//keeps the bytes of the previous pass, because they cannot have changed: only
//the position in the output moves, and that follows from the address each item
//was assembled at.
//
//The next address is remembered rather than recomputed from the byte count. The
//two are equal for everything that emits what it spans, but not for .org, which
//moves the address and emits nothing: recomputing made a skipped .org into a
//silent no-op.
//
//What must not happen is reassembling a subset of the items at the addresses
//recorded during an earlier pass. Label.Assemble writes the address it receives
//into the label, so feeding it a stale address does not fail to update the
//label, it actively overwrites the correct value with the old one, and every
//caller of that label ends up pointing before the real target.
func AssembleSource(m opcodes.VM, list []Compilable, ctx Context) ([]uint8, error) {
	result := make([]BinaryImage, len(list))
	lastAddr := make([]uint32, len(list))
	nextAddr := make([]uint32, len(list))
	known := make([]bool, len(list))

	for pass := 0; pass < maxPasses; pass++ {
		//the slots to redo are those registered on symbols that moved during
		//the previous pass: the retry queues already know which ones
		dirty := map[int]bool{}
		for _, queue := range ctx.RetryList() {
			for _, j := range queue.Slots() {
				dirty[j] = true
			}
		}
		ctx.ClearAll()

		work := 0
		addr := uint32(0)

		for j, item := range list {
			here := addr

			//the bytes cannot have changed when no symbol this item depends on
			//has moved and either its address is irrelevant or it did not move
			if known[j] && !dirty[j] && (here == lastAddr[j] || item.IsAddressInvariant()) {
				addr = nextAddr[j] + (here - lastAddr[j])
				lastAddr[j] = here
				continue
			}

			next, img, err := item.Assemble(m, here, j, ctx)
			if err != nil {
				return img, err
			}
			result[j] = BinaryImage{img}
			lastAddr[j] = here
			nextAddr[j] = next
			known[j] = true
			addr = next
			work++
		}

		if work != 0 {
			continue
		}

		//nothing was reassembled, so nothing can have changed: fixed point
		for j, item := range list {
			fmt.Println("@", lastAddr[j]+uint32(len(result[j].content)), ":", item, "->", result[j].content)
		}
		fmt.Println("Addresses stable, done in", pass+1, "pass(es)")

		//Placed by address, not by concatenation: an item's bytes belong at the
		//address it was assembled at, which is the only reading under which
		//.org means anything and under which two of them may run backwards.
		size := uint32(0)
		for j := range result {
			if end := lastAddr[j] + uint32(len(result[j].content)); end > size {
				size = end
			}
		}
		img := make([]uint8, size)
		for j, bin := range result {
			copy(img[lastAddr[j]:], bin.content)
		}
		return img, nil
	}

	return nil, errors.New("addresses did not stabilize in " + fmt.Sprint(maxPasses) + " passes")
}
