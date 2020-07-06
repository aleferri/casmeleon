package lexing

import (
	"container/list"
	"strings"
	"testing"
)

func CheckEquals(a *list.List, b []string, t *testing.T) bool {
	if a.Len() != len(b) {
		t.Errorf("Different in length: %d <-> %d\n", a.Len(), len(b))
		return false
	}
	var i = 0
	var size = a.Len()
	for i < size {
		var rElem = a.Remove(a.Front()).(string)
		var eElem = b[i]
		if strings.Compare(rElem, eElem) != 0 {
			t.Errorf("Different element: '%s' <-> '%s' \n", rElem, eElem)
			return false
		}
		i++
	}
	return true
}

func CheckEqualsSlice(a, b []string, t *testing.T) bool {
	if len(a) != len(b) {
		t.Errorf("Different in length: %d <-> %d\n", len(a), len(b))
		return false
	}
	var i = 0
	for i < len(a) {
		var rElem = a[i]
		var eElem = b[i]
		if strings.Compare(rElem, eElem) != 0 {
			t.Errorf("Different element: '%s' <-> '%s'\n", rElem, eElem)
			return false
		}
		i++
	}
	return true
}

func TestSourceReader(t *testing.T) {
	var source = ".type '#' literal\n" +
		".type '@+' stack"
	var sourceReader = NewSourceReader(strings.NewReader(source), "#@+' ", "Test.string")
	var line1 = []string{".type", " ", "'", "#", "'", " ", "literal"}
	var line2 = []string{".type", " ", "'", "@", "+", "'", " ", "stack"}
	var lines = [2][]string{line1, line2}
	var lineNumber = uint(0)
	for sourceReader.HasNext() {
		var line = sourceReader.NextLine()
		if lineNumber+1 != sourceReader.lineNumber {
			t.Error("Different line number")
			t.Fail()
		}
		if !CheckEquals(line, lines[lineNumber], t) {
			t.Error("Different list")
			t.Fail()
		}
		lineNumber++
	}
}

func TestSourceReader2(t *testing.T) {
	var source = ".number .hex prefix '$'\n"
	var sourceReader = NewSourceReader(strings.NewReader(source), " '", "Test.string")
	var line = []string{".number", " ", ".hex", " ", "prefix", " ", "'$'"}
	var ln = sourceReader.NormalizeLine(sourceReader.NextLine(), sourceReader.LineNumber(), "&& || ++ -- >= <= >> << >>>")
	if !CheckEqualsSlice(line, ln.AllStrings(), t) {
		t.Error("Invalid line\n")
		t.Errorf("%v -- %v\n", ln.AllStrings(), line)
		t.Fail()
	}
}

func TestSourceReaderComplete(t *testing.T) {
	var source = ".type '#' literal\n" +
		".type '@+' stack;!="
	var sourceReader = NewSourceReader(strings.NewReader(source), "#@+!=;' ", "Test.string")
	var line1 = []string{".type", " ", "'#'", " ", "literal"}
	var line2 = []string{".type", " ", "'@+'", " ", "stack", ";", "!="}
	var lines = [2][]string{line1, line2}
	var options = NewMatchingOptions("@ + # !=", "", "", ";")
	var lineNumber = uint(0)
	for sourceReader.HasNext() {
		var line = sourceReader.NextLine()
		var normalized = sourceReader.NormalizeLine(line, lineNumber, options.operators)
		if !CheckEqualsSlice(normalized.AllStrings(), lines[lineNumber], t) {
			t.Error("Invalid line\n")
			t.Fail()
		}
		sourceReader.FilterLine(normalized, options)
		lineNumber++
	}
}
