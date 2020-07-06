package langdef

import (
	"bitbucket.org/mrpink95/casmeleon/internal/parsing"
	"bitbucket.org/mrpink95/casmeleon/internal/text"
	"bitbucket.org/mrpink95/casmeleon/internal/ui"
	"strconv"
	"strings"
)

//OpcodeNode is a node of the opcode
type OpcodeNode interface {
	Assemble(values []OpcodeArg, flags uint) ([]uint, ui.SourceCodeError)
	DebugPrint()
}

//Block of nodes delimited by braces
type Block struct {
	nodes []OpcodeNode
}

//Assemble a block of nodes
func (b *Block) Assemble(values []OpcodeArg, flags uint) ([]uint, ui.SourceCodeError) {
	var binary = []uint{}
	for _, node := range b.nodes {
		temp, err := node.Assemble(values, flags)
		if err != nil {
			return binary, err
		}
		binary = append(binary, temp...)
	}
	return binary, nil
}

//DebugPrint print in the stdout the block structure
func (b *Block) DebugPrint() {
	for _, node := range b.nodes {
		node.DebugPrint()
	}
}

//LoopFill contains parsing result of a loop statement
//Loop <index> = this_address -> <limit expression> evaluating <expression>
//es loop i until this_address + 32 do i*i/4; //this is the fourth square table of 4 bit value
//useful for the multiplication algorithm a * b = ((a + b)**2)/4 - ((a - b)**2)/4
type LoopFill struct {
	indexName text.Token   //indexName used in the loop
	limit     *Expression  //limit of index
	deposit   DepositValue //expression to evaluate
}

//Assemble fill a slice of uint with result of repeated evaluation
func (l LoopFill) Assemble(values []OpcodeArg, flags uint) ([]uint, ui.SourceCodeError) {
	if flags&1 > 0 { //suppress loops
		return nil, nil
	}
	initialSize := len(values)
	loopIndex := OpcodeArg{l.indexName.Value(), values[1].value, l.indexName}
	var emptyValue OpcodeArg
	var restoreEmptyValue bool
	var vIndex int
	if len(values) >= 3 && strings.EqualFold(values[3].name, "") {
		emptyValue = values[3]
		values[3] = loopIndex
		restoreEmptyValue = true
		vIndex = 3
	} else {
		values = append(values, loopIndex)
		vIndex = initialSize
	}
	fill := []uint{}
	limit := l.limit.Eval(values)
	byteSize := uint(1)
	for i := uint(loopIndex.value); i < limit; i += byteSize {
		fillValue, err := l.deposit.Assemble(values, flags)
		if err != nil {
			return nil, err
		}
		fill = append(fill, fillValue...)
		byteSize = uint(len(fillValue))
		loopIndex.value += len(fillValue)
		values[vIndex] = loopIndex
	}
	if restoreEmptyValue {
		values[3] = emptyValue
	}
	return fill, nil
}

//DebugPrint print the loop statement to stdout
func (l LoopFill) DebugPrint() {

}

//DepositValue deposit a value in the destination file
type DepositValue struct {
	e        *Expression
	byteSize uint
}

//Assemble deposit the value of the expression in the buffer, the high part is not cleared
func (d DepositValue) Assemble(values []OpcodeArg, flags uint) ([]uint, ui.SourceCodeError) {
	val := d.e.Eval(values)
	if d.byteSize == 1 {
		return []uint{val}, nil
	} else if d.byteSize == 2 {
		return []uint{uint(val >> 8), uint(val & 255)}, nil
	}
	return []uint{uint(val >> 24), uint(val >> 16), uint(val >> 8), uint(val & 255)}, nil
}

//DebugPrint print the deposit statement to stdout
func (d DepositValue) DebugPrint() {
	print(d.byteSize, " => ")
	d.e.DebugPrint()
}

//EmitError is the error statement
type EmitError struct {
	message string
	refer   int
}

//DebugPrint print the error statement to stdout
func (es EmitError) DebugPrint() {
	println("error statement", es.refer, es.message)
}

//Assemble return an error in the token position
func (es EmitError) Assemble(values []OpcodeArg, flags uint) ([]uint, ui.SourceCodeError) {
	if flags&(1<<2) > 0 { //suppress errors
		return nil, nil
	}
	return nil, parsing.NewParsingError(values[es.refer].token, parsing.ErrorUserDefined, es.message+"; value: "+strconv.Itoa(values[es.refer].value))
}

//BranchStatement is a standard If else statement
type BranchStatement struct {
	condition       *Expression
	conditionMet    Block
	conditionNotMet Block
}

//Assemble branch statement
func (b BranchStatement) Assemble(values []OpcodeArg, flags uint) ([]uint, ui.SourceCodeError) {
	met := b.condition.Eval(values)
	if met > 0 {
		return b.conditionMet.Assemble(values, flags)
	}
	return b.conditionNotMet.Assemble(values, flags)
}

//DebugPrint print the branch statement to stdout
func (b BranchStatement) DebugPrint() {
	println("Condition: ")
	b.condition.DebugPrint()
	println("Branches: ")
	b.conditionMet.DebugPrint()
	b.conditionNotMet.DebugPrint()
}
