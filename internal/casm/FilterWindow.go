package casm

import "errors"

type FilterWindow struct {
	matched  []string
	filtered []Opcode
}

func (win FilterWindow) FilterByFormat(format []uint32, types []uint32) FilterWindow {
	wnd := FilterWindow{}
	for _, op := range win.filtered {
		if op.Accept(format, types) {
			wnd.filtered = append(wnd.filtered, op)
		}
	}
	return wnd
}

func (win FilterWindow) PickFirst() (Opcode, error) {
	for _, op := range win.filtered {
		return op, nil
	}
	return Opcode{}, errors.New("No opcode found for provided parameters")
}

func (win FilterWindow) Candidates() []Opcode {
	return win.filtered
}
