package lexing

import (
	"strings"
	"testing"

	"github.com/aleferri/casmeleon/internal/text"
)

func ExpectedTokenType(tType text.TokenType, expected string, t *testing.T) {
	if !strings.EqualFold(tType.String(), expected) {
		t.Errorf("Error: wrong token type; found %s, expected %s\n", tType.String(), expected)
		t.Fail()
	}
}

func TestFindTokenType(t *testing.T) {
	var test = []string{"add", "+", "8", ";", "0xFF", "'@$'"}
	var expectedType = []string{"IDENTIFIER", "BINARY", "NUMBER", "SEPARATOR", "NUMBER", "SINGLE_QUOTED"}
	for i := 0; i < len(test); i++ {
		var tType = FindRawTokenType(test[i], "+@", ",;")
		ExpectedTokenType(tType, expectedType[i], t)
	}
}

func TestTokenMatchingOptions_LineToTokens(t *testing.T) {
	var test = "add + 8, 5\t;0xFF,'@$'\n"
	var source = NewSourceReader(strings.NewReader(test), "+;,'@\t ", "Test.string")
	var ls = RegroupSymbols(JoinQuotes(source.NextLine()), "")
	var line = text.NewSourceLine(ls, source.LineNumber(), source.sourceName)
	var opt = NewMatchingOptions("+@", ";,", "$", ";")
	var slice = opt.LineToTokens(line, 1, "Test.string")
	var expected = []string{"IDENTIFIER", "BINARY", "NUMBER", "SEPARATOR", "NUMBER"}
	if len(slice) != len(expected) {
		t.Errorf("Invalid len %d <-> %d\n", len(slice), len(expected))
		t.Fail()
		for i := 0; i < len(slice); i++ {
			t.Logf("%s\n", slice[i].Value())
		}
		return
	}
	for i := 0; i < len(slice); i++ {
		var toString = slice[i].EnumType().String()
		if !strings.EqualFold(toString, expected[i]) {
			t.Errorf("Unexpected token type %s <-> %s\n", toString, expected[i])
			t.Fail()
		}
	}
}

func TestTokenMatchingOptions_LineToTokens2(t *testing.T) {
	var test = "//All this line is a comment"
	var source = NewSourceReader(strings.NewReader(test), "//\t ", "Test.string")
	var ls = RegroupSymbols(JoinQuotes(source.NextLine()), "//")
	var line = text.NewSourceLine(ls, source.LineNumber(), source.sourceName)
	var opt = NewMatchingOptions("//", "", "", "//")
	var slice = opt.LineToTokens(line, 1, "Test.string")
	if len(slice) != 0 {
		t.Error("Expected 0 length line")
		t.Fail()
	}
}
