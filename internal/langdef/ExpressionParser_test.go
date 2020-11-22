package langdef

import (
	"strings"
	"testing"

	"github.com/aleferri/casmeleon/internal/parsing"
	"github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/internal/ui"
)

func testPush(t text.Token) {
	print(t.Value() + " ")
}

func TestLangBuilder_ParseExpressionOperand(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	buffer := parsing.NewTokenBuffer("test_expression", strings.NewReader("b"), LangParsingDefaultOptions())
	pState := NewParserState(textUI, buffer, true, nil)
	_, r := matchOperand([]text.Token{}, pState)
	print("Start: ")
	for _, t := range pState.lastParsed {
		testPush(t)
	}
	println("End")
	if !r {
		t.Fail()
		return
	}
}

func TestLangBuilder_ParseExpressionUnary(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	buffer := parsing.NewTokenBuffer("test_expression", strings.NewReader("+ \nb"), LangParsingDefaultOptions())
	pState := NewParserState(textUI, buffer, true, nil)
	_, r := matchUnaryOperation([]text.Token{}, pState)
	print("Start: ")
	for _, t := range pState.lastParsed {
		testPush(t)
	}
	println("End")
	if !r {
		t.Fail()
		return
	}
}

func TestLangBuilder_ParseExpressionBinary(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	buffer := parsing.NewTokenBuffer("test_expression", strings.NewReader("+ \nb - a*4"), LangParsingDefaultOptions())
	pState := NewParserState(textUI, buffer, true, nil)
	_, r := matchUnaryOperation([]text.Token{}, pState)
	_, r2 := matchBinaryOperation([]text.Token{}, pState, 0)
	r = r && r2
	print("Start: ")
	for _, t := range pState.lastParsed {
		testPush(t)
	}
	println("End")
	if !r {
		t.Fail()
		return
	}
}

func TestLangBuilder_ParseExpression2(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	buffer := parsing.NewTokenBuffer("test_expression", strings.NewReader("adr < 0 || (adr > 0)"), LangParsingDefaultOptions())
	pState := NewParserState(textUI, buffer, true, nil)
	_, r := MatchExpression([]text.Token{}, pState)
	print("Start: ")
	for _, t := range pState.cst[0].content {
		testPush(t)
	}
	println("End")
	if !r {
		t.Fail()
		return
	}
}
