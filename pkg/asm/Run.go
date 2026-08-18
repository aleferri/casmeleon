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
//the position in the output moves, and that follows from the concatenation.
//
//What must not happen is reassembling a subset of the items at the addresses
//recorded during an earlier pass. Label.Assemble writes the address it receives
//into the label, so feeding it a stale address does not fail to update the
//label, it actively overwrites the correct value with the old one, and every
//caller of that label ends up pointing before the real target.
func AssembleSource(m opcodes.VM, list []Compilable, ctx Context) ([]uint8, error) {
	result := make([]BinaryImage, len(list))
	lastAddr := make([]uint32, len(list))
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
				lastAddr[j] = here
				addr = here + uint32(len(result[j].content))
				continue
			}

			next, img, err := item.Assemble(m, here, j, ctx)
			if err != nil {
				return img, err
			}
			result[j] = BinaryImage{img}
			lastAddr[j] = here
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

		img := []uint8{}
		for _, bin := range result {
			img = append(img, bin.content...)
		}
		return img, nil
	}

	return nil, errors.New("addresses did not stabilize in " + fmt.Sprint(maxPasses) + " passes")
}
