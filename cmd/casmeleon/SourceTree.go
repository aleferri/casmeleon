package main

import (
	"bitbucket.org/mrpink95/casmeleon/internal/langdef"
	"bitbucket.org/mrpink95/casmeleon/internal/parsing"
	"bitbucket.org/mrpink95/casmeleon/internal/text"
	"bitbucket.org/mrpink95/casmeleon/internal/ui"
)

//ASMSource is an interface to every type of source like not-yet-included source or opcode list
type ASMSource interface {
	Parent() ASMSource                           //GetParent return the parent source
	Size() int                                   //GetSize return the number of opcodes in this source
	AddChild(ASMSource)                          //NewChild create new child
	AddOpcode(opcode OpcodeInfo)                 //AddOpcode add opcode to the source
	AddLabel(label ASMLabel) LabelInfo           //AddLabel add a label to the source tree
	GetChildIndex(child ASMSource) int           //GetChildIndex return the index of the child
	NextOpcodes(i int) (ASMSource, []OpcodeInfo) //NextOpcodes return the list of opcode in the next source part
	NextLabels(i int) (ASMSource, []ASMLabel)    //NextLabels return the list of labels in the next source part
	FilePath() string                            //FileName return the file name of the source
}

//ASMSourceRoot is the first file
type ASMSourceRoot struct {
	filePath string
	sources  []ASMSource
	lang     *langdef.LangDef //lang is the language of the program
}

//NewASMSourceRoot return a new root source
func NewASMSourceRoot(filePath string, lang *langdef.LangDef) *ASMSourceRoot {
	return &ASMSourceRoot{filePath, []ASMSource{}, lang}
}

//FilePath return the source file path
func (p *ASMSourceRoot) FilePath() string {
	return p.filePath
}

//Parent return the parent of the source
func (p *ASMSourceRoot) Parent() ASMSource {
	return nil //no parent for the root node
}

//Size return size of the opcode source
func (p *ASMSourceRoot) Size() int {
	sourceSize := len(p.sources)
	if sourceSize == 0 {
		return 0
	}
	size := 0
	for _, source := range p.sources {
		size += source.Size()
	}
	return size
}

//AddChild add a child to the root
func (p *ASMSourceRoot) AddChild(source ASMSource) {
	p.sources = append(p.sources, source)
}

//AddOpcode insert an opcode in the last source tree
func (p *ASMSourceRoot) AddOpcode(opcode OpcodeInfo) {
	p.sources[len(p.sources)].AddOpcode(opcode)
}

//AddLabel insert a label in the last source tree
func (p *ASMSourceRoot) AddLabel(label ASMLabel) LabelInfo {
	return p.sources[len(p.sources)].AddLabel(label)
}

//GetChildIndex return the index of the child node
func (p *ASMSourceRoot) GetChildIndex(child ASMSource) int {
	for i, s := range p.sources {
		if s == child {
			return i
		}
	}
	return -1
}

//NextOpcodes return nil only if there are no more opcodes
func (p *ASMSourceRoot) NextOpcodes(i int) (ASMSource, []OpcodeInfo) {
	if i >= len(p.sources) {
		return p, nil
	}
	return p.sources[i].NextOpcodes(0)
}

//NextLabels return nil only if there are no more opcodes
func (p *ASMSourceRoot) NextLabels(i int) (ASMSource, []ASMLabel) {
	if i >= len(p.sources) {
		return p, nil
	}
	return p.sources[i].NextLabels(0)
}

//ASMSourceInclude is a node in the source tree, it loads source from "include" statements
type ASMSourceInclude struct {
	parent   ASMSource        //parent include file
	sources  []ASMSource      //list of asm source
	filePath text.Token       //filePath is the path of the included file
	textUI   ui.UI            //ui where to print errors
	lang     *langdef.LangDef //lang is the language of the program
}

//SetupBuilder setup an already used builder
func (p *ASMSourceInclude) SetupBuilder(oldBuilder ASMReader, settings Settings) ASMReader {
	buffer := parsing.NewTokenBufferFromFile(p.filePath.Value(), ASMDefaultOptions(settings.lang))
	buffer.SyncLines(GetSyncOfFile(p.filePath.Value()))
	child := NewASMSourceLeaf(p)
	p.AddChild(child)
	pState := NewParserState(oldBuilder.pState.pUI, buffer)
	pState.SetCustomIdentification(func(t text.Token) (text.Token, ui.SourceCodeError) { return settings.lang.IdentifyNumber(t) })
	builder := NewASMReader(pState, child)
	builder.pass = oldBuilder.pass
	return builder
}

