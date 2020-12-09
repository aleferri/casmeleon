package casm

import "errors"

type FilterWindow struct {
	matched    []string
	filtered   []Opcode
	lenMatches []bool
}

func (win FilterWindow) FilterByFormat(format []uint32) FilterWindow {
	wnd := FilterWindow{}
	for _, op := range win.filtered {
		valid := true
		for k, kind := range format {
			if k >= len(op.format) || op.format[k] != kind {
				valid = false
			}
		}
		if valid {
			wnd.filtered = append(wnd.filtered, op)
			wnd.lenMatches = append(wnd.lenMatches, len(format) == len(op.format))
		}
	}
	return wnd
}

func (win FilterWindow) PickFirst() (Opcode, error) {
	for r, op := range win.filtered {
		if win.lenMatches[r] {
			return op, nil
		}
	}
	return Opcode{}, errors.New("No opcode found for provided parameters")
}
