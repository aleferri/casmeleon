package casm

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

func TestParserSource(t *testing.T) {
	fileName := "../../tests/parser_test.casm"
	var file, fileErr = os.Open(fileName)
	if fileErr != nil {
		wnd, _ := os.Getwd()
		t.Errorf("Error during opening of file %s from %s\n", fileName, wnd)
		return
	}
	source := bufio.NewReader(file)
	repo := text.BuildSource("parser_test.casm")

	stream := BuildStream(source, &repo)

	errFound := false

	id := stream.Peek().ID()

	for id != text.EOF && !errFound {
		var cst parser.CSTNode
		var e error

		switch id {
		case text.KeywordInline:
			{
				cst, e = ParseInline(stream)
			}
		case text.KeywordOpcode:
			{
				cst, e = ParseOpcode(stream)
			}
		case text.KeywordNum:
			{
				cst, e = ParseNumberBase(stream)
			}
		case text.KeywordSet:
			{
				cst, e = ParseSet(stream)
			}
		default:
			{
				fmt.Println(len(idDescriptor))
				e = fmt.Errorf("Undefined symbol '%s'", idDescriptor[id])
			}
		}

		if e != nil {
			errFound = true
			parseErr, ok := e.(*ParserError)
			if !ok {
				fmt.Println("Unexpected Error")
			} else {
				parseErr.PrettyPrint(&repo)
			}
		} else {
			fmt.Println(cst.ID())
		}

		id = stream.Peek().ID()
	}
	t.Fail()

	if errFound {
		t.Errorf("Errors have been found")
	}
}
