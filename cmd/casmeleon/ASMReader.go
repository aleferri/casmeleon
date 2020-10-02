package main

import (
	"github.com/aleferri/casmeleon/internal/langdef"
	"github.com/aleferri/casmeleon/internal/parsing"
	"github.com/aleferri/casmeleon/internal/text"
	"path/filepath"
	"strings"
)

//ASMReader is a builder for the source tree
type ASMReader struct {
	program         ASMSource
	pState          *ASMParserState
	unresolved      []text.Token
	lastGlobalLabel *ASMLabel
	pass            PassResult
}

//NewASMReader return a new builder for the program
func NewASMReader(pState *ASMParserState, child ASMSource) ASMReader {
	pass := PassResult{make(map[string]LabelInfo), make([]uint, 1), make([]int, 1)}
	return ASMReader{child, pState, []text.Token{}, nil, pass}
}

//NewASMLabel create a new asm label with the name and the current information preventing inconsistent state in the label
//i.e. a label that is local but has no reference, a label having reference to a local label, etc
func (p *ASMReader) NewASMLabel(name string) ASMLabel {
	label := NewASMLabel(name, p.program.Size()-1, p)
	if label.global {
		p.lastGlobalLabel = &label
	}
	return label
}

//NewOpcodeAssemble create a new opcode assemble structure, it also insert this_address value in the arguments
func (p *ASMReader) NewOpcodeAssemble(name text.Token, opcodeIndex int, argsName []text.Token, lang *langdef.LangDef) {
	thisOpcode := name
	thisAddress := text.NewSpecialToken(name, "this_address", text.Number)
	argsValue, argsLabel := solveNames(argsName, lang, p)
	argsName = append(argsName, thisAddress, thisOpcode)
	argsValue = append(argsValue, 0, 0)
	argsLabel = append(argsLabel, false, false)
	opcode := NewOpcodeInfo(opcodeIndex, argsName, argsValue, lang)
	opcode.argsLabel = argsLabel
	p.program.AddOpcode(opcode)
}

//AddLabel parse a single label
func (p *ASMReader) AddLabel(name text.Token, lang *langdef.LangDef) bool {
	label := p.NewASMLabel(name.Value())
	if p.pass.FindLabelAddress(label.FullName()) != -1 {
		p.pState.ReportError(name, parsing.ErrorDuplicatedLabel, "label already in use, use local labels if you want to reuse the name")
		return false
	}
	lInfo := p.program.AddLabel(label)
	p.pass.labels[label.FullName()] = lInfo
	p.RemoveUnresolvedSymbol(label.FullName())
	return true
}

//AddInclude add a child in the source tree that will be deferred until the parsing of the current file end
func (p *ASMReader) AddInclude(tokens []text.Token, lang *langdef.LangDef, includeList *[]ASMSourceInclude) bool {
	filePath := tokens[1]
	absPath, err := filepath.Abs(filePath.Value()[1 : len(filePath.Value())-1])
	node := p.program
	noErr := err == nil
	for noErr && node != nil {
		if strings.EqualFold(node.FilePath(), absPath) {
			p.pState.ReportError(filePath, parsing.ErrorCyclicInclusion, "cyclic include found")
			noErr = false
		}
		node = node.Parent()
	}
	if noErr {
		parent := p.program.Parent()
		include := ASMSourceInclude{parent, []ASMSource{}, filePath.WithValue(absPath), p.pState.pUI, lang}
		*includeList = append(*includeList, include)
		parent.AddChild(&include)
		child := NewASMSourceLeaf(parent)
		parent.AddChild(child)
	}
	return noErr
}

//Parse parse the input file. Who would have guessed?
func (p *ASMReader) Parse(lang *langdef.LangDef, includeList *[]ASMSourceInclude) bool {
	includeParser := matchInclude()
	opcodeParser := matchOpcode(lang)
	labelParser := matchLabel()
	tokens := []text.Token{}
	noErr := true
	tokens, notEof := skipEOL([]text.Token{}, p.pState)
	for notEof {
		valid := true
		p.pState.Pull(&tokens)
		if tokens[0].EnumType() == text.KeywordInclude {
			tokens, valid = includeParser(tokens, p.pState)
			valid = valid && p.AddInclude(p.pState.lastParsed, lang, includeList)
		} else {
			isLabel := false
			tokens, isLabel = isASMLabel(tokens, p.pState)
			if isLabel {
				tokens, valid = labelParser(tokens, p.pState)
				valid = valid && p.AddLabel(p.pState.lastParsed[0], lang)
			} else {
				tokens, valid = opcodeParser(tokens, p.pState)
				args := []text.Token{}
				if valid && len(p.pState.lastParsed) > 1 {
					args = p.pState.lastParsed[1:]
				}
				if valid {
					opcodeIndex := p.pState.lastOpcode
					p.NewOpcodeAssemble(p.pState.lastParsed[0], opcodeIndex, args, lang)
				}
			}
		}
		noErr = noErr && valid
		p.pState.lastParsed = []text.Token{}
		tokens, notEof = skipEOL([]text.Token{}, p.pState)
	}
	return noErr
}

//ExpandName expand a label name to the full name
func (p *ASMReader) ExpandName(name string) string {
	if strings.HasPrefix(name, ".") {
		return p.lastGlobalLabel.name + name
	}
	return name
}

//AddUnresolvedSymbol add a currently unresolved symbol to the array
//of unresolved symbol
func (p *ASMReader) AddUnresolvedSymbol(symbol text.Token) {
	p.unresolved = append(p.unresolved, symbol)
}

//RemoveUnresolvedSymbol remove a symbol not known before the current statement
func (p *ASMReader) RemoveUnresolvedSymbol(name string) {
	for i, n := range p.unresolved {
		if strings.EqualFold(n.Value(), name) {
			p.unresolved[i] = text.NewInternalToken("")
		}
	}
}

//ReportUnresolvedSymbols report to the user symbols not found in the code
func (p *ASMReader) ReportUnresolvedSymbols() {
	for _, n := range p.unresolved {
		if len(n.Value()) > 0 {
			p.pState.ReportError(n, parsing.ErrorUnresolvedSymbol, "")
		}
	}
}
