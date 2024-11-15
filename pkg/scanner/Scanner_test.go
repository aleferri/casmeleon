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
	'&':  FollowSequence('&', '&'),
	'|':  FollowSequence('|', '|'),
	'^':  FollowNone,
	'!':  FollowSequence('!', '='),
	'<':  FollowComparison,
	'>':  FollowComparison,
	'=':  FollowSequence('=', '='),
	'*':  FollowCommentClose,
	'+':  FollowNone,
	'/':  FollowCommentOpen,
	'-':  FollowSequence('-', '>'),
	'%':  FollowNone,
	'@':  FollowNone,
	'#':  FollowNone,
	',':  FollowNone,
	';':  FollowNone,
	':':  FollowNone,
	'(':  FollowNone,
	')':  FollowNone,
	'[':  FollowNone,
	']':  FollowNone,
	'{':  FollowSequence('{', '{'),
	'}':  FollowSequence('}', '}'),
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

func TestFullPipeline(t *testing.T) {
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

	completed := []Token{}
	last := int32(-1)
	var merged *Token = nil

	for ioErr == nil {
		tokens, left := FastScan(buffer, false, FromMap(followMap))

		ClassifyMergeableTokens(tokens)
		completed, merged, last = Merge(map[int32]int32{1: 1, 2: 2, 3: 3, 4: 5}, tokens, merged, last)

		line, ioErr = source.ReadBytes('\n')
		if len(left) == 0 {
			buffer = bytes.Runes(line)
		} else {
			buffer = append(left, bytes.Runes(line)...)
		}

		for _, c := range completed {
			str := c.String()
			if str == "\n" || str == "\r" {
				fmt.Println()
			} else {
				fmt.Printf("%s", str)
			}
		}
	}

	tokens, _ := FastScan(buffer, true, FromMap(followMap))
	ClassifyMergeableTokens(tokens)
	completed, merged, last = Merge(map[int32]int32{1: 1, 2: 2, 3: 3, 4: 5}, tokens, merged, last)
}
