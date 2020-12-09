package main

import "github.com/aleferri/casmeleon/pkg/asm"

type SymbolTable struct {
	list []asm.Symbol
}

func (t *SymbolTable) Add(sym asm.Symbol) {
	t.list = append(t.list, sym)
}

func (t *SymbolTable) Search(name string) (asm.Symbol, bool) {
	for _, s := range t.list {
		if s.Name() == name {
			return s, true
		}
	}
	return nil, false
}

func MakeSymbolTable() SymbolTable {
	return SymbolTable{list: []asm.Symbol{}}
}
