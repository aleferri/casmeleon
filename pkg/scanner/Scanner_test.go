package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"
)

var followMap = map[rune]Follow{
	'\n': FollowNone,
	'\r': FollowNone,
	' ':  FollowSpaces,
	'\t': FollowNone,
	'&':  FollowSame('&'),
	'|':  FollowSame('|'),
	'^':  FollowNone,
	'!':  FollowSequence('!', '='),
	'<':  FollowComparison,
	'>':  FollowComparison,
	'=':  FollowSame('='),
	'*':  FollowNone,
	'+':  FollowNone,
	'/':  FollowComments,
	'-':  FollowSequence('-', '>'),
	'%':  FollowNone,
	'.':  FollowNone,
	'@':  FollowNone,
	'#':  FollowNone,
	',':  FollowNone,
	';':  FollowNone,
	':':  FollowNone,
	'(':  FollowNone,
	')':  FollowNone,
	'[':  FollowNone,
	']':  FollowNone,
	'{':  FollowSame('{'),
	'}':  FollowSame('}'),
	'"':  FollowNone,
	'\'': FollowNone,
}

func TestFastScan(t *testing.T) {
	fileName := "../../tests/parser_test.casm"
	var file, fileErr = os.Open(fileName)
	if fileErr != nil {
		wnd, _ := os.Getwd()
		t.Errorf("Error during opening of file %s from %s\n", fileName, wnd)
		return
	}
	source := bufio.NewReader(file)

	line, ioErr := source.ReadBytes('\n')
	buffer := bytes.Runes(line)

	count := 0

	for ioErr == nil {
		tokens, left := FastScan(buffer, false, FromMap(followMap))

		for _, t := range tokens {
			count += len(t.slice)
		}

		line, ioErr = source.ReadBytes('\n')
		if len(left) == 0 {
			buffer = bytes.Runes(line)
		} else {
			buffer = append(left, bytes.Runes(line)...)
		}
	}

	tokens, _ := FastScan(buffer, true, FromMap(followMap))
	for _, t := range tokens {
		count += len(t.slice)
	}
	fmt.Println("Complete")
}
