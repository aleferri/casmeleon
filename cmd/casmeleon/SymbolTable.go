package main

import (
	"github.com/aleferri/casmeleon/pkg/asm"
	"github.com/aleferri/casmeleon/pkg/text"
)

type SymbolTable struct {
	list            []asm.Symbol
	lastGlobalLabel *asm.Label
	watchList       []text.Symbol
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

func (t *SymbolTable) Watch(token text.Symbol) {
	t.watchList = append(t.watchList, token)
}

func (t *SymbolTable) UnWatch(name string) {
	rerun := []text.Symbol{}
	for _, w := range t.watchList {
		if w.Value() != name {
			rerun = append(rerun, w)
		}
	}
	t.watchList = rerun
}

func MakeSymbolTable() SymbolTable {
	return SymbolTable{list: []asm.Symbol{}, lastGlobalLabel: nil, watchList: []text.Symbol{}}
}
