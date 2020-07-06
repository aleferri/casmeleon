package langdef

import (
	"bitbucket.org/mrpink95/casmeleon/internal/parsing"
	"bitbucket.org/mrpink95/casmeleon/internal/text"
	"bitbucket.org/mrpink95/casmeleon/internal/ui"
	"strings"
	"testing"
)

func configureBuilder(buffer *parsing.TokenBuffer) (*LangBuilder, *LangParserState) {
	textUI := ui.NewConsoleUI(false, false)
	builder := NewLangBuilder()
	report := func(err *parsing.ParsingError) {
		textUI.ReportError(err.Error(), true)
	}
	state := NewParserState(textUI, buffer, true, report)
	return &builder, state
}

func parseTestAny(parser parsing.MatchRule, pState *LangParserState) []CSTNode {
	_, ok := parser([]text.Token{}, pState)
	if !ok {
		return nil
	}
	return pState.cst
}

func TestLangBuilder_Build(t *testing.T) {
	buffer := parsing.NewTokenBuffer("test_build", strings.NewReader("//VerySimpleCPU"), LangParsingDefaultOptions())
	builder, state := configureBuilder(buffer)
	builder.Build(state)
	if state.pUI.GetErrorCount() > 0 {
		t.Error("Expected no error\n")
		t.Fail()
	}
}

func TestLangBuilder_ParseNumericFormat(t *testing.T) {
	buffer := parsing.NewTokenBuffer("test_numbers_format", strings.NewReader(".number .hex prefix '$'\n.number .bin suffix 'b'"), LangParsingDefaultOptions())
	builder, pState := configureBuilder(buffer)
	numFormatParser := createNumberStatementParser()
	nodes := parseTestAny(numFormatParser, pState)
	if len(nodes[0].content) == 0 || !builder.ParseNumericFormat(nodes[0]) {
		t.Error("Not passed 1\n")
	}
	nodes = parseTestAny(numFormatParser, pState)
	if len(nodes[0].content) == 0 || !builder.ParseNumericFormat(nodes[0]) {
		t.Error("Not passed 2\n")
	}
	if len(builder.numberFormats) != 2 {
		t.Fail()
		return
	}
	println(builder.numberFormats[0].String())
	println(builder.numberFormats[1].String())
}

func TestLangBuilder_ParseEnum(t *testing.T) {
	buffer := parsing.NewTokenBuffer("test_enums", strings.NewReader(".enum regs\n{\nA\n,B,C,D\n,E,H\n,L}\n.enum ports\n{\nport_0,port_1}"), LangParsingDefaultOptions())
	builder, pState := configureBuilder(buffer)
	enumParser := createEnumDeclarationParser()
	nodes := parseTestAny(enumParser, pState)
	if len(nodes) == 0 || !builder.ParseEnum(nodes[0]) || len(builder.enums) != 1 {
		t.Fail()
		return
	}
	println("Enum found: ", builder.enums[0].String())
	nodes = parseTestAny(enumParser, pState)
	if len(nodes) < 2 || !builder.ParseEnum(nodes[1]) || len(builder.enums) != 2 {
		t.Fail()
		return
	}
	println("Enum found: ", builder.enums[1].String())
}

func TestLangBuilder_ParseError(t *testing.T) {
	buffer := parsing.NewTokenBuffer("test_error", strings.NewReader(".error\n \na \"invalid register\";"), LangParsingDefaultOptions())
	b, pState := configureBuilder(buffer)
	errorParser := createErrorStatementParser()
	nodes := parseTestAny(errorParser, pState)
	if nodes != nil {
		ee, r := b.ParseError(pState.cst[0])
		if r == nil {
			println("read error: " + ee.message)
			return
		}
	}
	t.Fail()
}

