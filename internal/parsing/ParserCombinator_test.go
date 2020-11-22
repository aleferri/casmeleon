package parsing

import (
	"strings"
	"testing"

	"github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/internal/ui"
)

//TestingState is a simplified ParserState used to test parser combinators
type TestingState struct {
	n        int
	tags     []int
	src      TokenSource
	reportUI ui.UI
}

func NewTestingState(s string) *TestingState {
	textUI := ui.NewConsole(false, false)
	buffer := NewTokenBuffer("test", strings.NewReader(s), DefaultTestOptions())
	return &TestingState{0, nil, buffer, textUI}
}

func (t *TestingState) Push(tokens []text.Token, tag int) {
	t.n += len(tokens)
	t.tags = append(t.tags, tag)
}

func (t *TestingState) Pull(tokens *[]text.Token) int {
	lex := t.src.NextToken(t.reportUI)
	lex = lex.WithType(FindTokenType(lex))
	*tokens = append(*tokens, lex)
	if lex.EnumType() == text.EOF {
		return 0
	}
	return 1
}

func (t *TestingState) ReportError(wrong text.Token, err ParsingErrorType, msg string) {
	if err == NoError {
		return
	}
	err.ReportToken(t.reportUI, wrong, t.src, msg)
}

func (t *TestingState) Mark() (int, int) {
	return 0, 0
}

func (t *TestingState) ToCST(int, int, int) {

}

func TestMatchAnyString(t *testing.T) {
	tState := NewTestingState("random")
	err := NewIncompleteError(ErrorExpectedToken, ", expected 'random'")
	parseAnyString := MatchAnyString(err, "random")
	parseAnyString([]text.Token{}, tState)
	if tState.n != 1 {
		t.Fail()
	}
}

func TestMatchAnyToken(t *testing.T) {
	tState := NewTestingState("-")
	err := NewIncompleteError(ErrorExpectedToken, ", expected '+' or '-'")
	parseAnyString := MatchAnyString(err, "+", "-")
	parseAnyString([]text.Token{}, tState)
	if tState.n != 1 {
		t.Fail()
	}
}

func TestMatchNotAnyToken(t *testing.T) {
	tState := NewTestingState("*")
	err := NewIncompleteError(ErrorExpectedToken, ", unexpected '*'")
	parseAnyString := MatchNotAnyToken(err, text.Number, text.SymbolHash, text.SymbolArrow)
	parseAnyString([]text.Token{}, tState)
	if tState.n != 1 {
		t.Fail()
	}
}

func TestMatchAll(t *testing.T) {
	tState := NewTestingState("nick 44")
	nameError := NewIncompleteError(ErrorExpectedToken, "expected name")
	numberError := NewIncompleteError(ErrorExpectedToken, "expected age")
	name := MatchToken(nameError, text.Identifier)
	age := MatchToken(numberError, text.Number)
	parsePerson := MatchAll(1, name, age)
	parsePerson([]text.Token{}, tState)
	if tState.n != 2 {
		t.Fail()
	}
}

func TestTryMatch(t *testing.T) {
	tState := NewTestingState("nick 44")
	nameError := NewIncompleteError(ErrorExpectedToken, "expected name")
	numberError := NewIncompleteError(ErrorExpectedToken, "expected age")
	testName := TryMatchToken(text.Identifier, false)
	name := MatchToken(nameError, text.Identifier)
	age := MatchToken(numberError, text.Number)
	parsePerson := MatchAll(1, name, age)
	tryParsePerson := TryMatch(1, testName, parsePerson)
	tryParsePerson([]text.Token{}, tState)
	if tState.n != 2 {
		t.Fail()
	}
}

func TestTryMatchRepeat(t *testing.T) {
	tState := NewTestingState("a, b, c, d,e, f, g")
	nameError := NewIncompleteError(ErrorExpectedToken, ", expected identifier")
	testComma := TryMatchToken(text.Comma, true)
	name := MatchToken(nameError, text.Identifier)
	parseLeft := TryMatchRepeat(1, testComma, name)
	parseList := MatchAll(2, name, parseLeft)
	parseList([]text.Token{}, tState)
	if tState.n != 7 {
		t.Fail()
	}
}
