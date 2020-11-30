package casm

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"unicode/utf8"

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
			switch e.(type) {
			case *parser.MatchError:
				{
					m := e.(*parser.MatchError)
					v := idDescriptor[m.Expected()]
					fileName, lineIndex, column := repo.FindPosition(m.Found())
					fmt.Printf("In file %s, error at %d, %d: ", fileName, lineIndex+1, column+1)
					fmt.Printf(m.Error(), m.Found().Value(), v)
					fmt.Println()
				}
			case *ParserError:
				{
					p := e.(*ParserError)
					m := p.internal
					v := idDescriptor[m.Expected()]
					fileName, lineIndex, column := repo.FindPosition(m.Found())
					fmt.Printf("In file %s, error at %d, %d: ", fileName, lineIndex+1, column+1)
					fmt.Printf(m.Error(), m.Found().Value(), v)
					fmt.Println()
					context := repo.SliceScope(m.Found(), p.context)
					for _, c := range context {
						fmt.Print(c.Value())
					}
					fmt.Println()
					line := repo.SliceLine(m.Found())
					for _, c := range line[1:] {
						char := ""
						if c.Equals(m.Found()) {
							char = "^"
						} else if c.ID() > 2 {
							char = "-"
						}
						runeCount := utf8.RuneCountInString(c.Value())
						for k := 0; k < runeCount; k++ {
							fmt.Print(char)
						}
					}
					fmt.Println()
				}
			default:
				{
					fmt.Println(e.Error())
					context := stream.Source().SliceScope(stream.Peek(), "{")
					for _, t := range context {
						fmt.Print(t.Value())
					}
					fmt.Println()
				}
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
