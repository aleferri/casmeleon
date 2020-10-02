package langdef

import (
	"github.com/aleferri/casmeleon/internal/parsing"
	"github.com/aleferri/casmeleon/internal/text"
	"sort"
	"strings"
)

//LangBuilder is the builder for the Lang struct
type LangBuilder struct {
	enums         []Enum
	numberFormats []NumberFormat
	opcodes       []Opcode
	underParsing  Opcode
}

//NewLangBuilder create a new LangBuilder
func NewLangBuilder() LangBuilder {
	return LangBuilder{nil, nil, nil, Opcode{}}
}

//Build LangDef object
func (builder *LangBuilder) Build(pState *LangParserState) *LangDef {
	parser := createLangParser()
	parser([]text.Token{}, pState)
	for _, node := range pState.cst {
		var noErr bool
		if node.tag == TagOpcode {
			noErr = builder.ParseOpcodeDeclaration(node, pState.report)
		} else if node.tag == TagNumberFormat {
			noErr = builder.ParseNumericFormat(node)
		} else if node.tag == TagEnumDeclaration {
			noErr = builder.ParseEnum(node)
		} else {
			noErr = false
		}
		if !noErr {
			break
		}
	}
	return &LangDef{builder.enums, builder.numberFormats, builder.opcodes}
}

//ParseNumericFormat parse the number format for the specified base
//always return true
func (builder *LangBuilder) ParseNumericFormat(node CSTNode) bool {
	var known = []string{".bin", ".dec", ".hex", ".oct"}
	var baseNumber = []int{2, 10, 16, 8}
	var index = sort.SearchStrings(known, node.content[1].Value())
	var format = NewNumberFormat(node.content[3].Value(), strings.EqualFold(node.content[2].Value(), "prefix"), baseNumber[index])
	builder.numberFormats = append(builder.numberFormats, format)
	return true
}

//ParseEnum parse the enums in the language definition
func (builder *LangBuilder) ParseEnum(node CSTNode) bool {
	name := node.content[1]
	values := []string{}
	for i := 2; i < len(node.content); i++ {
		values = append(values, node.content[i].Value())
	}
	builder.enums = append(builder.enums, Enum{name.Value(), values})
	return true
}

//ParseLoop parse loop in form: loop <identifier> until <limit expression> <deposit expression> ;
func (builder *LangBuilder) ParseLoop(node CSTNode) (*LoopFill, *parsing.ParsingError) {
	param := node.content[1]
	paramsName := builder.underParsing.paramsName
	limit, limitErr := NewExpression(node.children[0].content, paramsName, builder.enums)
	oldParamsName := builder.underParsing.paramsName
	builder.underParsing.paramsName = append(builder.underParsing.paramsName, param.Value())
	deposit, err := builder.ParseDeposit(node.children[1])
	builder.underParsing.paramsName = oldParamsName
	if limitErr != nil {
		err = limitErr
	}
	return &LoopFill{param, &limit, deposit}, err
}

//ParseDeposit parse deposit statement in the format <size directive> <expression>
//<size directive> ::= .db | .dw | .dd
func (builder *LangBuilder) ParseDeposit(node CSTNode) (DepositValue, *parsing.ParsingError) {
	sizeDirective := node.content[0]
	expression, err := NewExpression(node.children[0].content, builder.underParsing.paramsName, builder.enums)
	index := sort.SearchStrings([]string{".db", ".dd", ".dw"}, sizeDirective.Value())
	byteSize := [4]uint{1, 4, 2}
	deposit := DepositValue{&expression, byteSize[index]}
	return deposit, err
}

//ParseBranch parse if { <statements> } [else { <statements> }]
//[0 [Other] Expression:13 [IfBlock] BlockParser:12 [ElseBlock] BlockParser:12 ElseParser:8 BranchParser:10]
func (builder *LangBuilder) ParseBranch(node CSTNode) (BranchStatement, *parsing.ParsingError) {
	condition, conditionErr := NewExpression(node.children[0].content, builder.underParsing.paramsName, builder.enums)
	trueBlock, err := builder.ParseBlock(node.children[1])
	if err != nil {
		return BranchStatement{}, err
	}
	falseBlock := Block{}
	if len(node.children) > 2 {
		falseBlock, err = builder.ParseBlock(node.children[2].children[0])
		if err != nil {
			return BranchStatement{}, err
		}
	}
	err = conditionErr
	return BranchStatement{&condition, trueBlock, falseBlock}, err
}

//ParseError parse error message
func (builder *LangBuilder) ParseError(node CSTNode) (EmitError, *parsing.ParsingError) {
	msg := node.content[2]
	argRef := node.content[1]
	argIndex := indexOfString(argRef.Value(), builder.underParsing.paramsName)
	return EmitError{msg.Value(), argIndex}, nil
}

//ParseBlock parse statements between '{' and '}'
func (builder *LangBuilder) ParseBlock(cstNode CSTNode) (Block, *parsing.ParsingError) {
	var node OpcodeNode
	var nodes []OpcodeNode
	var err *parsing.ParsingError
	for _, child := range cstNode.children {
		tag := child.children[0].tag
		if tag == TagErrorStatement {
			node, err = builder.ParseError(child.children[0])
		} else if tag == TagBranchStatement {
			node, err = builder.ParseBranch(child.children[0])
		} else if tag == TagForStatement {
			node, err = builder.ParseLoop(child.children[0])
		} else if tag == TagDepositStatement {
			node, err = builder.ParseDeposit(child.children[0])
		} else {
			break
		}
		if err == nil {
			nodes = append(nodes, node)
		} else {
			break
		}
	}
	return Block{nodes}, err
}

//ParseOpcodeDeclaration parse the opcode declaration in the form name <token list> -> { <body> }
func (builder *LangBuilder) ParseOpcodeDeclaration(node CSTNode, report func(err *parsing.ParsingError)) bool {
	name := node.content[1]
	var header []text.Token
	for i := 2; i < len(node.content)-1; i++ {
		header = append(header, node.content[i])
	}
	tTypes := make([]text.TokenType, 0, len(header))
	tValues := make([]string, 0, len(header))
	isByteList := false
	for _, val := range header {
		tTypes = append(tTypes, val.EnumType())
		tValues = append(tValues, val.Value())
		if strings.EqualFold(val.Value(), ":") {
			err := parsing.NewParsingError(val, parsing.InvalidOpcodeArgument, ", invalid ':' found, allowed only after a label")
			report(&err)
		}
		if strings.EqualFold(val.Value(), byteListKeyword) {
			isByteList = true
		}
	}
	if isByteList && len(tValues) > 1 {
		err := parsing.NewParsingError(name, parsing.InvalidOpcodeFormat, ", the opcode accept a byte list, can't have another arguments")
		report(&err)
	}
	opcode := NewOpcode(name.Value(), tTypes, tValues, Block{})
	builder.underParsing = opcode
	oldParamsName := opcode.paramsName
	builder.underParsing.paramsName = append(builder.underParsing.paramsName, "this_address", "this_opcode")
	block, err := builder.ParseBlock(node.children[0])
	builder.underParsing.paramsName = oldParamsName
	if err == nil {
		opcode.nodes = block
		builder.opcodes = append(builder.opcodes, opcode)
	} else {
		report(err)
	}
	return err == nil
}