//FilePath return the source file path
func (p *ASMSourceInclude) FilePath() string {
	return p.filePath.Value()
}

//Size return size of the opcode source
func (p *ASMSourceInclude) Size() int {
	sourceSize := len(p.sources)
	if sourceSize == 0 {
		return 0
	}
	size := 0
	for _, source := range p.sources {
		size += source.Size()
	}
	return size
}

//AddChild add a child to the node
func (p *ASMSourceInclude) AddChild(source ASMSource) {
	p.sources = append(p.sources, source)
}

//AddOpcode add an opcode
func (p *ASMSourceInclude) AddOpcode(opcode OpcodeInfo) {
	p.sources[len(p.sources)].AddOpcode(opcode)
}

//AddLabel insert a label in the last source tree
func (p *ASMSourceInclude) AddLabel(label ASMLabel) LabelInfo {
	return p.sources[len(p.sources)].AddLabel(label)
}

//Parent return the parent of the source
func (p *ASMSourceInclude) Parent() ASMSource {
	return p.parent
}

//GetChildIndex return the index of the child node
func (p *ASMSourceInclude) GetChildIndex(child ASMSource) int {
	for i, s := range p.sources {
		if s == child {
			return i
		}
	}
	return -1
}

//NextOpcodes return the next child of parent or his own child
func (p *ASMSourceInclude) NextOpcodes(i int) (ASMSource, []OpcodeInfo) {
	if i >= len(p.sources) {
		return p.parent.NextOpcodes(p.parent.GetChildIndex(p) + 1)
	}
	return p.sources[i].NextOpcodes(0)
}

//NextLabels return the next child of parent or his own child
func (p *ASMSourceInclude) NextLabels(i int) (ASMSource, []ASMLabel) {
	if i >= len(p.sources) {
		return p.parent.NextLabels(p.parent.GetChildIndex(p) + 1)
	}
	return p.sources[i].NextLabels(0)
}

//ASMSourceLeaf contains the opcode list for the source
type ASMSourceLeaf struct {
	parent       ASMSource
	instructions []OpcodeInfo
	labels       []ASMLabel
}

//NewASMSourceLeaf return a child node of the parent tree
func NewASMSourceLeaf(parent ASMSource) *ASMSourceLeaf {
	return &ASMSourceLeaf{parent, nil, nil}
}

//FilePath return the source file path
func (p *ASMSourceLeaf) FilePath() string {
	return p.parent.FilePath()
}

//Size return size of the opcode source
func (p *ASMSourceLeaf) Size() int {
	return len(p.instructions)
}

//AddChild is invalid operation
func (p *ASMSourceLeaf) AddChild(source ASMSource) {
	//invalid operation
}

//AddOpcode add an opcode
func (p *ASMSourceLeaf) AddOpcode(opcode OpcodeInfo) {
	p.instructions = append(p.instructions, opcode)
}

//AddLabel add a label
func (p *ASMSourceLeaf) AddLabel(label ASMLabel) LabelInfo {
	p.labels = append(p.labels, label)
	return LabelInfo{&p.labels[len(p.labels)-1], 0}
}

//Parent return the parent of the source
func (p *ASMSourceLeaf) Parent() ASMSource {
	return p.parent
}

//GetChildIndex return the index of the child node
func (p *ASMSourceLeaf) GetChildIndex(child ASMSource) int {
	return -1
}

//NextOpcodes always return the child contents and the parent pointer
func (p *ASMSourceLeaf) NextOpcodes(i int) (ASMSource, []OpcodeInfo) {
	return p.parent, p.instructions
}

//NextLabels always return the child contents and the parent pointer
func (p *ASMSourceLeaf) NextLabels(i int) (ASMSource, []ASMLabel) {
	return p.parent, p.labels
}

//ASMSourceIterator iterate over the source tree
type ASMSourceIterator struct {
	source      ASMSource
	nextSection int
}

//Next return the next opcode
func (p *ASMSourceIterator) Next() []OpcodeInfo {
	source, opcode := p.source.NextOpcodes(p.nextSection)
	p.source = source
	p.nextSection++
	return opcode
}

//Labels return the associated labels
func (p *ASMSourceIterator) Labels() []ASMLabel {
	_, labels := p.source.NextLabels(p.nextSection - 1)
	return labels
}
