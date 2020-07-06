package lexing

import (
	"strings"
	"testing"
)

func TestSplitKeepSeparators(t *testing.T) {
	var result = SplitKeepSeparators("a; bb + ad", ";+")
	var expected = []string{"a", ";", " bb ", "+", " ad"}
	if !CheckEquals(result, expected, t) {
		t.Fail()
	}
}

func TestJoinQuote(t *testing.T) {
	var test = "a 'bnf ' + \"&'\""
	var source = NewSourceReader(strings.NewReader(test), "'\"+& ", "Test.string")
	var l = JoinQuote(source.NextLine())
	var expected = []string{"a", " ", "'bnf '", " ", "+", " ", "\"&'\""}
	if !CheckEquals(l, expected, t) {
		t.Fail()
	}
}

func TestJoinQuote2(t *testing.T) {
	var test = ".number .hex prefix '$'\n"
	var source = NewSourceReader(strings.NewReader(test), "'\"+& ", "Test.string")
	var l = JoinQuote(source.NextLine())
	var expected = []string{".number", " ", ".hex", " ", "prefix", " ", "'$'"}
	if !CheckEquals(l, expected, t) {
		t.Fail()
	}
}

func TestJoinMultiCharSymbols(t *testing.T) {
	var test = "a && b +++e %||>="
	var source = NewSourceReader(strings.NewReader(test), "&+%|>= ", "Test.string")
	var l = JoinMultiCharSymbols(source.NextLine(), "&&  ||  >=")
	var expected = []string{"a", " ", "&&", " ", "b", " ", "+", "+", "+", "e", " ", "%", "||", ">="}
	if !CheckEquals(l, expected, t) {
		t.Fail()
	}
}

func TestRegroupSymbols(t *testing.T) {
	var test = "a && b +++e %||>="
	var source = NewSourceReader(strings.NewReader(test), "&+%|>= ", "Test.string")
	var regroup = RegroupSymbols(source.NextLine(), "&& || >=")
	var expected = []string{"a", " ", "&&", " ", "b", " ", "+", "+", "+", "e", " ", "%", "||", ">="}
	if !CheckEqualsSlice(regroup, expected, t) {
		t.Errorf("%v\n", regroup)
		t.Fail()
	}
}

func TestRegroupSymbols2(t *testing.T) {
	var test = ".number .hex prefix '$'\n"
	var source = NewSourceReader(strings.NewReader(test), "'\"+& ", "Test.string")
	var slice = RegroupSymbols(JoinQuote(source.NextLine()), "|| && >= <= == !=")
	var expected = []string{".number", " ", ".hex", " ", "prefix", " ", "'$'"}
	if !CheckEqualsSlice(slice, expected, t) {
		t.Fail()
	}
}

func TestRegroupSymbols3(t *testing.T) {
	var test = "+ \nb"
	var source = NewSourceReader(strings.NewReader(test), "+ ", "Test.string")
	var slice = RegroupSymbols(JoinQuote(source.NextLine()), "+ ")
	slice = append(slice, RegroupSymbols(JoinQuote(source.NextLine()), "+ ")...)
	var expected = []string{"+", " ", "b"}
	if !CheckEqualsSlice(slice, expected, t) {
		t.Fail()
	}
}
