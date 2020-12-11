package main

import "github.com/aleferri/casmeleon/pkg/asm"

type SelfPatchSymbol struct {
	fqn      string
	sym      asm.Symbol
	patched  bool
	symTable *SymbolTable
}

func MakePatchSymbol(fqn string, symTable *SymbolTable) *SelfPatchSymbol {
	return &SelfPatchSymbol{fqn: fqn, sym: nil, patched: false, symTable: symTable}
}

func (p *SelfPatchSymbol) Address() uint32 {
	if p.patched {
		return p.sym.Address()
	}
	return 0
}

func (p *SelfPatchSymbol) Value() int64 {
	if !p.patched {
		p.sym, _ = p.symTable.Search(p.fqn)
		p.patched = true
	}
	return p.sym.Value()
}

func (p *SelfPatchSymbol) Name() string {
	return p.fqn
}

func (p *SelfPatchSymbol) IsDynamic() bool {
	if p.patched {
		return p.sym.IsDynamic()
	}
	return true
}
