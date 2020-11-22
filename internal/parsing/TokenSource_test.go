package parsing

import (
	"strings"
	"testing"

	"github.com/aleferri/casmeleon/internal/lexing"
	"github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/internal/ui"
)

func DefaultTestOptions() lexing.TokenMatchingOptions {
	const separators = ",;:{}()[]'\"\t "
	const operators = "+ - * / % ^ & | ~ > < ! <= >= == && || != << >> "
	const lineComment = "//"
	return lexing.NewMatchingOptions(operators, separators, operators+lineComment, lineComment)
}

func TestNewTokenBuffer(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	buffer := NewTokenBuffer("test", strings.NewReader(".number .hex prefix '$'\n"), DefaultTestOptions())
	number := buffer.NextToken(textUI)
	if number.EnumType() != text.Identifier {
		t.Errorf("Invalid match, expected String, found: %s\n", number.EnumType().String())
		t.Fail()
		return
	}
	if len(buffer.lastLine) < 3 {
		t.Error("Error, required 3\n")
		t.Fail()
		return
	}
	base := buffer.NextToken(textUI)
	position := buffer.NextToken(textUI)
	str := buffer.NextToken(textUI)
	if base.EnumType() != text.Identifier || position.EnumType() != text.Identifier || str.EnumType() != text.SingleQuotedString {
		t.Error("Invalid match, expected Strings, found other\n")
		t.Fail()
		return
	}
}

func TestTokenBuffer_NextToken(t *testing.T) {
	test := "+ \nb - a*4"
	userInterface := ui.NewConsole(false, false)
	buffer := NewTokenBuffer("test_expression", strings.NewReader(test), DefaultTestOptions())
	expected := []string{"+", "EOL", "b", "-", "a", "*", "4"}
	for i, s := range expected {
		found := buffer.NextToken(userInterface)
		t.Log(found.Value())
		if !strings.EqualFold(found.Value(), s) {
			t.Errorf("at index %d, expected %s, found %s\n", i, s, found.Value())
			t.Fail()
		}
	}
}
