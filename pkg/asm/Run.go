package asm

import (
	"errors"
	"fmt"

	"github.com/aleferri/casmvm/pkg/opcodes"
)

type BinaryImage struct {
	content []uint8
}

type Slot struct {
	index int
	addr  uint32
}

func AssembleSource(m opcodes.VM, list []Compilable, ctx Context) ([]uint8, error) {
	addr := uint32(0)
	var err error
	var img []uint8

	fixLater := []Slot{}

	result := []BinaryImage{}

	for j, a := range list {
		init := addr
		addr, img, err = a.Assemble(m, addr, j, ctx)

		if !a.IsAddressInvariant() {
			fixLater = append(fixLater, Slot{j, init})
		}

		if err != nil {
			return img, err
		}

		result = append(result, BinaryImage{img})
	}

	fmt.Println("First pass done, checking things that must rerun")

	retry := ctx.RetryList()
	passes := 0

	for len(retry) != 0 && passes < 2 {
		ctx.ClearAll()

		fmt.Println("Ouput oscillation, pending ", len(retry), "opcode rebuild")

		for _, r := range retry {
			fix, err := r.ReAssemble(ctx, m, &result)
			if err != nil {
				return img, err
			}
			if fix != 0 {
				for _, v := range fixLater {
					_, img, err = list[v.index].Assemble(m, v.addr, v.index, ctx)
					result[v.index] = BinaryImage{img}
				}
			}
		}

		retry = ctx.RetryList()
		passes++
	}

	fmt.Println("Ouput stabilized, done")

	img = []uint8{}
	for _, bin := range result {
		img = append(img, bin.content...)
	}

	if passes < 2 {
		return img, nil
	}
	return img, errors.New("Infinite loop was stopped")
}
