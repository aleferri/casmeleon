package main

import (
	"github.com/aleferri/casmeleon/internal/langdef"
	"github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/internal/ui"
	"strconv"
)

//LabelInfo contains a pointer to the static part and the dynamic part of a label
type LabelInfo struct {
	ptr     *ASMLabel //pointer to the static label part
	address int       //last known address
}

//PassResult is the output of the current pass over the source
type PassResult struct {
	labels    map[string]LabelInfo //labels is a map: full name of label -> labels info
	output    []uint               //output of the pass
	outputLen []int                //len of output of individual opcode
}

func (p *PassResult) updateLabelAddress(fullName string, address int) {
	label := p.labels[fullName]
	label.address = address
	p.labels[fullName] = label
}

//makePass make another pass on the source
func (p *PassResult) makePass(oldPass *PassResult, source ASMSource, flags uint) (uint, ui.SourceCodeError) {
	lastAddress, nextLabel := 0, 0
	iterator := ASMSourceIterator{source, 0}
	opcodes := iterator.Next()
	labels := iterator.Labels()
	diff := uint(0)
	for opcodes != nil {
		for len(labels) > nextLabel && labels[nextLabel].index == -1 {
			p.updateLabelAddress(labels[nextLabel].FullName(), lastAddress)
			nextLabel++
		}
		for i, o := range opcodes {
			oldPass.UpdateLabels(&o)
			argsLen := len(o.argsValue)
			o.argsValue[argsLen-1] = 0
			o.argsValue[argsLen-2] = lastAddress
			val, err := o.Assemble(flags)
			if err == nil {
				for len(labels) > nextLabel && labels[nextLabel].index == i {
					p.updateLabelAddress(labels[nextLabel].FullName(), lastAddress)
					nextLabel++
				}
				p.output = append(p.output, val...)
				lastAddress += len(val)
				p.outputLen = append(p.outputLen, len(val))
				if len(oldPass.outputLen) <= i {
					diff += uint(len(val))
				} else {
					diff += uint(len(val) - oldPass.outputLen[i])
				}
			} else {
				return 1, err
			}
		}
		opcodes = iterator.Next()
		labels = iterator.Labels()
		nextLabel = 0
	}
	return diff, nil
}

//FindLabelAddress find the address of a label if exist
//if not return -1
func (p *PassResult) FindLabelAddress(fullName string) int {
	lInfo, exist := p.labels[fullName]
	if !exist {
		return -1
	}
	return lInfo.address
}

//UpdateLabels in the opcode, require tokens with full name
//so be sure to update token content with the label full name
func (p *PassResult) UpdateLabels(opcode *OpcodeInfo) {
	for i, v := range opcode.argsName {
		if opcode.argsLabel[i] {
			opcode.argsValue[i] = p.FindLabelAddress(v.Value())
		}
	}
}

//PrintResult of pass
func (p *PassResult) PrintResult(source ASMSource) {
	iterator := ASMSourceIterator{source, 0}
	opcodes := iterator.Next()
	outputIndex := 0
	for opcodes != nil {
		for i, o := range opcodes {
			print(outputIndex, ": ")
			print(o.String(), " -> ")
			for k := 0; k < p.outputLen[i]; k++ {
				print(p.output[outputIndex+k], " ")
			}
			println()
			outputIndex += p.outputLen[i]
		}
		opcodes = iterator.Next()
	}
	println("Labels:")
	for _, label := range p.labels {
		println(label.ptr.FullName(), "@", label.address)
	}
}

//solveNames solve names found in arguments
func solveNames(names []text.Token, lang *langdef.LangDef, a *ASMReader) ([]int, []bool) {
	var values []int
	var isLabel []bool
	for i, v := range names {
		n := 0
		label := false
		if v.EnumType() == text.Number {
			n, _ = strconv.Atoi(v.Value())
		} else if v.EnumType() == text.DoubleQuotedString {
			n = 0
		} else if lang.IsEnumValue(v.Value()) {
			n = lang.GetEnumValue(v.Value())
		} else {
			label = true
			fullName := a.ExpandName(v.Value())
			names[i] = names[i].WithValue(fullName)
			k, exist := a.pass.labels[fullName]
			if exist {
				n = k.address
			} else {
				a.AddUnresolvedSymbol(text.NewSpecialToken(v, fullName, text.Identifier))
			}
		}
		isLabel = append(isLabel, label)
		values = append(values, n)
	}
	return values, isLabel
}