func TestLangBuilder_ParseLoop(t *testing.T) {
	buffer := parsing.NewTokenBuffer("test_loop", strings.NewReader("for\n  i\nuntil\n this_address +\n 4 .db i*4;"), LangParsingDefaultOptions())
	builder, pState := configureBuilder(buffer)
	builder.underParsing.paramsName = append(builder.underParsing.paramsName, "this_address", "this_opcode")
	loopParser := createLoopParser(MatchExpression, createDepositStatementParser(MatchExpression))
	nodes := parseTestAny(loopParser, pState)
	if nodes != nil {
		loop, r := builder.ParseLoop(pState.cst[0])
		if r == nil {
			print("Index name: ", loop.indexName.Value(), " -- ")
			loop.limit.DebugPrint()
			return
		}
		println(r.Error())
	}
	t.Fail()
}

func TestLangBuilder_ParseDeposit(t *testing.T) {
	buffer := parsing.NewTokenBuffer("test_deposit", strings.NewReader(".dw 16*a % 256;"), LangParsingDefaultOptions())
	builder, pState := configureBuilder(buffer)
	builder.underParsing.paramsName = []string{"a"}
	depositParser := createDepositStatementParser(MatchExpression)
	nodes := parseTestAny(depositParser, pState)
	if nodes != nil {
		ee, r := builder.ParseDeposit(pState.cst[0])
		if r == nil {
			ee.e.DebugPrint()
			return
		}
	}
	t.Fail()
}

func TestLangBuilder_ParseBlock(t *testing.T) {
	buffer := parsing.NewTokenBuffer("test_block", strings.NewReader("{ .dw 16*a % 256; .error a\"invalid a\"; }"), LangParsingDefaultOptions())
	builder, pState := configureBuilder(buffer)
	builder.underParsing.paramsName = []string{"a"}
	blockParser, _ := createBlockParser()
	nodes := parseTestAny(blockParser, pState)
	if nodes != nil {
		block, r := builder.ParseBlock(pState.cst[0])
		if r == nil && len(block.nodes) == 2 {
			println("Block parsed successfully")
			return
		}
		println("Block: error during tree building")
	}
	t.Fail()
}

func TestLangBuilder_ParseBranch(t *testing.T) {
	buffer := parsing.NewTokenBuffer("test_branch", strings.NewReader("if a > b { .error a\"invalid a\"; } else { .db a * b; }"), LangParsingDefaultOptions())
	builder, pState := configureBuilder(buffer)
	builder.underParsing.paramsName = []string{"a", "b"}
	_, branchParser := createBlockParser()
	nodes := parseTestAny(branchParser, pState)
	if nodes != nil {
		branch, r := builder.ParseBranch(pState.cst[0])
		if r == nil && len(branch.conditionMet.nodes) == 1 && len(branch.conditionNotMet.nodes) == 1 {
			return
		}
	}
	t.Fail()
}

func TestLangBuilder_ParseOpcodeDeclaration(t *testing.T) {
	buffer := parsing.NewTokenBuffer("test_opcode", strings.NewReader(".opcode inx -> { .db 0xFF; }"), LangParsingDefaultOptions())
	builder, pState := configureBuilder(buffer)
	blockParser, _ := createBlockParser()
	opcodeParser := createOpcodeParser(blockParser)
	nodes := parseTestAny(opcodeParser, pState)
	if nodes != nil {
		r := builder.ParseOpcodeDeclaration(pState.cst[0], pState.report)
		if r {
			println(builder.opcodes[0].name)
			builder.opcodes[0].nodes.DebugPrint()
			return
		}
	}
	t.Fail()
}

func TestLangBuilder_ParseOpcodeDeclaration2(t *testing.T) {
	buffer := parsing.NewTokenBuffer("test_opcode", strings.NewReader(".opcode add #imm -> { .db 0xFF; .db imm; }"), LangParsingDefaultOptions())
	builder, pState := configureBuilder(buffer)
	blockParser, _ := createBlockParser()
	opcodeParser := createOpcodeParser(blockParser)
	nodes := parseTestAny(opcodeParser, pState)
	if nodes != nil {
		r := builder.ParseOpcodeDeclaration(pState.cst[0], pState.report)
		if r {
			for i, k := range builder.underParsing.paramsString {
				println(k, builder.underParsing.paramsFormat[i].String())
			}
			return
		}
	}
	t.Fail()
}
