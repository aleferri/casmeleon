package lexer

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

var mapOfScanners = map[rune]Scanner{
	'0':  ScanNumber,
	'1':  ScanBase10Number,
	'2':  ScanBase10Number,
	'3':  ScanBase10Number,
	'4':  ScanBase10Number,
	'5':  ScanBase10Number,
	'6':  ScanBase10Number,
	'7':  ScanBase10Number,
	'8':  ScanBase10Number,
	'9':  ScanBase10Number,
	'\n': ScanWhite,
	'\r': ScanWhite,
	' ':  ScanWhite,
	'\t': ScanWhite,
	'&':  ScanOperator,
	'|':  ScanOperator,
	'^':  ScanOperator,
	'!':  ScanOperator,
	'<':  ScanOperator,
	'>':  ScanOperator,
	'=':  ScanOperatorSingle,
	'*':  ScanOperatorSingle,
	'+':  ScanOperatorSingle,
	'/':  ScanSlash,
	'-':  ScanOperatorSingle,
	'%':  ScanOperatorSingle,
	'.':  ScanOperatorSingle,
	'@':  ScanSeparator,
	'#':  ScanSeparator,
	',':  ScanSeparator,
	';':  ScanSeparator,
	':':  ScanSeparator,
	'(':  ScanSeparator,
	')':  ScanSeparator,
	'[':  ScanSeparator,
	']':  ScanSeparator,
	'{':  ScanSeparator,
	'}':  ScanSeparator,
	'"':  ScanQuotedString,
}

func TestScanNumber(t *testing.T) {
	runes := []rune("0b10\n")
	id, left, err := ScanNumber(runes, true)
	t.Logf("Scanned id of %d and left in the buffer %v", id, left)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestLexer(t *testing.T) {
	source := bufio.NewReader(strings.NewReader("0b10 * aaa - ccc/ 0xF& 0xD|| 0o7., 0q3= >=a_% $r@#++^ 100_000\n"))
	lexer := Of(mapOfScanners, ScanIdentifier)

	line, _ := source.ReadBytes('\n')
	runes := bytes.Runes(line)

	lexed, err := lexer.Scan(0, 0, runes, true)
	if err != nil {
		t.Error(err)
		return
	}

	symbols := lexed.Symbols()
	if len(symbols) < 20 {
		t.Errorf("%v \nExpected 20 symbols, found %d", symbols, len(symbols))
		return
	}

}

func TestLexerSource(t *testing.T) {
	fileName := "../../tests/lexer_dump.txt"
	var file, fileErr = os.Open(fileName)
	if fileErr != nil {
		wnd, _ := os.Getwd()
		t.Errorf("Error during opening of file %s from %s\n", fileName, wnd)
		return
	}
	source := bufio.NewReader(file)
	lexer := Of(mapOfScanners, ScanIdentifier)

	line, ioErr := source.ReadBytes('\n')
	buffer := bytes.Runes(line)

	var lexed *Result
	var err error

	for ioErr == nil {
		lexed, err = lexer.Scan(0, 0, buffer, false)
		if err != nil {
			t.Errorf("Got an error %s from line\n%s\n", err, string(lexed.Unlexed()))
			return
		}
		line, ioErr = source.ReadBytes('\n')
		if len(lexed.Unlexed()) == 0 {
			buffer = bytes.Runes(line)
		} else {
			buffer = append(lexed.Unlexed(), bytes.Runes(line)...)
		}
	}

	lexed, err = lexer.Scan(0, 0, buffer, true)
	if err != nil {
		t.Error(err)
		return
	}
}

func BenchmarkLexerSource(b *testing.B) {
	fileName := "../../tests/lexer_dump.txt"
	var file, fileErr = os.Open(fileName)
	if fileErr != nil {
		wnd, _ := os.Getwd()
		b.Errorf("Error during opening of file %s from %s\n", fileName, wnd)
		return
	}
	source := bufio.NewReader(file)
	lexer := Of(mapOfScanners, ScanIdentifier)

	line, ioErr := source.ReadBytes('\n')
	buffer := bytes.Runes(line)

	var lexed *Result
	var err error

	for ioErr == nil {
		lexed, err = lexer.Scan(0, 0, buffer, false)
		if err != nil {
			b.Errorf("Got an error %s from line\n%s\n", err, string(lexed.Unlexed()))
			return
		}
		line, ioErr = source.ReadBytes('\n')
		if len(lexed.Unlexed()) == 0 {
			buffer = bytes.Runes(line)
		} else {
			slice := []rune{}
			slice = append(slice, lexed.Unlexed()...)
			buffer = append(slice, bytes.Runes(line)...)
			b.Errorf("Probable error, since we are reading line by line")
			return
		}
	}

	lexed, err = lexer.Scan(0, 0, buffer, true)
	if err != nil {
		b.Error(err)
		return
	}
}
