package asm

type BinaryImage struct {
	content []uint8
}

type Slot struct {
	index int
	addr  uint32
}

func AssembleSource(list []Compilable, ctx Context) ([]uint8, error) {
	addr := uint32(0)
	var err error
	var img []uint8

	fixLater := []Slot{}

	result := []BinaryImage{}

	for j, a := range list {
		init := addr
		addr, img, err = a.Assemble(addr, j, ctx)

		if !a.IsAddressInvariant() {
			fixLater = append(fixLater, Slot{j, init})
		}

		if err != nil {
			return img, err
		}

		result = append(result, BinaryImage{img})
	}

	retry := ctx.RetryList()

	for len(retry) != 0 {
		ctx.ClearAll()

		for _, r := range retry {
			fix, err := r.ReAssemble(ctx, &result)
			if err != nil {
				return img, err
			}
			if fix != 0 {
				for _, v := range fixLater {
					addr, img, err = list[v.index].Assemble(addr, v.index, ctx)
				}
			}
		}
	}

	img = []uint8{}
	for _, bin := range result {
		img = append(img, bin.content...)
	}
	return img, nil
}
